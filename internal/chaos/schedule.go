package chaos

import "time"

type Stage struct {
	AfterSec   int  `json:"after_sec"`
	Enabled    bool `json:"enabled"`
	LatencyMs  int  `json:"latency_ms"`
	ErrorRate  int  `json:"error_rate"`
	BurstPause bool `json:"burst_pause"`
}

type Schedule struct {
	Stages []Stage `json:"stages"`
}

func (s Schedule) Current(start time.Time) Config {
	elapsed := int(time.Since(start).Seconds())

	var active Config
	for _, stage := range s.Stages {
		if elapsed >= stage.AfterSec {
			active = Config{
				Enabled:    stage.Enabled,
				LatencyMs:  stage.LatencyMs,
				ErrorRate:  stage.ErrorRate,
				BurstPause: stage.BurstPause,
			}
		}
	}
	return active
}
