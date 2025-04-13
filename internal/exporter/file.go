package exporter

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/trinhminhtriet/docker-activity/internal/model"
)

// FileConfig configures file-based export.
type FileConfig struct {
	Format string
	Path   string
}

// Flags returns the flag set for FileConfig.
func (c *FileConfig) Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("file", pflag.ContinueOnError)
	fs.StringVar(&c.Format, "file-format", "json", "Format of the output records (json)")
	fs.StringVar(&c.Path, "file-output", "", "Path to write the file")
	return fs
}

// Exporter creates a file exporter.
func (c *FileConfig) Exporter() Exporter {
	file, err := os.OpenFile(c.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		panic(fmt.Sprintf("couldn't open output file: %v", err))
	}
	return &fileExporter{
		file:      file,
		formatter: newFormatter(c.Format),
	}
}

type fileExporter struct {
	file      *os.File
	formatter formatter
}

func (e *fileExporter) Handle(record model.Record) error {
	line, err := e.formatter.format(record)
	if err != nil {
		return err
	}
	_, err = e.file.WriteString(line + "\n")
	return err
}
