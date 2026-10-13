package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) StartCrawlerRunHandler(c *gin.Context) {
	run, err := h.crawlerRuns.CreateCrawlerRun(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "Failed to initialize new tracker session initialization block", err)
		return
	}
	c.JSON(http.StatusCreated, run)
}

func (h *Handler) CompleteCrawlerRunHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
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

	if err := h.crawlerRuns.CompleteCrawlerRun(c.Request.Context(), uint(id), req.Status); err != nil {
		respondRepoError(c, "Failed to execute closure metrics logic sequence", err)
		return
	}

	c.JSON(http.StatusOK, messageResponse{Message: fmt.Sprintf("Crawler run %d finalized with state: %s", id, req.Status)})
}
