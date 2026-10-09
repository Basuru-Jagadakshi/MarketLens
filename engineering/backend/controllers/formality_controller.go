package controllers

import (
	"errors"
	"fmt"
	"marketlens-go-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Formality CRUD
func (ctrl *JobController) CreateFormalityHandler(c *gin.Context) {
	var item models.Formality
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}
	if err := ctrl.repo.CreateFormality(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create formality", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (ctrl *JobController) GetAllFormalitiesHandler(c *gin.Context) {
	items, err := ctrl.repo.GetAllFormalities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch formalities", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": len(items), "formalities": items})
}

func (ctrl *JobController) GetFormalityByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	item, err := ctrl.repo.GetFormalityByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Formality not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch formality", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) UpdateFormalityHandler(c *gin.Context) {
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
	item, err := ctrl.repo.UpdateFormality(id, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update formality", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) DeleteFormalityHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	if err := ctrl.repo.DeleteFormality(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete formality", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Formality deleted successfully"})
}
