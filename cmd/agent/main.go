package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"cdn-load-platform/internal/chaos"
	"cdn-load-platform/internal/geo"
	"cdn-load-platform/internal/load"
	"cdn-load-platform/internal/metrics"
	"cdn-load-platform/internal/state"
)

var (
	currentChaos chaos.Config
	chaosMu      sync.RWMutex
)

func main() {
	log.Println("[agent] starting")

	// --------------------
	// ENV
	// --------------------
	profileBucket := os.Getenv("PROFILE_BUCKET")
	profileKey := os.Getenv("PROFILE_KEY")
	testID := os.Getenv("TEST_ID")
	awsRegion := os.Getenv("AWS_REGION")

	if profileBucket == "" || profileKey == "" || testID == "" {
		log.Fatal("PROFILE_BUCKET, PROFILE_KEY or TEST_ID not set")
	}

	// --------------------
	// Metrics
	// --------------------
	metrics.Start()

	// --------------------
	// State store (DynamoDB)
	// --------------------
	stateStore, err := state.NewDynamoStore(awsRegion)
	if err != nil {
		log.Fatalf("state store init failed: %v", err)
	}

	// --------------------
	// Load profile
	// --------------------
	profile, err := load.LoadProfileFromS3(profileBucket, profileKey)
	if err != nil {
		log.Fatalf("profile load failed: %v", err)
	}

	// --------------------
	// Adaptive load engine
	// --------------------
	engine := load.NewAdaptiveEngine(
		profile.MinRPS,
		profile.MaxRPS,
		profile.Step,
	)

	// --------------------
	// Chaos config refresher
	// --------------------
	go refreshTestState(stateStore, testID)

	// --------------------
	// HTTP client
	// --------------------
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// --------------------
	// Rate limiter ticker
	// --------------------
	rateTicker := time.NewTicker(time.Second)
	defer rateTicker.Stop()

	// --------------------
	// geo region
	// --------------------
	region := geo.DetectRegion()
	log.Printf("[agent] region=%s", region)

	var testStateMu sync.RWMutex
	var currentTestState *state.TestState

	// --------------------
	// Main loop
	// --------------------
	for range rateTicker.C {

		rps := engine.CurrentRPS()
		interval := time.Second / time.Duration(rps)

		reqTicker := time.NewTicker(interval)

		go func() {
			for range reqTicker.C {
				go executeRequest(client, profile.TargetURL)
			}
		}()

		// Adaptive adjustment once per second
		m := metrics.Snapshot()
		newRPS := engine.Adjust(m)
		log.Printf("[agent] RPS adjusted to %d", newRPS)
	}
}

func executeRequest(client *http.Client, url string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		metrics.RecordError()
		return
	}

	// GEO simulation
	geo.Apply(region)

	// --------------------
	// CHAOS INJECTION
	// --------------------
	chaosMu.RLock()
	err = chaos.Apply(currentChaos)
	chaosMu.RUnlock()

	if err != nil {
		metrics.RecordError()
		return
	}

	// --------------------
	// HTTP request
	// --------------------
	start := time.Now()
	resp, err := client.Do(req)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		metrics.RecordError()
		return
	}
	defer resp.Body.Close()

	// --------------------
	// Edge metrics
	// --------------------
	edge := resp.Header.Get("X-Cache-Node")
	ip := resp.Request.URL.Hostname()

	//metrics.RecordEdge(edge, ip, latency)
	metrics.RecordLatency(edge, string(region), latency)
	//metrics.RecordLatency(latency)
	metrics.RecordError(string(region))
}

func refreshTestState(store *state.DynamoStore, testID string) {
	for {
		test, err := store.GetTest(context.Background(), testID)
		if err == nil {
			testStateMu.Lock()
			currentTestState = &test
			testStateMu.Unlock()
		}
		testStateMu.RLock()
		ts := currentTestState
		testStateMu.RUnlock()

		if ts == nil || ts.Status != "running" {
			time.Sleep(200 * time.Millisecond)
			return
		}
		time.Sleep(5 * time.Second)
	}
}
