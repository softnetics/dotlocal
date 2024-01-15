package dnsproxy

import "context"

type DNSProxy interface {
	Start(ctx context.Context) error
	SetHosts(hosts map[string]struct{}) error
	Stop() error
}
