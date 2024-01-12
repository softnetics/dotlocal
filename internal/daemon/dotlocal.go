package daemon

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/dchest/uniuri"
	"github.com/samber/lo"
	"github.com/softnetics/dotlocal/internal"
	api "github.com/softnetics/dotlocal/internal/api/proto"
	"github.com/softnetics/dotlocal/internal/daemon/dnsproxy"
	"github.com/softnetics/dotlocal/internal/daemon/mdnsproxy"
	"github.com/softnetics/dotlocal/internal/util"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
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

	dnsProxy, err := mdnsproxy.NewMDNSProxy(logger.Named("dnsproxy"))
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

func (d *DotLocal) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	d.ctx = ctx
	d.cancel = cancel

	preferences, err := d.loadPreferences()
	if err != nil {
		return err
	}
	for _, mapping := range preferences.Mappings {
		key := internal.MappingKey{
			Host:       *mapping.Host,
			PathPrefix: *mapping.PathPrefix,
		}
		state := &internal.MappingState{
			ID:        *mapping.Id,
			Target:    *mapping.Target,
			ExpiresAt: mapping.ExpiresAt.AsTime(),
		}
		if state.ExpiresAt.Before(time.Now()) {
			continue
		}
		d.mappings[key] = state
	}

	var t tomb.Tomb
	t.Go(func() error {
		return d.nginx.Start(ctx)
	})
	t.Go(func() error {
		return d.dnsProxy.Start(ctx)
	})

	err = t.Wait()
	if err != nil {
		return err
	}

	err = d.UpdateMappings()
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

func (d *DotLocal) Stop() error {
	var t tomb.Tomb
	t.Go(func() error {
		return d.nginx.Stop()
	})
	t.Go(func() error {
		return d.dnsProxy.Stop()
	})
	err := t.Wait()
	if err != nil {
		return err
	}
	d.logger.Info("Stopped")
	err = d.savePreferences()
	if err != nil {
		return err
	}
	return nil
}

func (d *DotLocal) loadPreferences() (*api.Preferences, error) {
	json, err := os.ReadFile(util.GetPreferencesPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &api.Preferences{}, nil
		}
		return nil, err
	}
	preference := &api.Preferences{}
	err = protojson.Unmarshal(json, preference)
	if err != nil {
		return &api.Preferences{}, nil
	}
	return preference, nil
}

func (d *DotLocal) savePreferences() error {
	mappings := lo.Map(d.GetMappings(), func(mapping internal.Mapping, _ int) *api.Mapping {
		return &api.Mapping{
			Id:         &mapping.ID,
			Host:       &mapping.Host,
			PathPrefix: &mapping.PathPrefix,
			Target:     &mapping.Target,
			ExpiresAt:  &timestamppb.Timestamp{Seconds: mapping.ExpresAt.Unix()},
		}
	})
	preference := &api.Preferences{
		Mappings: mappings,
	}
	json, err := protojson.Marshal(preference)
	if err != nil {
		return err
	}
	err = os.WriteFile(util.GetPreferencesPath(), json, 0644)
	if err != nil {
		return err
	}
	d.logger.Info("Saved preferences", zap.String("path", util.GetPreferencesPath()))
	return nil
}

func (d *DotLocal) GetMappings() []internal.Mapping {
	return lo.MapToSlice(d.mappings, func(key internal.MappingKey, state *internal.MappingState) internal.Mapping {
		return internal.NewMapping(key, state)
	})
}

func (d *DotLocal) CreateMapping(opts internal.MappingOptions) (internal.Mapping, error) {
	if !strings.HasSuffix(opts.Host, ".local") {
		opts.Host += ".local"
	}
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
	d.logger.Info("Updated mappings", zap.Any("mappings", mappings))
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
