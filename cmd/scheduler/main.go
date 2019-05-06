package main

import (
	"net/http"
	"os"

	ginLogger "github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/tioxy/scheduler/pkg/api"
	"github.com/zsais/go-gin-prometheus"
)

var err error

func main() {
	r := gin.New()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Custom logger
	subLog := zerolog.New(os.Stdout).With().Str("app", "scheduler").Logger()
	r.Use(
		ginLogger.SetLogger(
			ginLogger.Config{
				Logger:   &subLog,
				UTC:      true,
				SkipPath: []string{"/healthz"},
			},
		),
	)

	// Custom metrics
	prom := ginprometheus.NewPrometheus("gin")
	prom.Use(r)

	// Routes
	r.GET("/", root)
	r.GET("/healthz", healthCheck)

	simpleV1 := r.Group("/api/v1/jobs/simple")
	{
		simpleV1.GET("/", api.ListAllSimpleJobs)
		simpleV1.POST("/", api.CreateSimpleJob)
		simpleV1.GET("/:namespace", api.ListSimpleJobsFromNamespace)
		simpleV1.GET("/:namespace/:name", api.FetchSimpleJob)
		simpleV1.DELETE("/:namespace/:name", api.DeleteSimpleJob)
	}

	scheduledV1 := r.Group("/api/v1/jobs/scheduled")
	{
		scheduledV1.GET("/", api.ListAllScheduledSimpleJobs)
		scheduledV1.POST("/", api.CreateScheduledSimpleJob)
		scheduledV1.GET("/:namespace", api.ListScheduledSimpleJobsFromNamespace)
		scheduledV1.GET("/:namespace/:name", api.FetchScheduledSimpleJob)
		scheduledV1.DELETE("/:namespace/:name", api.DeleteScheduledSimpleJob)
		scheduledV1.PUT("/:namespace/:name", api.UpdateScheduledSimpleJob)
	}

	r.Run()
}

func root(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{"status": http.StatusOK, "message": "Scheduler Running :D"},
	)
}

func healthCheck(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  http.StatusOK,
			"message": "ok",
		},
	)
}
