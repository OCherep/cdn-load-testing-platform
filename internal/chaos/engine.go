package chaos

import (
	"errors"
	"math/rand"
	"time"
)

type Config struct {
	Enabled    bool `json:"enabled"`
	LatencyMs  int  `json:"latency_ms"`
	ErrorRate  int  `json:"error_rate"` // percent (0-100)
	BurstPause bool `json:"burst_pause"`
}

func Apply(cfg Config) error {
	if !cfg.Enabled {
		return nil
	}

	if cfg.LatencyMs > 0 {
		time.Sleep(time.Duration(cfg.LatencyMs) * time.Millisecond)
	}

	if cfg.BurstPause {
		time.Sleep(time.Duration(rand.Intn(300)) * time.Millisecond)
	}

	if cfg.ErrorRate > 0 && rand.Intn(100) < cfg.ErrorRate {
		return errors.New("chaos: injected error")
	}

	return nil
}
