package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	//scheduler "github.com/tioxy/scheduler/pkg"
)

var err error

func main() {
	r := gin.Default()

	r.GET("/", Root)
	r.GET("/healthz", healthCheck)
	r.GET("/metrics", exportMetrics)

	v1 := r.Group("/api/v1/jobs")
	{
		v1.GET("/", getSimpleJobs)
		v1.POST("/:name", createSimpleJob)
		v1.GET("/:name", fetchSimpleJob)
		v1.PUT("/:name", updateSimpleJob)
		v1.DELETE("/:name", deleteSimpleJob)
	}

	r.Run(":8080")
}

func Root(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{"status": http.StatusOK, "message": "Scheduler Running :D"},
	)
}

func getSimpleJobs(c *gin.Context) {

}

func createSimpleJob(c *gin.Context) {

}

func fetchSimpleJob(c *gin.Context) {

}

func deleteSimpleJob(c *gin.Context) {

}

func updateSimpleJob(c *gin.Context) {

}

func healthCheck(c *gin.Context) {

}

func exportMetrics(c *gin.Context) {

}
