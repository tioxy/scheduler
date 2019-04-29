package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	scheduler "github.com/tioxy/scheduler/pkg"
)

var err error

func main() {
	r := gin.Default()

	r.GET("/", root)
	r.GET("/healthz", healthCheck)
	r.GET("/metrics", exportMetrics)

	v1 := r.Group("/api/v1/jobs")
	{
		v1.GET("/", getSimpleJobs)
		v1.POST("/", createSimpleJob)
		v1.GET("/:name", fetchSimpleJob)
		v1.PUT("/:name", updateSimpleJob)
		v1.DELETE("/:name", deleteSimpleJob)
	}

	r.Run()
}

func root(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{"status": http.StatusOK, "message": "Scheduler Running :D"},
	)
}

func getSimpleJobs(c *gin.Context) {

}

func createSimpleJob(c *gin.Context) {
	sj := scheduler.SimpleJob{}

	if err := c.BindJSON(sj); err != nil {
		fmt.Println(err)
		c.AbortWithStatus(400)
		return
	}

	k8s := scheduler.CreateKubernetesAPI()

	if sj.IsScheduled() {
		k8s.CreateCronJob(sj)
	} else {
		k8s.CreateJob(sj)
	}
}

func fetchSimpleJob(c *gin.Context) {

}

func deleteSimpleJob(c *gin.Context) {

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
