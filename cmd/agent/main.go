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

/*
GLOBAL STATE
*/
var (
	currentChaos chaos.Config
	chaosMu      sync.RWMutex

	currentTestState *state.TestState
	testStateMu      sync.RWMutex
)

func main() {
	log.Println("[agent] starting")

	/*
		ENV
	*/
	profileBucket := os.Getenv("PROFILE_BUCKET")
	profileKey := os.Getenv("PROFILE_KEY")
	testID := os.Getenv("TEST_ID")
	awsRegion := os.Getenv("AWS_REGION")

	if profileBucket == "" || profileKey == "" || testID == "" {
		log.Fatal("PROFILE_BUCKET, PROFILE_KEY or TEST_ID not set")
	}

	/*
		METRICS
	*/
	metrics.Start()

	/*
		STATE STORE
	*/
	stateStore, err := state.NewDynamoStore(awsRegion)
	if err != nil {
		log.Fatalf("state store init failed: %v", err)
	}

	/*
		LOAD PROFILE
	*/
	profile, err := load.LoadProfileFromS3(profileBucket, profileKey)
	if err != nil {
		log.Fatalf("profile load failed: %v", err)
	}

	/*
		ADAPTIVE LOAD ENGINE
	*/
	engine := load.NewAdaptiveEngine(
		profile.MinRPS,
		profile.MaxRPS,
		profile.Step,
	)

	/*
		GEO
	*/
	region := geo.DetectRegion()
	log.Printf("[agent] geo region = %s", region)

	/*
		STICKINESS
	*/
	stickiness := load.NewStickinessTracker()

	/*
		HTTP CLIENT
	*/
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	/*
		REFRESH TEST STATE (status, chaos, rps)
	*/
	go refreshTestState(stateStore, testID)

	/*
		MAIN LOAD LOOP
	*/
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {

		testStateMu.RLock()
		ts := currentTestState
		testStateMu.RUnlock()

		if ts == nil || ts.Status != "running" {
			continue
		}

		rps := engine.CurrentRPS()
		if rps <= 0 {
			continue
		}

		violated, reason := load.CheckSLA(load.SLA{
			LatencyMs:  currentTestState.SLA.LatencyMs,
			ErrorRate:  currentTestState.SLA.ErrorRate,
			Stickiness: currentTestState.SLA.Stickiness,
		})

		if violated {
			stateStore.MarkSLAViolation(
				context.Background(),
				testID,
				reason,
			)
			log.Println("[SLA] VIOLATED:", reason)
		}

		interval := time.Second / time.Duration(rps)

		for i := 0; i < rps; i++ {
			go executeRequest(
				client,
				profile.TargetURL,
				region,
				stickiness,
			)
			time.Sleep(interval)
		}

		/*
			ADAPTIVE ADJUSTMENT
		*/
		snapshot := metrics.Snapshot()
		newRPS := engine.Adjust(snapshot)
		log.Printf("[agent] RPS adjusted -> %d", newRPS)
	}
}

/*
EXECUTE SINGLE REQUEST
*/
func executeRequest(
	client *http.Client,
	url string,
	region string,
	stickiness *load.StickinessTracker,
) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		metrics.RecordError(region)
		return
	}

	/*
		GEO SIMULATION
	*/
	geo.Apply(region)

	/*
		CHAOS
	*/
	testStateMu.RLock()
	cfg := resolveChaos(currentTestState)
	testStateMu.RUnlock()

	if err := chaos.Apply(cfg); err != nil {
		metrics.RecordError(region)
		return
	}

	if err != nil {
		metrics.RecordError(region)
		return
	}

	/*
		HTTP REQUEST
	*/
	start := time.Now()
	resp, err := client.Do(req)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		metrics.RecordError(region)
		return
	}
	defer resp.Body.Close()

	/*
		EDGE METRICS
	*/
	edge := resp.Header.Get("X-Cache-Node")
	if edge == "" {
		edge = "unknown"
	}

	clientID := req.RemoteAddr
	stickiness.Record(clientID, edge)

	metrics.RecordLatency(edge, region, latency)
}

/*
REFRESH TEST STATE FROM CONTROLLER (DYNAMO)
*/
func refreshTestState(store *state.DynamoStore, testID string) {
	for {
		test, err := store.GetTest(context.Background(), testID)
		if err == nil {

			testStateMu.Lock()
			currentTestState = &test
			testStateMu.Unlock()

			chaosMu.Lock()
			currentChaos = chaos.Config{
				Enabled:    test.ChaosConfig.Enabled,
				LatencyMs:  test.ChaosConfig.LatencyMs,
				ErrorRate:  test.ChaosConfig.ErrorRate,
				BurstPause: test.ChaosConfig.BurstPause,
			}
			chaosMu.Unlock()
		}

		time.Sleep(5 * time.Second)
	}
}

func resolveChaos(test *state.TestState) chaos.Config {
	if test == nil {
		return chaos.Config{}
	}

	cfg := test.ChaosSchedule.Active(
		time.Now(),
		time.Unix(test.StartedAt, 0),
	)

	if cfg == nil {
		return chaos.Config{}
	}

	return *cfg
}
