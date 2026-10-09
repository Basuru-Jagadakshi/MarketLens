package controllers

import (
	"marketlens-go-backend/crawler"
	"marketlens-go-backend/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JobController struct {
	repo      *repositories.JobRepository
	ingestion *crawler.IngestionService
}

func NewJobController(repo *repositories.JobRepository, ingestion *crawler.IngestionService) *JobController {
	return &JobController{repo: repo, ingestion: ingestion}
}

// This function confirms the process itself is up and responding to HTTP - used for Kubernetes liveness probes
func (ctrl *JobController) HealthzHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// This function confirms the database connection is reachable - used for Kubernetes readiness probes
func (ctrl *JobController) ReadyzHandler(c *gin.Context) {
	if err := ctrl.repo.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}
