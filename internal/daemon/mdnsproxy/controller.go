package mdnsproxy

import (
	"context"

	dnssd "github.com/softnetics/dotlocal/dns-sd"
	"github.com/softnetics/dotlocal/internal/daemon/dnsproxy"
	"go.uber.org/zap"
)

var nginxImage = "nginx:1.24.0-alpine"

type MDNSProxy struct {
	logger          *zap.Logger
	dnsService      dnssd.DNSService
	registeredHosts map[string]dnssd.DNSRecord

	cancelProcess context.CancelFunc
}

func NewMDNSProxy(logger *zap.Logger) (dnsproxy.DNSProxy, error) {
	return &MDNSProxy{
		logger:          logger,
		registeredHosts: make(map[string]dnssd.DNSRecord),
	}, nil
}

func (p *MDNSProxy) Start(ctx context.Context) error {
	p.logger.Debug("Connecting to dns service")
	service, err := dnssd.NewConnection()
	if err != nil {
		return err
	}
	p.dnsService = service
	p.logger.Info("Ready")

	ctx, cancel := context.WithCancel(ctx)
	p.cancelProcess = cancel
	go func() {
		err := service.Process(ctx)
		if err != nil {
			p.logger.Error("Failed to process dns service", zap.Error(err))
		}
	}()

	return nil
}

func (p *MDNSProxy) Stop() error {
	p.logger.Info("Stopping")
	p.cancelProcess()
	p.dnsService.Deallocate()
	return nil
}

func (p *MDNSProxy) SetHosts(hostsMap map[string]struct{}) error {
	p.logger.Debug("Setting hosts", zap.Any("hosts", hostsMap))

	for host := range hostsMap {
		if _, ok := p.registeredHosts[host]; ok {
			continue
		}
		p.logger.Debug("Adding host", zap.String("host", host))
		record, err := p.dnsService.RegisterProxyAddressRecord(host, "127.0.0.1", 0)
		if err != nil {
			return err
		}
		p.registeredHosts[host] = record
	}

	for host, record := range p.registeredHosts {
		if _, ok := hostsMap[host]; !ok {
			p.logger.Debug("Removing host", zap.String("host", host))
			err := p.dnsService.RemoveRecord(record, 0)
			if err != nil {
				return err
			}
			delete(p.registeredHosts, host)
		}
	}
	return nil
}
