package daemon

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"

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
	port, err := util.FindAvailablePort()
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

	fmt.Printf("nginx -c %s\n", n.configFile)
	cmd := exec.Command("nginx", "-c", n.configFile)
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = os.Stderr
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	nginxStarted := false
	go func() {
		func() {
			defer wg.Done()
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				line := scanner.Text()
				println(line)
				if strings.Contains(line, "start worker processes") {
					nginxStarted = true
					return
				}
			}
		}()
		io.Copy(os.Stdout, stdout)
	}()

	err = cmd.Start()
	if err != nil {
		return err
	}
	wg.Wait()
	if !nginxStarted {
		err := cmd.Wait()
		if err != nil {
			return err
		}
		return errors.New("nginx failed to start")
	}
	n.cmd = cmd

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
	n.cmd.Wait()
	return nil
}

func (n *Nginx) Port() int {
	return n.port
}

func (n *Nginx) writeConfig() error {
	p := parser.NewStringParser(`
		daemon off;
		error_log /dev/stdout info;
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
		directives := []gonginx.IDirective{
			&gonginx.Directive{
				Name:       "listen",
				Parameters: []string{strconv.Itoa(n.port)},
			},
			&gonginx.Directive{
				Name:       "server_name",
				Parameters: []string{host},
			},
		}
		locations := lo.SliceToMap(hostMappings, func(mapping internal.Mapping) (string, []gonginx.IDirective) {
			return mapping.PathPrefix, []gonginx.IDirective{
				&gonginx.Directive{
					Name:       "proxy_pass",
					Parameters: []string{mapping.Target},
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
			}
		})
		_, hasRoot := locations["/"]
		if !hasRoot {
			locations["/"] = []gonginx.IDirective{
				&gonginx.Directive{
					Name:       "return",
					Parameters: []string{"404"},
				},
			}
		}
		for pathPrefix, locationDirectives := range locations {
			directives = append(directives, &gonginx.Directive{
				Name:       "location",
				Parameters: []string{pathPrefix},
				Block: &gonginx.LuaBlock{
					Directives: locationDirectives,
				},
			})
		}
		http.Directives = append(http.Directives, &gonginx.Directive{
			Name:       "server",
			Parameters: []string{},
			Block: &gonginx.LuaBlock{
				Directives: directives,
			},
		})
	}
	http.Directives = append(http.Directives, &gonginx.Directive{
		Name:       "server",
		Parameters: []string{},
		Block: &gonginx.LuaBlock{
			Directives: []gonginx.IDirective{
				&gonginx.Directive{
					Name:       "listen",
					Parameters: []string{strconv.Itoa(n.port), "default_server"},
				},
				&gonginx.Directive{
					Name:       "return",
					Parameters: []string{"404"},
				},
			},
		},
	})

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
