package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type statusResponse struct {
	Status string `json:"status"`
}

// HealthzHandler confirms the process itself is up and responding to HTTP - used for Kubernetes liveness probes
func (h *Handler) HealthzHandler(c *gin.Context) {
	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}

// ReadyzHandler confirms the database connection is reachable - used for Kubernetes readiness probes
func (h *Handler) ReadyzHandler(c *gin.Context) {
	if err := h.repo.Ping(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, statusResponse{Status: "ready"})
}
