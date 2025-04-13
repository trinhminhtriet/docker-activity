package formatter

import (
	"fmt"
	"time"

	"github.com/trinhminhtriet/docker-activity/model"
)

type TextFormatter struct{}

func (f *TextFormatter) Format(record model.Record) (string, error) {
	ts := time.Unix(record.Ts, 0).Format(time.RFC3339)
	return fmt.Sprintf(
		"Container: %s (%s)\n"+
			"Time: %s\n"+
			"CPU: %.2f%% of %d cores\n"+
			"Memory: %s/%s\n"+
			"PIDs: %s/%s\n"+
			"Energy: %s\n",
		record.ContainerName, record.ContainerID,
		ts,
		record.CPUPercent, record.CPUCount,
		formatBytes(record.MemoryUsage), formatBytes(record.MemoryLimit),
		formatUintPtr(record.PidCount), formatUintPtr(record.PidLimit),
		formatEnergy(record.CPUEnergy),
	), nil
}
