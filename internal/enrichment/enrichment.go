package enrichment

import "github.com/trinhminhtriet/docker-activity/internal/model"

// Enricher enriches a Record.
type Enricher interface {
	Enrich(model.Record) model.Record
	Reset()
}

// Builder builds enrichers.
type Builder struct{}

// Build creates enrichers.
func (b *Builder) Build() []Enricher {
	return nil // Stub
}
