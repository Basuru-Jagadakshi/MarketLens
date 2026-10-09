package controllers

import (
	"marketlens-go-backend/crawler"
	"marketlens-go-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	var rawJobs []crawler.RawJobInput
	if err := c.ShouldBindJSON(&rawJobs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body structural mapping",
			"details": err.Error(),
		})
		return
	}

	if len(rawJobs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The job batch cannot be empty"})
		return
	}

	failedJobs := ctrl.ingestion.ProcessBatch(c.Request.Context(), rawJobs)

	c.JSON(http.StatusOK, gin.H{
		"message":      "Batch processed",
		"submitted":    len(rawJobs),
		"failed_count": len(failedJobs),
		"failed_jobs":  failedJobs,
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
