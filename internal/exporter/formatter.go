package exporter

import (
	"encoding/json"

	"github.com/trinhminhtriet/docker-activity/internal/model"
)

type formatter interface {
	format(model.Record) (string, error)
}

type jsonFormatter struct{}

func (f *jsonFormatter) format(record model.Record) (string, error) {
	data, err := json.Marshal(record)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func newFormatter(format string) formatter {
	// Only JSON supported for now
	return &jsonFormatter{}
}
