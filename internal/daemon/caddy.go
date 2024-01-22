package daemon

import (
	"crypto/x509"
	"encoding/json"
	"errors"
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
	"golang.org/x/exp/slices"
)

type Caddy struct {
	logger   *zap.Logger
	mappings []internal.Mapping
	prefs    *preferences
}

func NewCaddy(logger *zap.Logger, prefs *preferences) (*Caddy, error) {
	return &Caddy{
		logger:   logger,
		mappings: nil,
		prefs:    prefs,
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

func (c *Caddy) reload() error {
	cfgJSON := encodeJson(c.config())
	err := caddy.Load(cfgJSON, false)
	if err != nil {
		return err
	}
	return nil
}

func (c *Caddy) SetMappings(mappings []internal.Mapping) error {
	c.mappings = mappings
	c.logger.Debug("Setting mappings", zap.Any("mappings", mappings))
	return c.reload()
}

func (c *Caddy) getRootCertificate() (*x509.Certificate, error) {
	caddyCtx := caddy.ActiveContext()
	pki, ok := caddyCtx.AppIfConfigured("pki").(*caddypki.PKI)
	if !ok {
		return nil, errors.New("pki module not found")
	}
	localCA := pki.CAs["local"]
	if localCA == nil {
		return nil, errors.New("local CA not found")
	}
	rootCert := localCA.RootCertificate()
	if rootCert == nil {
		return nil, errors.New("root certificate not found")
	}
	return rootCert, nil
}

func (c *Caddy) config() *caddy.Config {
	boolFalse := false
	routes := c.routes()

	servers := make(map[string]*caddyhttp.Server)
	if !c.prefs.disableHTTPS {
		servers["https"] = &caddyhttp.Server{
			Listen: []string{"127.0.0.1:443"},
			Routes: routes,
			AutoHTTPS: &caddyhttp.AutoHTTPSConfig{
				DisableRedir: !c.prefs.redirectHTTPS,
			},
		}
	}
	if c.prefs.disableHTTPS || !c.prefs.redirectHTTPS {
		servers["http"] = &caddyhttp.Server{
			Listen: []string{"127.0.0.1:80"},
			Routes: routes,
		}
	}

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
				Servers:   servers,
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
