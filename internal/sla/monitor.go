package sla

import (
	"log"
	"time"

	"cdn-load-platform/internal/report"
	"cdn-load-platform/internal/state"
)

func Monitor(
	store *state.Store,
	onBreach func(evidence report.SLAEvidence),
) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		tests, _ := store.List()

		for _, test := range tests {
			if test.Status != "running" {
				continue
			}

			m := store.GetMetricsSnapshot(test.TestID)

			evidence := report.SLAEvidence{
				TestID:          test.TestID,
				TargetURL:       test.TargetURL,
				StartTime:       time.Unix(test.StartedAt, 0),
				EndTime:         time.Now(),
				AvgLatencyMs:    m.AvgLatency,
				P95LatencyMs:    m.P95Latency,
				ErrorRate:       m.ErrorRate,
				StickinessRatio: m.StickinessRatio,

				LatencySLAms:  test.SLA.LatencyMs,
				ErrorRateSLA:  test.SLA.ErrorRate,
				StickinessSLA: test.SLA.Stickiness,
			}

			Evaluate(&evidence)

			if IsBreached(evidence) {
				log.Printf("[SLA] BREACH test=%s", test.TestID)
				onBreach(evidence)
			}
		}
	}
}
