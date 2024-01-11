package dnsproxy

type DNSProxy interface {
	Start() error
	SetHosts(hosts map[string]struct{}) error
	Stop() error
}
