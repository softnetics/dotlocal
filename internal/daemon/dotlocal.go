package daemon

import (
	"context"
	"time"

	"github.com/dchest/uniuri"
	"github.com/samber/lo"
	"github.com/softnetics/dotlocal/internal"
	"github.com/softnetics/dotlocal/internal/daemon/dnsproxy"
	"github.com/softnetics/dotlocal/internal/daemon/orbdnsproxy"
	"go.uber.org/zap"
	"gopkg.in/tomb.v2"
)

type DotLocal struct {
	logger   *zap.Logger
	nginx    *Nginx
	dnsProxy dnsproxy.DNSProxy
	mappings map[internal.MappingKey]*internal.MappingState
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewDotLocal(logger *zap.Logger) (*DotLocal, error) {
	nginx, err := NewNginx(logger.Named("nginx"))
	if err != nil {
		return nil, err
	}

	dnsProxy, err := orbdnsproxy.NewOrbstackDNSProxy(logger.Named("orbdnsproxy"))
	if err != nil {
		return nil, err
	}

	return &DotLocal{
		logger:   logger,
		nginx:    nginx,
		dnsProxy: dnsProxy,
		mappings: make(map[internal.MappingKey]*internal.MappingState),
	}, nil
}

func (d *DotLocal) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	d.ctx = ctx
	d.cancel = cancel

	var t tomb.Tomb
	t.Go(func() error {
		return d.nginx.Start()
	})
	t.Go(func() error {
		return d.dnsProxy.Start(d.nginx.Port())
	})

	err := t.Wait()
	if err != nil {
		return err
	}

	d.logger.Info("Ready")

	go func() {
		for {
			timer := time.NewTimer(1 * time.Minute)
			select {
			case <-timer.C:
				err := d.removeExpiredMappings()
				if err != nil {
					d.logger.Error("Failed to update mappings", zap.Error(err))
				}
			case <-d.ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (d *DotLocal) GetMappings() []internal.Mapping {
	return lo.MapToSlice(d.mappings, func(key internal.MappingKey, state *internal.MappingState) internal.Mapping {
		return internal.NewMapping(key, state)
	})
}

func (d *DotLocal) CreateMapping(opts internal.MappingOptions) (internal.Mapping, error) {
	if opts.PathPrefix == "" {
		opts.PathPrefix = "/"
	}
	key := internal.MappingKey{
		Host:       opts.Host,
		PathPrefix: opts.PathPrefix,
	}
	expiresAt := time.Now().Add(2 * time.Minute)

	state, ok := d.mappings[key]
	if ok {
		state.ExpiresAt = expiresAt
	} else {
		state = &internal.MappingState{
			ID:        uniuri.NewLen(6),
			Target:    "",
			ExpiresAt: expiresAt,
		}
		d.mappings[key] = state
	}
	previousMapping := internal.NewMapping(key, state)
	state.Target = opts.Target

	mapping := internal.NewMapping(key, state)

	if previousMapping == mapping {
		return mapping, nil
	}
	d.logger.Info("Created mapping", zap.Any("mapping", mapping))
	return mapping, d.UpdateMappings()
}

func (d *DotLocal) RemoveMapping(keys ...internal.MappingKey) error {
	for _, key := range keys {
		delete(d.mappings, key)
	}
	return d.UpdateMappings()
}

func (d *DotLocal) UpdateMappings() error {
	mappings := d.GetMappings()
	err := d.nginx.SetMappings(mappings)
	if err != nil {
		return err
	}
	err = d.dnsProxy.SetHosts(lo.SliceToMap(mappings, func(mapping internal.Mapping) (string, struct{}) {
		return mapping.Host, struct{}{}
	}))
	if err != nil {
		return err
	}
	return nil
}

func (d *DotLocal) removeExpiredMappings() error {
	var expiredMappings []internal.MappingKey
	for key, state := range d.mappings {
		if state.ExpiresAt.Before(time.Now()) {
			expiredMappings = append(expiredMappings, key)
		}
	}
	if len(expiredMappings) == 0 {
		return nil
	}
	return d.RemoveMapping(expiredMappings...)
}
