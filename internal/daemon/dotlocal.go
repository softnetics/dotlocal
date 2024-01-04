package daemon

import (
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
	mappings map[string]internal.Mapping
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
		mappings: make(map[string]internal.Mapping),
	}, nil
}

func (d *DotLocal) Start() error {
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

	return nil
}

func (d *DotLocal) GetMappings() []internal.Mapping {
	return lo.MapToSlice(d.mappings, func(host string, mapping internal.Mapping) internal.Mapping {
		return mapping
	})
}

func (d *DotLocal) CreateMapping(opts internal.MappingOptions) (internal.Mapping, error) {
	id := uniuri.NewLen(6)
	if opts.PathPrefix == "" {
		opts.PathPrefix = "/"
	}
	mapping := internal.Mapping{
		ID:         id,
		Host:       opts.Host,
		PathPrefix: opts.PathPrefix,
		Target:     opts.Target,
	}
	d.mappings[id] = mapping
	d.logger.Info("Created mapping", zap.Any("mapping", mapping))
	return mapping, d.UpdateMappings()
}

func (d *DotLocal) RemoveMapping(ids ...string) error {
	for _, id := range ids {
		delete(d.mappings, id)
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

func (d *DotLocal) Stop() error {
	var t tomb.Tomb
	t.Go(func() error {
		return d.nginx.Stop()
	})
	return t.Wait()
}
