package orbdnsproxy

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/softnetics/dotlocal/internal/daemon/dnsproxy"
	"github.com/softnetics/dotlocal/internal/util"
	"github.com/tufanbarisyildirim/gonginx"
	"go.uber.org/zap"
	"gopkg.in/tomb.v2"
)

var nginxImage = "nginx:1.24.0-alpine"

type OrbstackDNSProxy struct {
	logger          *zap.Logger
	docker          *client.Client
	port            int
	containers      map[string]*Container
	nginxConfigFile string
}

func NewOrbstackDNSProxy(logger *zap.Logger) (dnsproxy.DNSProxy, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	nginxConfigFile, err := util.CreateTmpFile()
	if err != nil {
		return nil, err
	}

	return &OrbstackDNSProxy{
		logger:          logger,
		docker:          cli,
		containers:      make(map[string]*Container),
		nginxConfigFile: nginxConfigFile,
	}, nil
}

func (p *OrbstackDNSProxy) Start(port int) error {
	p.port = port
	p.logger.Debug("Ensuring nginx image exists", zap.String("image", nginxImage))
	err := p.writeNginxConfig()
	if err != nil {
		return err
	}
	err = ensureImageExists(p.docker, nginxImage)
	if err != nil {
		return err
	}

	p.logger.Debug("Cleaning up existing containers")
	containers, err := p.docker.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(filters.Arg("label", "managed-dotlocal")),
	})
	if err != nil {
		return err
	}
	for _, container := range containers {
		p.logger.Debug("Removing container", zap.String("id", container.ID))
		err := p.docker.ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{
			Force: true,
		})
		if err != nil {
			return err
		}
	}

	p.logger.Info("Ready")

	return nil
}

func (p *OrbstackDNSProxy) Stop() error {
	p.logger.Info("Stopping")
	var t tomb.Tomb
	for _, container := range p.containers {
		t.Go(func() error {
			return container.Remove()
		})
	}
	t.Go(func() error {
		return os.Remove(p.nginxConfigFile)
	})
	return t.Wait()
}

func (p *OrbstackDNSProxy) SetHosts(hosts map[string]struct{}) error {
	p.logger.Debug("Setting hosts", zap.Any("hosts", hosts))

	var t tomb.Tomb

	for host := range hosts {
		_, exists := p.containers[host]
		if exists {
			continue
		}

		container, err := NewContainer(p.logger, p.docker, p.nginxConfigFile, host)
		if err != nil {
			return err
		}
		p.containers[host] = container
		t.Go(func() error {
			return container.CreateAndStart()
		})
	}

	for _host, _container := range p.containers {
		host := _host
		container := _container
		_, exists := hosts[host]
		if exists {
			continue
		}
		t.Go(func() error {
			delete(p.containers, host)
			return container.Remove()
		})
	}

	return t.Wait()
}

func (p *OrbstackDNSProxy) writeNginxConfig() error {
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
										Name:       "proxy_pass",
										Parameters: []string{fmt.Sprintf("http://host.docker.internal:%d", p.port)},
									},
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

func ensureImageExists(cli *client.Client, containerID string) error {
	_, _, err := cli.ImageInspectWithRaw(context.Background(), containerID)
	if err == nil {
		return nil
	}
	out, err := cli.ImagePull(context.Background(), containerID, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, out)

	defer out.Close()

	return nil
}
