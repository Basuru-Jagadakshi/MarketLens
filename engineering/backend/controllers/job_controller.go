package controllers

import (
	"marketlens-go-backend/models"
	"marketlens-go-backend/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)



type JobController struct {
	repo *repositories.JobRepository
}

func NewJobController(repo *repositories.JobRepository) *JobController {
	return &JobController{repo: repo}
}

func (ctrl *JobController) CreateJobHandler(c *gin.Context) {
	var input models.JobPost

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	createdJob, err := ctrl.repo.CreateJob(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to process job publication data",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, createdJob)
}


func (ctrl *JobController) GetAllJobsHandler(c *gin.Context) {

	jobs, err := ctrl.repo.GetAllJobs()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve job listings",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, jobs)
}