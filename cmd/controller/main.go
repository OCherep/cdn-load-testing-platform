package main

import (
	"cdn-load-platform/internal/chaos"
	"cdn-load-platform/internal/cost"
	"net/http"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"cdn-load-platform/internal/auth"
	"cdn-load-platform/internal/state"
)

var upgrader = websocket.Upgrader{}

func main() {
	r := gin.Default()
	store := state.NewStore("cdn-load-tests")

	r.POST("/auth/login", func(c *gin.Context) {
		token, _ := auth.Generate("admin")
		c.JSON(200, gin.H{"token": token})
	})

	api := r.Group("/tests")
	api.Use(JWTMiddleware())

	api.POST("", func(c *gin.Context) {
		var req struct {
			ProfileKey string `json:"profile_key"`
			Nodes      int    `json:"nodes"`
			Sessions   int    `json:"sessions"`
		}
		c.BindJSON(&req)

		id := uuid.New().String()

		store.PutTest(state.TestState{
			TestID:     id,
			Status:     "CREATED",
			ProfileKey: req.ProfileKey,
			Nodes:      req.Nodes,
			Sessions:   req.Sessions,
			StartedAt:  time.Now().Unix(),
			TTL:        3600,
		})

		c.JSON(201, gin.H{"test_id": id})
	})

	api.GET("", func(c *gin.Context) {
		tests, _ := store.List()
		c.JSON(200, tests)
	})

	api.POST("/:id/start", func(c *gin.Context) {
		id := c.Param("id")

		exec.Command("terraform",
			"-chdir=terraform/load-nodes",
			"apply", "-auto-approve",
			"-var", "test_id="+id,
		).Run()

		store.UpdateStatus(id, "RUNNING")
		c.Status(204)
	})

	api.POST("/:id/stop", func(c *gin.Context) {
		id := c.Param("id")

		exec.Command("terraform",
			"-chdir=terraform/load-nodes",
			"destroy", "-auto-approve").Run()

		store.UpdateStatus(id, "FINISHED")
		c.Status(204)
	})

	api.POST("/:id/chaos", func(c *gin.Context) {
		id := c.Param("id")

		var cfg state.ChaosConfig
		c.BindJSON(&cfg)

		store.UpdateChaos(id, cfg)
		c.Status(204)
	})

	api.POST("/:id/chaos/schedule", func(c *gin.Context) {
		id := c.Param("id")

		var schedule chaos.Schedule
		if err := c.BindJSON(&schedule); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		store.UpdateChaosSchedule(id, schedule)
		c.Status(204)
	})

	api.POST("/:id/rps", func(c *gin.Context) {
		id := c.Param("id")

		var body struct {
			RPS int `json:"rps"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		store.UpdateDesiredRPS(id, body.RPS)
		c.Status(204)
	})

	api.POST("/:id/pause", func(c *gin.Context) {
		store.UpdateStatus(c.Param("id"), "paused")
		c.Status(204)
	})

	api.POST("/:id/resume", func(c *gin.Context) {
		store.UpdateStatus(c.Param("id"), "running")
		c.Status(204)
	})

	api.POST("/:id/extend", func(c *gin.Context) {
		id := c.Param("id")

		var body struct {
			Seconds int64 `json:"seconds"`
		}
		c.BindJSON(&body)

		store.ExtendTTL(id, body.Seconds)
		c.Status(204)
	})

	r.GET("/ws/tests/:id", func(c *gin.Context) {
		ws, _ := upgrader.Upgrade(c.Writer, c.Request, nil)
		handleWS(ws)
	})

	costUSD := cost.Estimate(
		req.Nodes,
		"c6i.large",
		float64(req.DurationSec)/3600,
	)

	if costUSD > req.MaxBudgetUSD {
		c.JSON(400, gin.H{"error": "budget exceeded"})
		return
	}

	store.PutTest(state.TestState{
		TestID:          id,
		Status:          "CREATED",
		CostEstimateUSD: costUSD,
	})

	r.Run(":8080")
}
