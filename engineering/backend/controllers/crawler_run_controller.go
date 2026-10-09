package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
