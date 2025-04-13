package model

import (
	"encoding/json"
	"strings"

	"github.com/docker/docker/api/types"
)

// Record holds container statistics.
type Record struct {
	ContainerID   string  `json:"containerId"`
	ContainerName string  `json:"containerName"`
	Timestamp     int64   `json:"ts"`
	PIDCount      uint64  `json:"pidCount"`
	PIDLimit      uint64  `json:"pidLimit"`
	MemoryUsage   uint64  `json:"memoryUsage"`
	MemoryLimit   uint64  `json:"memoryLimit"`
	CPUPercent    float64 `json:"cpuPercent"`
	CPUCount      uint64  `json:"cpuCount"`
}

// NewRecordFromStats converts Docker stats to a Record.
func NewRecordFromStats(stats types.StatsJSON) Record {
	var cpuDelta, systemDelta uint64
	if stats.CPUStats.CPUUsage.TotalUsage > stats.PreCPUStats.CPUUsage.TotalUsage {
		cpuDelta = stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage
	}
	if stats.CPUStats.SystemCPUUsage > stats.PreCPUStats.SystemCPUUsage {
		systemDelta = stats.CPUStats.SystemCPUUsage - stats.PreCPUStats.SystemCPUUsage
	}
	cpuCount := uint64(1)
	if stats.CPUStats.OnlineCPUs > 0 {
		cpuCount = uint64(stats.CPUStats.OnlineCPUs)
	}
	var cpuPercent float64
	if systemDelta > 0 {
		cpuPercent = float64(cpuDelta) / float64(systemDelta)
	}

	return Record{
		ContainerID:   stats.ID,
		ContainerName: strings.TrimPrefix(stats.Name, "/"),
		Timestamp:     stats.Read.UTC().Unix(),
		PIDCount:      stats.PidsStats.Current,
		PIDLimit:      stats.PidsStats.Limit,
		MemoryUsage:   stats.MemoryStats.Usage,
		MemoryLimit:   stats.MemoryStats.Limit,
		CPUPercent:    cpuPercent,
		CPUCount:      cpuCount,
	}
}

// MarshalJSON customizes JSON output, omitting zero values.
func (r Record) MarshalJSON() ([]byte, error) {
	type Alias Record
	aux := struct {
		Alias
		PIDCount    *uint64 `json:"pidCount,omitempty"`
		PIDLimit    *uint64 `json:"pidLimit,omitempty"`
		MemoryUsage *uint64 `json:"memoryUsage,omitempty"`
		MemoryLimit *uint64 `json:"memoryLimit,omitempty"`
	}{
		Alias: Alias(r),
	}
	if r.PIDCount > 0 {
		aux.PIDCount = &r.PIDCount
	}
	if r.PIDLimit > 0 {
		aux.PIDLimit = &r.PIDLimit
	}
	if r.MemoryUsage > 0 {
		aux.MemoryUsage = &r.MemoryUsage
	}
	if r.MemoryLimit > 0 {
		aux.MemoryLimit = &r.MemoryLimit
	}
	return json.Marshal(aux)
}
