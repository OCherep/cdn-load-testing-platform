package autoscale

import (
	"log"
	"time"

	"cdn-load-platform/internal/state"
)

func Run(
	store *state.Store,
	metricsProvider func(testID string) Metrics,
	apply func(testID string, decision Decision),
) {
	predictor := NewPredictor()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		tests, err := store.List()
		if err != nil {
			continue
		}

		for _, test := range tests {
			if test.Status != "running" {
				continue
			}

			m := metricsProvider(test.TestID)
			decision := predictor.Decide(m, test.Nodes)

			if decision.ScaleNodes != test.Nodes {
				log.Printf(
					"[autoscale] test=%s nodes %d → %d (%s)",
					test.TestID,
					test.Nodes,
					decision.ScaleNodes,
					decision.Reason,
				)
				apply(test.TestID, decision)
			}
		}
	}
}
