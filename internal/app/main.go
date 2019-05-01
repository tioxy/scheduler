package main

import (
	"fmt"
	"net/http"
	"os"

	ginLogger "github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	scheduler "github.com/tioxy/scheduler/pkg"
)

var err error

func main() {
	r := gin.New()

	// Control zerolog INFO/DEBUG by Gin env
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	r.Use(ginLogger.SetLogger())

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

	r.GET("/", root)
	r.GET("/healthz", healthCheck)
	r.GET("/metrics", exportMetrics)

	v1 := r.Group("/api/v1/jobs")
	{
		v1.GET("/", getAllSimpleJobs)
		v1.GET("/:namespace", getSimpleJobsFromNamespace)
		v1.GET("/:namespace/:name", fetchSimpleJob)

		v1.POST("/", createSimpleJob)
		v1.PUT("/", updateSimpleJob)
		v1.DELETE("/", deleteSimpleJob)
	}

	r.Run()
}

func root(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{"status": http.StatusOK, "message": "Scheduler Running :D"},
	)
}

func getAllSimpleJobs(c *gin.Context) {
}

func createSimpleJob(c *gin.Context) {
	sj, err := generateSimpleJobFromJSON(c)

	if err != nil {
		log.Warn().Msg(fmt.Sprintf("Could not parse JSON for SimpleJob | %v", err))
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": "could not parse json",
			},
		)
		return
	}

	k8s := scheduler.CreateKubernetesAPI()

	if sj.IsScheduled() {
		log.Info().Msg("Creating CronJob from SimpleJob=" + sj.Name)
		err = k8s.CreateCronJob(sj)
	} else {
		log.Info().Msg("Creating Job from SimpleJob=" + sj.Name)
		err = k8s.CreateJob(sj)
	}

	if err != nil {
		log.Error().Msg(
			fmt.Sprintf("Failed creating SimpleJob=%s | %v", sj.Name, err),
		)
		c.JSON(
			http.StatusUnprocessableEntity,
			gin.H{
				"status":  http.StatusUnprocessableEntity,
				"message": "could not create SimpleJob=" + sj.Name,
			},
		)
		return
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"status":  http.StatusCreated,
			"message": "created SimpleJob=" + sj.Name,
		},
	)
}

func getSimpleJobsFromNamespace(c *gin.Context) {

}

func fetchSimpleJob(c *gin.Context) {

}

func deleteSimpleJob(c *gin.Context) {
	sj, err := generateSimpleJobFromJSON(c)

	if err != nil {
		log.Warn().Msg(fmt.Sprintf("Could not parse JSON for SimpleJob | %v", err))
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": "could not parse json",
			},
		)
		return
	}

	k8s := scheduler.CreateKubernetesAPI()

	if sj.IsScheduled() {
		log.Info().Msg("Deleting CronJob from SimpleJob=" + sj.Name)
		err = k8s.DeleteCronJob(sj)
	} else {
		log.Info().Msg("Deleting Job from SimpleJob=" + sj.Name)
		err = k8s.DeleteJob(sj)
	}

	if err != nil {
		log.Error().Msg(
			fmt.Sprintf("Failed deleting SimpleJob=%s | %v", sj.Name, err),
		)
		c.JSON(
			http.StatusUnprocessableEntity,
			gin.H{
				"status":  http.StatusUnprocessableEntity,
				"message": "could not delete SimpleJob=" + sj.Name,
			},
		)
		return
	}

	c.JSON(
		http.StatusAccepted,
		gin.H{
			"status":  http.StatusAccepted,
			"message": "Marked for deletion SimpleJob=" + sj.Name,
		},
	)

}

func updateSimpleJob(c *gin.Context) {

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

func exportMetrics(c *gin.Context) {

}

func generateSimpleJobFromJSON(c *gin.Context) (scheduler.SimpleJob, error) {
	sj := &scheduler.SimpleJob{}
	err := c.BindJSON(sj)

	if err != nil {
		return *sj, err
	}

	return *sj, nil
}
