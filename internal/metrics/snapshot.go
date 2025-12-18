package metrics

type Snapshot struct {
	AvgLatency      float64
	P95Latency      float64
	ErrorRate       float64
	StickinessRatio float64
	RPS             int
}

func SnapshotNow() Snapshot {
	return Snapshot{
		AvgLatency:      AvgLatency(),
		P95Latency:      P95Latency(),
		ErrorRate:       ErrorRate(),
		StickinessRatio: StickinessRatio(),
		RPS:             CurrentRPS(),
	}
}
