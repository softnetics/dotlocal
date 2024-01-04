package internal

import "time"

type Mapping struct {
	ID         string
	Host       string
	PathPrefix string
	Target     string
}

func NewMapping(key MappingKey, state *MappingState) Mapping {
	return Mapping{
		ID:         state.ID,
		Host:       key.Host,
		PathPrefix: key.PathPrefix,
		Target:     state.Target,
	}
}

type MappingKey struct {
	Host       string
	PathPrefix string
}

type MappingState struct {
	ID        string
	Target    string
	ExpiresAt time.Time
}

type MappingOptions struct {
	Host       string
	PathPrefix string
	Target     string
}
