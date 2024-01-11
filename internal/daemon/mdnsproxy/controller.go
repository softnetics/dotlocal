package mdnsproxy

import (
	"os"
	"os/exec"

	"github.com/softnetics/dotlocal/internal/daemon/dnsproxy"
	"github.com/softnetics/dotlocal/internal/util"
	"github.com/tufanbarisyildirim/gonginx"
	"go.uber.org/zap"
	"gopkg.in/tomb.v2"
)

var nginxImage = "nginx:1.24.0-alpine"

type MDNSProxy struct {
	logger          *zap.Logger
	port            int
	nginxConfigFile string
	command         *exec.Cmd
}

func NewMDNSProxy(logger *zap.Logger) (dnsproxy.DNSProxy, error) {
	nginxConfigFile, err := util.CreateTmpFile()
	if err != nil {
		return nil, err
	}

	return &MDNSProxy{
		logger:          logger,
		nginxConfigFile: nginxConfigFile,
	}, nil
}

func (p *MDNSProxy) Start() error {
	p.logger.Debug("Ensuring nginx image exists", zap.String("image", nginxImage))
	err := p.writeNginxConfig()
	if err != nil {
		return err
	}
	p.logger.Info("Ready")
	return nil
}

func (p *MDNSProxy) Stop() error {
	p.logger.Info("Stopping")
	var t tomb.Tomb
	t.Go(func() error {
		return os.Remove(p.nginxConfigFile)
	})
	return t.Wait()
}

func (p *MDNSProxy) SetHosts(hostsMap map[string]struct{}) error {
	p.logger.Debug("Setting hosts", zap.Any("hosts", hostsMap))

	hosts := make([]string, len(hostsMap))
	i := 0
	for host := range hostsMap {
		hosts[i] = host
		i++
	}

	p.logger.Debug("Setting hosts", zap.Any("hosts", hosts))
	cmd := exec.Command("./c/build/dns-sd", hosts...)
	err := cmd.Start()
	if err != nil {
		return err
	}
	if p.command != nil {
		err := p.command.Process.Kill()
		if err != nil {
			return err
		}
	}
	p.command = cmd
	return nil
}

func (p *MDNSProxy) writeNginxConfig() error {
	conf := &gonginx.Block{
		Directives: []gonginx.IDirective{
			&gonginx.Directive{
				Name: "server",
				Block: &gonginx.Block{
					Directives: []gonginx.IDirective{
						&gonginx.Directive{
							Name:       "listen",
							Parameters: []string{"80"},
						},
						&gonginx.Directive{
							Name:       "location",
							Parameters: []string{"/"},
							Block: &gonginx.Block{
								Directives: []gonginx.IDirective{
									&gonginx.Directive{
										Name:       "proxy_http_version",
										Parameters: []string{"1.1"},
									},
									&gonginx.Directive{
										Name:       "proxy_set_header",
										Parameters: []string{"Upgrade", "$http_upgrade"},
									},
									&gonginx.Directive{
										Name:       "proxy_set_header",
										Parameters: []string{"Connection", "\"Upgrade\""},
									},
									&gonginx.Directive{
										Name:       "proxy_set_header",
										Parameters: []string{"Host", "$host"},
									},
									&gonginx.Directive{
										Name:       "proxy_set_header",
										Parameters: []string{"X-Forwarded-For", "$remote_addr"},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	configString := gonginx.DumpBlock(conf, gonginx.IndentedStyle)
	p.logger.Debug("Writing nginx config", zap.String("config", configString))
	err := os.WriteFile(p.nginxConfigFile, []byte(configString), 0644)
	if err != nil {
		return err
	}
	return nil
}
