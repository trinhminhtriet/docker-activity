package enrichment

import "github.com/trinhminhtriet/docker-activity/model"

// EnrichmentStack manages multiple enrichers in sequence
type EnrichmentStack struct {
	enrichers []Enricher
}

// NewEnrichmentStack creates a new stack with the given enrichers
func NewEnrichmentStack(enrichers ...Enricher) *EnrichmentStack {
	return &EnrichmentStack{
		enrichers: enrichers,
	}
}

// Enrich applies all enrichers in sequence
func (s *EnrichmentStack) Enrich(record model.Record) model.Record {
	for _, enricher := range s.enrichers {
		record = enricher.Enrich(record)
	}
	return record
}

// Reset resets all enrichers in the stack
func (s *EnrichmentStack) Reset() {
	for _, enricher := range s.enrichers {
		enricher.Reset()
	}
}

// Add appends an enricher to the stack
func (s *EnrichmentStack) Add(enricher Enricher) {
	s.enrichers = append(s.enrichers, enricher)
}
