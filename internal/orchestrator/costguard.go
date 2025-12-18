package orchestrator

import "time"

func MonitorCost(max float64, burnRateFn func() float64, stopFn func()) {
	ticker := time.NewTicker(time.Minute)
	var spent float64

	for range ticker.C {
		spent += burnRateFn()
		if spent >= max {
			stopFn()
			return
		}
	}
}
