package daemon

import (
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
	mappings []internal.Mapping
}

func NewDotLocal(logger *zap.Logger) (*DotLocal, error) {
	nginx, err := NewNginx(logger.Named("nginx"))
	if err != nil {
		return nil, err
	}

	dnsProxy, err := orbdnsproxy.NewOrbstackDNSProxy(logger.Named("orbdnsproxxy"))
	if err != nil {
		return nil, err
	}

	return &DotLocal{
		logger:   logger,
		nginx:    nginx,
		dnsProxy: dnsProxy,
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

func (d *DotLocal) SetMappings(mappings []internal.Mapping) error {
	d.mappings = mappings
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

func (d *DotLocal) Wait() error {
	var t tomb.Tomb
	t.Go(func() error {
		return d.nginx.Wait()
	})
	return t.Wait()
}
