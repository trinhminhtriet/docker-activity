package exporter

import (
	"os"

	"github.com/spf13/pflag"
	"github.com/trinhminhtriet/docker-activity/internal/model"
)

// StdoutConfig configures stdout-based export.
type StdoutConfig struct {
	Format string
}

// Flags returns the flag set for StdoutConfig.
func (c *StdoutConfig) Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("stdout", pflag.ContinueOnError)
	fs.StringVar(&c.Format, "format", "json", "Format of the output records (json)")
	return fs
}

// Exporter creates a stdout exporter.
func (c *StdoutConfig) Exporter() Exporter {
	return &stdoutExporter{
		writer:    os.Stdout,
		formatter: newFormatter(c.Format),
	}
}

type stdoutExporter struct {
	writer    *os.File
	formatter formatter
}

func (e *stdoutExporter) Handle(record model.Record) error {
	line, err := e.formatter.format(record)
	if err != nil {
		return err
	}
	_, err = e.writer.WriteString(line + "\n")
	return err
}
