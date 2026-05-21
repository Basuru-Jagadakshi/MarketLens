package controllers

import (
	"fmt"
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


func (ctrl *JobController) UpdateJobHandler(c *gin.Context) {
	
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID format parameter"})
		return
	}

	var input models.JobPost
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload configuration",
			"details": err.Error(),
		})
		return
	}

	updatedJob, err := ctrl.repo.UpdateJob(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to execute modifications on targeted job profile",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updatedJob)
}


func (ctrl *JobController) DeleteJobHandler(c *gin.Context) {
	
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID format parameter"})
		return
	}

	id, err := ctrl.repo.DeleteJob(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to execute deletion on targeted job profile",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Job post with ID %d has been successfully deleted", id),
	})
}