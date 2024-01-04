package dnsproxy

type DNSProxy interface {
	Start(port int) error
	SetHosts(hosts map[string]struct{}) error
	Stop() error
}
