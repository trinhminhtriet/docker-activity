package exporter

import (
	"github.com/spf13/pflag"
	"github.com/trinhminhtriet/docker-activity/internal/model"
)

// Config defines an exporter configuration.
type Config interface {
	Flags() *pflag.FlagSet
	Exporter() Exporter
}

// Exporter handles Record exporting.
type Exporter interface {
	Handle(model.Record) error
}
