package internal

import "time"

type Mapping struct {
	Host       string
	PathPrefix string
	Target     string
}
type MappingState struct {
	ID        string
	ExpiresAt time.Time
}

type MappingOptions struct {
	Host       string
	PathPrefix string
	Target     string
}
