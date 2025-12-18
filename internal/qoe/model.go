package qoe

import "cdn-load-platform/internal/metrics"

type Score struct {
	Value float64
}

func Compute(s metrics.Snapshot) Score {
	score := 100.0

	if s.P95Latency > 200 {
		score -= (s.P95Latency - 200) * 0.1
	}
	if s.ErrorRate > 0.01 {
		score -= s.ErrorRate * 100
	}
	if s.StickinessRatio < 0.9 {
		score -= (0.9 - s.StickinessRatio) * 50
	}

	if score < 0 {
		score = 0
	}

	return Score{Value: score}
}
