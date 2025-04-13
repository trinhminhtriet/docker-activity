package enrichment

import "github.com/trinhminhtriet/docker-activity/model"

// Enricher defines the interface for record enrichment
type Enricher interface {
	Enrich(record model.Record) model.Record
	Reset()
}

// BaseEnricher provides a default no-op implementation
type BaseEnricher struct{}

func (e *BaseEnricher) Enrich(record model.Record) model.Record {
	return record
}

func (e *BaseEnricher) Reset() {}
