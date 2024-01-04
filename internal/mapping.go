package internal

type Mapping struct {
	ID         string
	Host       string
	PathPrefix string
	Target     string
}

type MappingOptions struct {
	Host       string
	PathPrefix string
	Target     string
}
