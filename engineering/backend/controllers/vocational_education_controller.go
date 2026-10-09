package controllers

import (
	"errors"
	"fmt"
	"marketlens-go-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Vocational Education CRUD
func (ctrl *JobController) CreateVocationalEducationHandler(c *gin.Context) {
	var item models.VocationalEducation
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}
	if err := ctrl.repo.CreateVocationalEducation(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vocational education", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (ctrl *JobController) GetAllVocationalEducationsHandler(c *gin.Context) {
	items, err := ctrl.repo.GetAllVocationalEducations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vocational education levels", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": len(items), "vocational_educations": items})
}

func (ctrl *JobController) GetVocationalEducationByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	item, err := ctrl.repo.GetVocationalEducationByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vocational education not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vocational education", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) UpdateVocationalEducationHandler(c *gin.Context) {
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
	item, err := ctrl.repo.UpdateVocationalEducation(id, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vocational education", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) DeleteVocationalEducationHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	if err := ctrl.repo.DeleteVocationalEducation(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vocational education", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Vocational education deleted successfully"})
}
