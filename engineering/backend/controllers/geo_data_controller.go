package controllers

import (
	"errors"
	"fmt"
	"marketlens-go-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Geo Data CRUD
func (ctrl *JobController) CreateGeoDataHandler(c *gin.Context) {
	var item models.GeoData
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error})
		return
	}
	if err := ctrl.repo.CreateGeoData(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create geo data", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (ctrl *JobController) GetAllGeoDataHandler(c *gin.Context) {
	items, err := ctrl.repo.GetAllGeoData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch geo data", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": len(items), "geo_data": items})
}

func (ctrl *JobController) GetGeoDataByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	item, err := ctrl.repo.GetGeoDataByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Geo data not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch geo data", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) UpdateGeoDataHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}
	item, err := ctrl.repo.UpdateGeoData(id, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update geo data", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) DeleteGeoDataHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	if err := ctrl.repo.DeleteGeoData(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete geo data", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Geo data deleted successfully"})
}
