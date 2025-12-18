package state

func (s *Store) GetMetricsSnapshot(testID string) MetricsSnapshot {
	// ⚠️ MVP implementation
	// Реально тут буде Prometheus query

	return MetricsSnapshot{
		AvgLatency:      180,
		P95Latency:      240,
		ErrorRate:       0.012,
		StickinessRatio: 0.87,
	}
}
