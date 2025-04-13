package formatter

import (
	"encoding/json"
	"fmt"

	"github.com/trinhminhtriet/docker-activity/model"
)

type JSONFormatter struct{}

func (f *JSONFormatter) Format(record model.Record) (string, error) {
	data, err := json.Marshal(record)
	if err != nil {
		return "", fmt.Errorf("json marshal error: %w", err)
	}
	return string(data), nil
}
