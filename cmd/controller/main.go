package main

import (
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"cdn-load-platform/internal/auth"
	"cdn-load-platform/internal/chaos"
	"cdn-load-platform/internal/cost"
	"cdn-load-platform/internal/report"
	"cdn-load-platform/internal/state"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	r := gin.Default()

	/*
		STATE STORE
	*/
	store := state.NewStore("cdn-load-tests")

	/*
		AUTH
	*/
	r.POST("/auth/login", func(c *gin.Context) {
		token, err := auth.Generate("admin")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"token": token})
	})

	api := r.Group("/tests")
	api.Use(auth.JWTMiddleware())

	/*
		CREATE TEST
	*/
	api.POST("", func(c *gin.Context) {
		var req struct {
			ProfileKey   string  `json:"profile_key"`
			Nodes        int     `json:"nodes"`
			Sessions     int     `json:"sessions"`
			DurationSec  int64   `json:"duration_sec"`
			MaxBudgetUSD float64 `json:"max_budget_usd"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		id := uuid.New().String()

		/*
			COST GUARD
		*/
		costUSD := cost.Estimate(
			req.Nodes,
			"c6i.large",
			float64(req.DurationSec)/3600,
		)

		if req.MaxBudgetUSD > 0 && costUSD > req.MaxBudgetUSD {
			c.JSON(400, gin.H{"error": "budget exceeded", "estimate": costUSD})
			return
		}

		store.PutTest(state.TestState{
			TestID:          id,
			Status:          "created",
			ProfileKey:      req.ProfileKey,
			Nodes:           req.Nodes,
			Sessions:        req.Sessions,
			StartedAt:       time.Now().Unix(),
			TTL:             req.DurationSec,
			CostEstimateUSD: costUSD,
		})

		c.JSON(201, gin.H{
			"test_id":       id,
			"cost_estimate": costUSD,
		})
	})

	/*
		LIST TESTS
	*/
	api.GET("", func(c *gin.Context) {
		tests, err := store.List()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, tests)
	})

	/*
		START TEST
	*/
	api.POST("/:id/start", func(c *gin.Context) {
		id := c.Param("id")

		err := exec.Command(
			"terraform",
			"-chdir=terraform/load-nodes",
			"apply",
			"-auto-approve",
			"-var", "test_id="+id,
		).Run()

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		store.UpdateStatus(id, "running")
		c.Status(204)
	})

	/*
		STOP TEST
	*/
	api.POST("/:id/stop", func(c *gin.Context) {
		id := c.Param("id")

		_ = exec.Command(
			"terraform",
			"-chdir=terraform/load-nodes",
			"destroy",
			"-auto-approve",
		).Run()

		store.UpdateStatus(id, "finished")
		c.Status(204)
	})

	/*
		PAUSE / RESUME
	*/
	api.POST("/:id/pause", func(c *gin.Context) {
		store.UpdateStatus(c.Param("id"), "paused")
		c.Status(204)
	})

	api.POST("/:id/resume", func(c *gin.Context) {
		store.UpdateStatus(c.Param("id"), "running")
		c.Status(204)
	})

	/*
		UPDATE RPS
	*/
	api.POST("/:id/rps", func(c *gin.Context) {
		var body struct {
			RPS int `json:"rps"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		store.UpdateDesiredRPS(c.Param("id"), body.RPS)
		c.Status(204)
	})

	/*
		CHAOS CONFIG
	*/
	api.POST("/:id/chaos", func(c *gin.Context) {
		var cfg state.ChaosConfig
		if err := c.BindJSON(&cfg); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		store.UpdateChaos(c.Param("id"), cfg)
		c.Status(204)
	})

	/*
		CHAOS SCHEDULE
	*/
	api.POST("/:id/chaos/schedule", func(c *gin.Context) {
		var schedule chaos.Schedule
		if err := c.BindJSON(&schedule); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		store.UpdateChaosSchedule(c.Param("id"), schedule)
		c.Status(204)
	})

	/*
		EXTEND TTL
	*/
	api.POST("/:id/extend", func(c *gin.Context) {
		var body struct {
			Seconds int64 `json:"seconds"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		store.ExtendTTL(c.Param("id"), body.Seconds)
		c.Status(204)
	})

	/*
		WEBSOCKET (LIVE UI)
	*/
	r.GET("/ws/tests/:id", func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer ws.Close()

		for {
			time.Sleep(1 * time.Second)
			ws.WriteJSON(gin.H{"status": "alive"})
		}
	})

	api.POST("/:id/report", func(c *gin.Context) {
		testID := c.Param("id")

		test, err := store.Get(testID)
		if err != nil {
			c.JSON(404, gin.H{"error": "test not found"})
			return
		}

		metrics := store.GetMetricsSnapshot(testID)

		evidence := report.SLAEvidence{
			TestID:          testID,
			TargetURL:       test.TargetURL,
			StartTime:       time.Unix(test.StartedAt, 0),
			EndTime:         time.Now(),
			AvgLatencyMs:    metrics.AvgLatency,
			P95LatencyMs:    metrics.P95Latency,
			ErrorRate:       metrics.ErrorRate,
			StickinessRatio: metrics.StickinessRatio,

			LatencySLAms:  200,
			ErrorRateSLA:  0.01,
			StickinessSLA: 0.90,
		}

		file, err := report.ExportSLAEvidencePDF(evidence)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"file": file,
		})
	})

	log.Println("[controller] listening on :8080")
	r.Run(":8080")
}
