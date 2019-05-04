package api

import (
	"github.com/gin-gonic/gin"
	scheduler "github.com/tioxy/scheduler/pkg"
)

func generateSimpleJobFromJSON(c *gin.Context) (scheduler.SimpleJob, error) {
	sj := &scheduler.SimpleJob{}
	err := c.BindJSON(sj)

	if err != nil {
		return *sj, err
	}

	return *sj, nil
}

func convertSimpleJobToJSON() {

}
