package formatter

import (
	"encoding/xml"
	"fmt"

	"github.com/trinhminhtriet/docker-activity/model"
)

type XMLFormatter struct{}

type xmlRecord struct {
	XMLName       xml.Name `xml:"record"`
	ContainerID   string   `xml:"container_id"`
	ContainerName string   `xml:"container_name"`
	Timestamp     int64    `xml:"timestamp"`
	CPU           cpu      `xml:"cpu"`
	Memory        memory   `xml:"memory"`
	Processes     procs    `xml:"processes"`
}

type cpu struct {
	Percent float64  `xml:"percent"`
	Cores   uint64   `xml:"cores"`
	Energy  *float64 `xml:"energy,omitempty"`
}

type memory struct {
	Usage *uint64 `xml:"usage,omitempty"`
	Limit *uint64 `xml:"limit,omitempty"`
}

type procs struct {
	Count *uint64 `xml:"count,omitempty"`
	Limit *uint64 `xml:"limit,omitempty"`
}

func (f *XMLFormatter) Format(record model.Record) (string, error) {
	xr := xmlRecord{
		ContainerID:   record.ContainerID,
		ContainerName: record.ContainerName,
		Timestamp:     record.Ts,
		CPU: cpu{
			Percent: record.CPUPercent,
			Cores:   record.CPUCount,
			Energy:  record.CPUEnergy,
		},
		Memory: memory{
			Usage: record.MemoryUsage,
			Limit: record.MemoryLimit,
		},
		Processes: procs{
			Count: record.PidCount,
			Limit: record.PidLimit,
		},
	}

	data, err := xml.MarshalIndent(xr, "", "  ")
	if err != nil {
		return "", fmt.Errorf("xml marshal error: %w", err)
	}
	return string(data), nil
}
