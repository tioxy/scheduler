package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	scheduler "github.com/tioxy/scheduler/pkg"
	"github.com/tioxy/scheduler/pkg/k8s"
)

func CreateScheduledSimpleJob(c *gin.Context) {
	sj, err := generateSimpleJobFromJSON(c)

	if err != nil {
		log.Warn().Msg("Could not parse JSON for scheduled SimpleJob | " + err.Error())
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": "could not parse json",
				"error":   err.Error(),
			},
		)
		return
	}

	api := k8s.CreateKubernetesAPI()

	log.Info().Msg("Creating CronJob from scheduled SimpleJob")
	err = api.CreateCronJob(sj)

	if err != nil {
		log.Error().Msg("Failed creating scheduled SimpleJob | " + err.Error())
		c.JSON(
			http.StatusUnprocessableEntity,
			gin.H{
				"status":  http.StatusUnprocessableEntity,
				"message": "could not create scheduled simplejob",
				"error":   err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"status":  http.StatusCreated,
			"message": "created scheduled simplejob",
		},
	)
}

func FetchScheduledSimpleJob(c *gin.Context) {
	sj := &scheduler.SimpleJob{
		Name:      c.Params.ByName("name"),
		Namespace: c.Params.ByName("namespace"),
	}

	api := k8s.CreateKubernetesAPI()

	log.Info().Msg("Fetching scheduled SimpleJob")
	cronJob, err := api.FetchCronJob(sj.Name, sj.Namespace)

	if err != nil {
		log.Error().Msg("Failed fetching scheduled SimpleJob | " + err.Error())
		c.JSON(
			http.StatusUnprocessableEntity,
			gin.H{
				"status":  http.StatusUnprocessableEntity,
				"message": "could not fetch scheduled simplejob",
				"error":   err.Error(),
			},
		)
		return
	}

	simpleJob := scheduler.ConvertCronJobToSimpleJob(*cronJob)

	c.JSON(
		http.StatusAccepted,
		gin.H{
			"status":  http.StatusAccepted,
			"message": "fetched scheduled simplejob",
			"data":    simpleJob,
		},
	)
}

func DeleteScheduledSimpleJob(c *gin.Context) {
	sj := &scheduler.SimpleJob{
		Name:      c.Params.ByName("name"),
		Namespace: c.Params.ByName("namespace"),
	}

	api := k8s.CreateKubernetesAPI()

	log.Info().Msg("Deleting scheduled SimpleJob")
	err := api.DeleteCronJob(sj.Name, sj.Namespace)

	if err != nil {
		log.Error().Msg("Failed deleting scheduled SimpleJob | " + err.Error())
		c.JSON(
			http.StatusUnprocessableEntity,
			gin.H{
				"status":  http.StatusUnprocessableEntity,
				"message": "could not delete scheduled simplejob",
				"error":   err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusAccepted,
		gin.H{
			"status":  http.StatusAccepted,
			"message": "scheduled simplejob marked for deletion",
		},
	)
}

func UpdateScheduledSimpleJob(c *gin.Context) {
	sj, err := generateSimpleJobFromJSON(c)

	if err != nil {
		log.Warn().Msg("Could not parse JSON for scheduled SimpleJob | " + err.Error())
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": "could not parse json",
				"error":   err.Error(),
			},
		)
		return
	}

	if !sj.IsScheduled() {
		log.Warn().Msg("Could not update scheduled SimpleJob because it is missing 'cron' key.")
		c.JSON(
			http.StatusAccepted,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": "could not update scheduled simplejob",
				"error":   "missing 'cron' key",
			},
		)
		return
	}

	api := k8s.CreateKubernetesAPI()

	log.Info().Msg("Updating CronJob from scheduled SimpleJob")
	err = api.UpdateCronJob(sj)

	if err != nil {
		log.Error().Msg("Failed updating scheduled SimpleJob | " + err.Error())
		c.JSON(
			http.StatusUnprocessableEntity,
			gin.H{
				"status":  http.StatusUnprocessableEntity,
				"message": "could not update scheduled simplejob",
				"error":   err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusAccepted,
		gin.H{
			"status":  http.StatusAccepted,
			"message": "scheduled simplejob marked for update",
		},
	)
}
