package formatter

import (
	"strconv"

	"github.com/trinhminhtriet/docker-activity/model"
)

type CSVFormatter struct{}

func (f *CSVFormatter) Format(record model.Record) (string, error) {
	base := formatCSVBase(record)
	if record.CPUEnergy != nil {
		return base + "," + formatFloat(*record.CPUEnergy), nil
	}
	return base, nil
}

func formatCSVBase(record model.Record) string {
	return joinStrings(
		record.ContainerID,
		record.ContainerName,
		formatInt(record.Ts),
		formatUintPtr(record.PidCount),
		formatUintPtr(record.PidLimit),
		formatUintPtr(record.MemoryUsage),
		formatUintPtr(record.MemoryLimit),
		formatFloat(record.CPUPercent),
		formatUint(record.CPUCount),
	)
}

func joinStrings(values ...string) string {
	result := ""
	for i, v := range values {
		if i > 0 {
			result += ","
		}
		result += v
	}
	return result
}

func formatUint(v uint64) string {
	return strconv.FormatUint(v, 10)
}
