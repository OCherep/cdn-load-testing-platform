package geo

import (
	"math/rand"
	"time"
)

type Region string

const (
	EU   Region = "EU"
	US   Region = "US"
	ASIA Region = "ASIA"
)

type Config struct {
	BaseLatencyMs int
	JitterMs      int
}

var regionProfiles = map[Region]Config{
	EU: {
		BaseLatencyMs: 30,
		JitterMs:      20,
	},
	US: {
		BaseLatencyMs: 120,
		JitterMs:      50,
	},
	ASIA: {
		BaseLatencyMs: 250,
		JitterMs:      80,
	},
}

func Apply(region Region) {
	cfg, ok := regionProfiles[region]
	if !ok {
		return
	}

	jitter := rand.Intn(cfg.JitterMs)
	delay := cfg.BaseLatencyMs + jitter

	time.Sleep(time.Duration(delay) * time.Millisecond)
}
