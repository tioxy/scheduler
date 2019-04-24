package main

import (
	"github.com/gin-gonic/gin"
	//scheduler "github.com/tioxy/scheduler/pkg"
)

var err error

func main() {
	r := gin.Default()
	r.GET("/", Root)

	r.GET("/jobs/", GetSimpleJobs)
	r.POST("/jobs/:name", CreateSimpleJob)
	r.GET("/jobs/:name", ReadSimpleJob)
	r.PUT("/jobs/:name", UpdateSimpleJob)
	r.DELETE("/jobs/:name", DeleteSimpleJob)

	//r.GET("/healthz", ListSimpleJobs)
	//r.GET("/metrics", ListSimpleJobs)

	r.Run(":8080")
}

func Root(c *gin.Context) {

}

func GetSimpleJobs(c *gin.Context) {

}

func CreateSimpleJob(c *gin.Context) {

}

func ReadSimpleJob(c *gin.Context) {

}

func DeleteSimpleJob(c *gin.Context) {

}

func UpdateSimpleJob(c *gin.Context) {

}
