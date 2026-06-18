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

func (ctrl *JobController) StartCrawlerRunHandler(c *gin.Context) {
	run, err := ctrl.repo.CreateCrawlerRun()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to initialize new tracker session initialization block",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, run)
}

func (ctrl *JobController) CompleteCrawlerRunHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid crawler run ID parameter"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"` // 'COMPLETED' or 'FAILED'
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing state status flag criteria"})
		return
	}

	if err := ctrl.repo.CompleteCrawlerRun(id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to execute closure metrics logic sequence",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Crawler run %d finalized with state: %s", id, req.Status)})
}

func (ctrl *JobController) GetJobsByBucketKeysHandler(c *gin.Context) {
	var req struct {
		BucketKeys []string `json:"bucket_keys" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request structure payload mapping",
			"details": err.Error(),
		})
		return
	}

	matchingJobs, err := ctrl.repo.GetJobsByBucketKeys(req.BucketKeys)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "LSH bucket radar execution sequence encountered a processing exception",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, matchingJobs)
}

func (ctrl *JobController) BatchSaveJobsHandler(c *gin.Context) {
	var payload models.BatchSavePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body structural mapping",
			"details": err.Error(),
		})
		return
	}

	if len(payload.NewJobs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The new_jobs collection buffer cannot be empty"})
		return
	}

	err := ctrl.repo.BatchSaveNewJobs(payload.NewJobs, payload.LshIndexes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Bulk insertion transaction routine failed execution",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "Successfully persisted unique job block chunk",
		"inserted_records": len(payload.NewJobs),
	})
}

func (ctrl *JobController) BatchUpdateDuplicatesHandler(c *gin.Context) {
	var payload models.BatchUpdatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload configuration mapping",
			"details": err.Error(),
		})
		return
	}

	if len(payload.Duplicates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The duplicates reference buffer cannot be empty"})
		return
	}

	err := ctrl.repo.BatchUpdateDuplicateJobs(payload.Duplicates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Bulk update checkpoint modifications failed execution",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "Successfully refreshed duplicate job keep-alive markers",
		"refreshed_records": len(payload.Duplicates),
	})
}

func (ctrl *JobController) ReconcileStaleVacanciesHandler(c *gin.Context) {
	var payload models.ReconciliationPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing active execution tracking sequence identifier criteria",
			"details": err.Error(),
		})
		return
	}

	closedCount, err := ctrl.repo.ReconcileStaleVacancies(payload.CrawlerRunID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Global snapshot reconciliation sweeping routine failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":               "System-wide stale vacancy cleanup sweep finalized",
		"reconciled_stale_jobs": closedCount,
	})
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