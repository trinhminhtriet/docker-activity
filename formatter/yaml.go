package formatter

import (
	"fmt"

	"github.com/trinhminhtriet/docker-activity/model"
	"gopkg.in/yaml.v3"
)

type YAMLFormatter struct{}

func (f *YAMLFormatter) Format(record model.Record) (string, error) {
	data, err := yaml.Marshal(record)
	if err != nil {
		return "", fmt.Errorf("yaml marshal error: %w", err)
	}
	return string(data), nil
}
