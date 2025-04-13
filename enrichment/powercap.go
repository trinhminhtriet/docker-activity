package enrichment

import (
	"sync"
	"time"

	"github.com/trinhminhtriet/docker-activity/model"
	"github.com/trinhminhtriet/docker-activity/pkg/powercap"
)

// PowerCapEnricher adds power consumption metrics to records
type PowerCapEnricher struct {
	reader    powercap.Reader
	lastValue *uint64
	lastTime  time.Time
	mu        sync.Mutex
}

// NewPowerCapEnricher creates a new powercap enricher
func NewPowerCapEnricher(reader powercap.Reader) *PowerCapEnricher {
	return &PowerCapEnricher{
		reader:    reader,
		lastValue: nil,
		lastTime:  time.Time{},
	}
}

// Enrich adds power consumption data to the record
func (e *PowerCapEnricher) Enrich(record model.Record) model.Record {
	e.mu.Lock()
	defer e.mu.Unlock()

	current, err := e.reader.Read()
	if err != nil {
		return record
	}

	if e.lastValue != nil {
		timeDelta := time.Since(e.lastTime).Seconds()
		if timeDelta > 0 {
			energyDelta := float64(current - *e.lastValue)
			power := energyDelta / timeDelta
			record.CPUEnergy = &power
		}
	}

	e.lastValue = &current
	e.lastTime = time.Now()

	return record
}

// Reset clears the powercap enricher state
func (e *PowerCapEnricher) Reset() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.lastValue = nil
	e.lastTime = time.Time{}
}
