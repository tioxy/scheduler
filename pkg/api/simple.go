package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	scheduler "github.com/tioxy/scheduler/pkg"
	"github.com/tioxy/scheduler/pkg/k8s"
)

func CreateSimpleJob(c *gin.Context) {
	sj, err := generateSimpleJobFromJSON(c)

	if err != nil {
		log.Warn().Msg("Could not parse JSON for SimpleJob | " + err.Error())
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

	log.Info().Msg("Creating Job from SimpleJob")
	err = api.CreateJob(sj)

	if err != nil {
		log.Error().Msg("Failed creating SimpleJob | " + err.Error())
		c.JSON(
			http.StatusUnprocessableEntity,
			gin.H{
				"status":  http.StatusUnprocessableEntity,
				"message": "could not create simplejob",
				"error":   err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"status":  http.StatusCreated,
			"message": "created simplejob",
		},
	)
}

func FetchSimpleJob(c *gin.Context) {
	sj := &scheduler.SimpleJob{
		Name:      c.Params.ByName("name"),
		Namespace: c.Params.ByName("namespace"),
	}

	api := k8s.CreateKubernetesAPI()

	log.Info().Msg("Fetching SimpleJob")
	job, err := api.FetchJob(sj.Name, sj.Namespace)

	if err != nil {
		log.Error().Msg("Failed fetching SimpleJob | " + err.Error())
		c.JSON(
			http.StatusUnprocessableEntity,
			gin.H{
				"status":  http.StatusUnprocessableEntity,
				"message": "could not fetch simplejob",
				"error":   err.Error(),
			},
		)
		return
	}

	simpleJob := scheduler.ConvertJobToSimpleJob(*job)

	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  http.StatusOK,
			"message": "fetched simplejob",
			"data":    simpleJob,
		},
	)
}

func DeleteSimpleJob(c *gin.Context) {
	sj := &scheduler.SimpleJob{
		Name:      c.Params.ByName("name"),
		Namespace: c.Params.ByName("namespace"),
	}

	api := k8s.CreateKubernetesAPI()

	log.Info().Msg("Deleting SimpleJob")
	err := api.DeleteJob(sj.Name, sj.Namespace)

	if err != nil {
		log.Error().Msg("Failed deleting SimpleJob | " + err.Error())
		c.JSON(
			http.StatusUnprocessableEntity,
			gin.H{
				"status":  http.StatusUnprocessableEntity,
				"message": "could not delete simplejob",
				"error":   err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusAccepted,
		gin.H{
			"status":  http.StatusAccepted,
			"message": "simplejob marked for deletion",
		},
	)
}

func ListAllSimpleJobs(c *gin.Context) {
	api := k8s.CreateKubernetesAPI()

	log.Info().Msg("Listing SimpleJobs from all namespaces")
	allJobs, err := api.ListJobs(k8s.AllNamespaces)

	if err != nil {
		log.Error().Msg("Failed listing SimpleJobs from all namespaces | " + err.Error())
		c.JSON(
			http.StatusUnprocessableEntity,
			gin.H{
				"status":  http.StatusUnprocessableEntity,
				"message": "could not list simplejobs",
				"error":   err.Error(),
			},
		)
		return
	}

	simpleJobs := []scheduler.SimpleJob{}

	for _, job := range allJobs {
		simpleJobs = append(simpleJobs, scheduler.ConvertJobToSimpleJob(job))
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  http.StatusOK,
			"message": "listed simplejobs",
			"data":    simpleJobs,
		},
	)
}

func ListSimpleJobsFromNamespace(c *gin.Context) {
	namespace := c.Params.ByName("namespace")

	api := k8s.CreateKubernetesAPI()

	log.Info().Msg("Listing SimpleJobs from namespace " + namespace)
	namespacedJobs, err := api.ListJobs(namespace)

	if err != nil {
		log.Error().Msg("Failed listing SimpleJobs from single namespace | " + err.Error())
		c.JSON(
			http.StatusUnprocessableEntity,
			gin.H{
				"status":  http.StatusUnprocessableEntity,
				"message": "could not list simplejobs from namespace",
				"error":   err.Error(),
			},
		)
		return
	}

	simpleJobs := []scheduler.SimpleJob{}

	for _, job := range namespacedJobs {
		simpleJobs = append(simpleJobs, scheduler.ConvertJobToSimpleJob(job))
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  http.StatusOK,
			"message": "listed simplejobs from namespace",
			"data":    simpleJobs,
		},
	)
}
