package daemon

import (
	"encoding/json"
	"net"
	"net/url"
	"path"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp/reverseproxy"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp/rewrite"
	"github.com/caddyserver/caddy/v2/modules/caddypki"
	"github.com/caddyserver/caddy/v2/modules/caddytls"
	_ "github.com/caddyserver/caddy/v2/modules/filestorage"
	"github.com/samber/lo"
	"github.com/softnetics/dotlocal/internal"
	"github.com/softnetics/dotlocal/internal/util"
	"go.uber.org/zap"
)

type Caddy struct {
	logger   *zap.Logger
	mappings []internal.Mapping
}

func NewCaddy(logger *zap.Logger) (*Caddy, error) {
	return &Caddy{
		logger:   logger,
		mappings: nil,
	}, nil
}

func (c *Caddy) Start() error {
	cfgJSON := encodeJson(c.config())
	err := caddy.Load(cfgJSON, true)
	if err != nil {
		return err
	}
	return nil
}

func (c *Caddy) SetMappings(mappings []internal.Mapping) error {
	c.mappings = mappings
	c.logger.Debug("Setting mappings", zap.Any("mappings", mappings))
	cfgJSON := encodeJson(c.config())
	err := caddy.Load(cfgJSON, true)
	if err != nil {
		return err
	}
	return nil
}

func (c *Caddy) config() *caddy.Config {
	boolFalse := false
	routes := c.routes()
	return &caddy.Config{
		Admin: &caddy.AdminConfig{
			Disabled: true,
			Config: &caddy.ConfigSettings{
				Persist: &boolFalse,
			},
		},
		Logging: &caddy.Logging{
			Logs: map[string]*caddy.CustomLog{
				"default": {
					BaseLog: caddy.BaseLog{
						Level: zap.DebugLevel.String(),
					},
				},
			},
		},
		StorageRaw: encodeJson(map[string]string{
			"module": "file_system",
			"root":   caddyStoragePath(),
		}),
		AppsRaw: caddy.ModuleMap{
			"pki": encodeJson(&caddypki.PKI{
				CAs: map[string]*caddypki.CA{
					"local": {
						Name:                   "Local",
						InstallTrust:           &boolFalse,
						RootCommonName:         "DotLocal Local Authority",
						IntermediateCommonName: "DotLocal Local Intermediate Authority",
					},
				},
			}),
			"tls": encodeJson(&caddytls.TLS{
				Automation: &caddytls.AutomationConfig{
					Policies: []*caddytls.AutomationPolicy{
						{
							IssuersRaw: []json.RawMessage{
								encodeJson(map[string]string{
									"module": "internal",
								}),
							},
							OnDemand: true,
						},
					},
				},
			}),
			"http": encodeJson(&caddyhttp.App{
				HTTPPort:  80,
				HTTPSPort: 443,
				Servers: map[string]*caddyhttp.Server{
					"srv0": {
						Listen: []string{"127.0.0.1:443"},
						Routes: routes,
						AutoHTTPS: &caddyhttp.AutoHTTPSConfig{
							DisableRedir: true,
						},
					},
					"srv1": {
						Listen: []string{"127.0.0.1:80"},
						Routes: routes,
					},
				},
			}),
		},
	}
}

func (c *Caddy) routes() []caddyhttp.Route {
	hosts := lo.GroupBy(c.mappings, func(m internal.Mapping) string {
		return m.Host
	})
	hostRoutes := lo.MapToSlice(hosts, func(host string, mappings []internal.Mapping) caddyhttp.Route {
		var routes caddyhttp.RouteList
		hasRoot := false
		slices.SortFunc(mappings, func(i, j internal.Mapping) bool {
			return len(i.PathPrefix) > len(j.PathPrefix)
		})
		for _, mapping := range mappings {
			targetUri, err := url.Parse(mapping.Target)
			if err != nil {
				panic(err)
			}
			targetHostPort := targetUri.Host
			if !strings.Contains(targetHostPort, ":") {
				targetHostPort += ":"
			}
			targetHost, targetPort, err := net.SplitHostPort(targetHostPort)
			if err != nil {
				panic(err)
			}
			if targetPort == "" {
				if targetUri.Scheme == "https" {
					targetPort = "443"
				} else {
					targetPort = "80"
				}
			}
			handlerRewrite := &rewrite.Rewrite{}
			if targetUri.Path != "" && targetUri.Path != "/" {
				handlerRewrite.URI = targetUri.Path + "{http.request.uri.path}"
			}

			handler := encodeJson(map[string]any{
				"handler": "reverse_proxy",
				"upstreams": reverseproxy.UpstreamPool{
					{Dial: targetHost + ":" + targetPort},
				},
				"rewrite": handlerRewrite,
			})
			if mapping.PathPrefix == "/" {
				hasRoot = true
				routes = append(routes, caddyhttp.Route{
					Group:       host,
					HandlersRaw: []json.RawMessage{handler},
					MatcherSetsRaw: []caddy.ModuleMap{
						{
							"path": encodeJson([]string{"/*"}),
						},
					},
				})
			} else {
				routes = append(routes, caddyhttp.Route{
					Group:       host,
					HandlersRaw: []json.RawMessage{handler},
					MatcherSetsRaw: []caddy.ModuleMap{
						{
							"path": encodeJson([]string{mapping.PathPrefix}),
						},
					},
				})
				routes = append(routes, caddyhttp.Route{
					Group:       host,
					HandlersRaw: []json.RawMessage{handler},
					MatcherSetsRaw: []caddy.ModuleMap{
						{
							"path": encodeJson([]string{mapping.PathPrefix + "/*"}),
						},
					},
				})
			}
		}
		if !hasRoot {
			routes = append(routes, caddyhttp.Route{
				Group: host,
				HandlersRaw: []json.RawMessage{
					encodeJson(map[string]any{
						"handler":     "static_response",
						"status_code": 404,
					}),
				},
				MatcherSetsRaw: []caddy.ModuleMap{
					{
						"path": encodeJson([]string{"/*"}),
					},
				},
			})
		}
		return caddyhttp.Route{
			HandlersRaw: []json.RawMessage{
				encodeJson(map[string]any{
					"handler": "subroute",
					"routes":  routes,
				}),
			},
			MatcherSetsRaw: []caddy.ModuleMap{
				{
					"host": encodeJson([]string{host}),
				},
			},
			Terminal: true,
		}
	})
	return append(hostRoutes, caddyhttp.Route{
		HandlersRaw: []json.RawMessage{
			encodeJson(map[string]any{
				"handler":     "static_response",
				"status_code": 404,
			}),
		},
	})
}

func encodeJson(value any) json.RawMessage {
	return lo.Must1(json.Marshal(value))
}

func caddyStoragePath() string {
	return path.Join(util.GetDotlocalPath(), "caddy")
}
