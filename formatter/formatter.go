package formatter

import (
	"fmt"
	"strconv"

	"github.com/trinhminhtriet/docker-activity/model"
)

type Format string

const (
	FormatCSV  Format = "csv"
	FormatJSON Format = "json"
	FormatText Format = "text"
	FormatXML  Format = "xml"
	FormatYAML Format = "yaml"
)

func ParseFormat(s string) (Format, error) {
	switch s {
	case "csv":
		return FormatCSV, nil
	case "json":
		return FormatJSON, nil
	case "text":
		return FormatText, nil
	case "xml":
		return FormatXML, nil
	case "yaml":
		return FormatYAML, nil
	default:
		return FormatCSV, fmt.Errorf("unknown format: %s", s)
	}
}

func (f Format) Formatter() Formatter {
	switch f {
	case FormatJSON:
		return &JSONFormatter{}
	case FormatText:
		return &TextFormatter{}
	case FormatXML:
		return &XMLFormatter{}
	case FormatYAML:
		return &YAMLFormatter{}
	default: // Default to CSV
		return &CSVFormatter{}
	}
}

type Formatter interface {
	Format(record model.Record) (string, error)
}

// Helper functions
func formatUintPtr(v *uint64) string {
	if v == nil {
		return "0"
	}
	return strconv.FormatUint(*v, 10)
}

func formatInt(v int64) string {
	return strconv.FormatInt(v, 10)
}

func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', 2, 64)
}

func formatBytes(v *uint64) string {
	if v == nil {
		return "N/A"
	}
	const unit = 1024
	if *v < unit {
		return fmt.Sprintf("%d B", *v)
	}
	div, exp := uint64(unit), 0
	for n := *v / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(*v)/float64(div), "KMGTPE"[exp])
}

func formatEnergy(v *float64) string {
	if v == nil {
		return "N/A"
	}
	return fmt.Sprintf("%.2f J", *v)
}
