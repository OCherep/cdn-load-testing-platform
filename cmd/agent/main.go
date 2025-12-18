package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"cdn-load-platform/internal/chaos"
	"cdn-load-platform/internal/load"
	"cdn-load-platform/internal/metrics"
	"cdn-load-platform/internal/state"
)

var (
	stateMu      sync.RWMutex
	currentState *state.TestState
)

func main() {
	log.Println("[agent] starting")

	// =====================
	// ENVIRONMENT
	// =====================
	testID := os.Getenv("TEST_ID")
	if testID == "" {
		log.Fatal("TEST_ID env not set")
	}

	profileBucket := os.Getenv("PROFILE_BUCKET")
	profileKey := os.Getenv("PROFILE_KEY")
	if profileBucket == "" || profileKey == "" {
		log.Fatal("PROFILE_BUCKET or PROFILE_KEY env not set")
	}

	// =====================
	// METRICS
	// =====================
	metrics.Start()

	// =====================
	// STATE STORE
	// =====================
	store := state.NewDynamoStoreFromEnv()
	go refreshTestState(store, testID)

	// =====================
	// LOAD PROFILE
	// =====================
	profile, err := load.LoadProfileFromS3(profileBucket, profileKey)
	if err != nil {
		log.Fatalf("failed to load profile: %v", err)
	}

	// =====================
	// LOAD ENGINE
	// =====================
	engine := load.NewAdaptiveEngine(
		profile.MinRPS,
		profile.MaxRPS,
	)

	// =====================
	// EDGE STICKINESS
	// =====================
	tracker := load.NewStickinessTracker()

	// =====================
	// HTTP CLIENT
	// =====================
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// =====================
	// MAIN LOOP
	// =====================
	for {
		stateMu.RLock()
		ts := currentState
		stateMu.RUnlock()

		// wait until state appears
		if ts == nil {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		// pause / stop handling
		if ts.Status != "running" {
			if ts.Status == "stopped" {
				log.Println("[agent] test stopped")
				return
			}
			time.Sleep(300 * time.Millisecond)
			continue
		}

		// TTL expiration
		if time.Now().Unix() > ts.ExpiresAt {
			log.Println("[agent] test expired")
			return
		}

		// live RPS override
		if ts.DesiredRPS != nil {
			engine.SetTarget(*ts.DesiredRPS)
		}

		// =====================
		// EXECUTE ONE STEP
		// =====================
		engine.Step(func() {
			// Apply chaos BEFORE request
			chaos.Apply(ts.ChaosConfig)

			req, err := http.NewRequest("GET", profile.URL, nil)
			if err != nil {
				return
			}

			start := time.Now()
			resp, err := client.Do(req)
			latency := time.Since(start).Milliseconds()

			if err != nil {
				metrics.RecordError()
				return
			}
			defer resp.Body.Close()

			// =====================
			// EDGE METRICS
			// =====================
			edge := resp.Header.Get("X-Cache-Node")
			if edge == "" {
				edge = "unknown"
			}

			host := resp.Request.URL.Hostname()
			clientID := host
			if !profile.StickyMode {
				clientID = strconv.FormatInt(time.Now().UnixNano(), 10)
			} // достатньо для stickiness

			tracker.Record(clientID, edge)
			ratio := tracker.Ratio(clientID)
			metrics.RecordStickiness(clientID, ratio)

			metrics.RecordEdge(edge, host, latency)
			metrics.RecordLatency(latency)
		})
	}
}

// =====================
// STATE REFRESH LOOP
// =====================
func refreshTestState(store *state.DynamoStore, testID string) {
	for {
		ts, err := store.GetTest(context.Background(), testID)
		if err == nil {
			stateMu.Lock()
			currentState = &ts
			stateMu.Unlock()
		}
		time.Sleep(5 * time.Second)
	}
}
