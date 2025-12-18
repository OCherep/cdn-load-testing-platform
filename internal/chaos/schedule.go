package chaos

import "time"

type Schedule struct {
	Enabled bool `json:"enabled"`

	Windows []Window `json:"windows"`
}

type Window struct {
	StartOffsetSec int64  `json:"start_offset_sec"`
	DurationSec    int64  `json:"duration_sec"`
	Type           string `json:"type"` // latency | error | pause

	Config Config `json:"config"`
}

func (s Schedule) Active(now time.Time, testStart time.Time) *Config {
	if !s.Enabled {
		return nil
	}

	elapsed := int64(now.Sub(testStart).Seconds())

	for _, w := range s.Windows {
		if elapsed >= w.StartOffsetSec &&
			elapsed <= w.StartOffsetSec+w.DurationSec {
			return &w.Config
		}
	}

	return nil
}
