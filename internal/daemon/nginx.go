package daemon

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/samber/lo"
	"github.com/softnetics/dotlocal/internal"
	"github.com/softnetics/dotlocal/internal/util"
	"github.com/tufanbarisyildirim/gonginx"
	"github.com/tufanbarisyildirim/gonginx/parser"
	"go.uber.org/zap"
)

type Nginx struct {
	logger     *zap.Logger
	configFile string
	port       int
	cmd        *exec.Cmd
	mappings   []internal.Mapping
}

func NewNginx(logger *zap.Logger) (*Nginx, error) {
	configFile, err := util.CreateTmpFile()
	if err != nil {
		return nil, err
	}
	port, err := findAvailablePort()
	if err != nil {
		return nil, err
	}
	return &Nginx{
		logger:     logger,
		configFile: configFile,
		port:       port,
		cmd:        nil,
		mappings:   nil,
	}, nil
}

func (n *Nginx) Start() error {
	n.writeConfig()
	n.logger.Debug("Starting nginx", zap.Int("port", n.port))

	fmt.Printf("nginx -t -c %s\n", n.configFile)
	cmd := exec.Command("nginx", "-c", n.configFile)
	err := cmd.Start()
	if err != nil {
		return err
	}
	n.cmd = cmd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	time.Sleep(100)

	n.logger.Info("Ready")

	return nil
}

func (n *Nginx) SetMappings(mappings []internal.Mapping) error {
	n.mappings = mappings
	err := n.writeConfig()
	if err != nil {
		return err
	}
	return n.reloadConfig()
}

func (n *Nginx) Stop() error {
	n.logger.Info("Stopping nginx")
	err := n.cmd.Process.Signal(syscall.SIGTERM)
	if err != nil {
		return err
	}
	return n.cmd.Wait()
}

func (n *Nginx) Port() int {
	return n.port
}

func (n *Nginx) writeConfig() error {
	p := parser.NewStringParser(`
		daemon off;
		events {
			worker_connections 1024;
		}
		http {
		}
	`)

	conf, err := p.Parse()
	if err != nil {
		return err
	}

	http := conf.FindDirectives("http")[0].GetBlock().(*gonginx.HTTP)
	hosts := lo.GroupBy(n.mappings, func(m internal.Mapping) string {
		return m.Host
	})
	for host, hostMappings := range hosts {
		var locations []gonginx.IDirective
		for _, mapping := range hostMappings {
			locations = append(locations, &gonginx.Directive{
				Name:       "location",
				Parameters: []string{mapping.PathPrefix},
				Block: &gonginx.LuaBlock{
					Directives: []gonginx.IDirective{
						&gonginx.Directive{
							Name:       "proxy_pass",
							Parameters: []string{mapping.Target},
						},
					},
				},
			})
		}
		http.Directives = append(http.Directives, &gonginx.Directive{
			Name:       "server",
			Parameters: []string{},
			Block: &gonginx.LuaBlock{
				Directives: append([]gonginx.IDirective{
					&gonginx.Directive{
						Name:       "listen",
						Parameters: []string{fmt.Sprintf("%d", n.port)},
					},
					&gonginx.Directive{
						Name:       "server_name",
						Parameters: []string{host},
					},
				}, locations...),
			},
		})
	}

	configString := gonginx.DumpBlock(conf, gonginx.IndentedStyle)
	// println(configString)
	n.logger.Debug("Writing nginx config", zap.String("config", configString))
	err = os.WriteFile(n.configFile, []byte(configString), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (n *Nginx) reloadConfig() error {
	err := n.cmd.Process.Signal(syscall.SIGHUP)
	if err != nil {
		return err
	}
	n.logger.Info("Reloaded nginx config")
	return nil
}

func findAvailablePort() (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	port := listener.Addr().(*net.TCPAddr).Port
	err = listener.Close()
	if err != nil {
		return 0, err
	}
	return port, nil
}
