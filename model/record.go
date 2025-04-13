package model

import "encoding/json"

type Record struct {
	ContainerID   string   `json:"containerId" yaml:"containerId" xml:"container_id"`
	ContainerName string   `json:"containerName" yaml:"containerName" xml:"container_name"`
	Ts            int64    `json:"timestamp" yaml:"timestamp" xml:"timestamp"`
	PidCount      *uint64  `json:"pidCount,omitempty" yaml:"pidCount,omitempty" xml:"pid_count,omitempty"`
	PidLimit      *uint64  `json:"pidLimit,omitempty" yaml:"pidLimit,omitempty" xml:"pid_limit,omitempty"`
	MemoryUsage   *uint64  `json:"memoryUsage,omitempty" yaml:"memoryUsage,omitempty" xml:"memory_usage,omitempty"`
	MemoryLimit   *uint64  `json:"memoryLimit,omitempty" yaml:"memoryLimit,omitempty" xml:"memory_limit,omitempty"`
	CPUPercent    float64  `json:"cpuPercent" yaml:"cpuPercent" xml:"cpu_percent"`
	CPUCount      uint64   `json:"cpuCount" yaml:"cpuCount" xml:"cpu_count"`
	CPUEnergy     *float64 `json:"cpuEnergy,omitempty" yaml:"cpuEnergy,omitempty" xml:"cpu_energy,omitempty"`
}

// MarshalJSON ensures proper handling of null values
func (r Record) MarshalJSON() ([]byte, error) {
	type Alias Record
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&r),
	})
}
