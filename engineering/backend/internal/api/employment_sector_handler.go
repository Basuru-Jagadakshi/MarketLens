package api

import (
	"marketlens-go-backend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createEmploymentSectorInput struct {
	Sector string `json:"sector" binding:"required,max=100"`
}

type updateEmploymentSectorInput struct {
	Sector string `json:"sector" binding:"required,max=100"`
}

type employmentSectorListResponse struct {
	Count             int                        `json:"count"`
	EmploymentSectors []models.EmploymentSector `json:"employment_sectors"`
}

// Employment Sector CRUD
func (h *Handler) CreateEmploymentSectorHandler(c *gin.Context) {
	var in createEmploymentSectorInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}
	item := models.EmploymentSector{Sector: in.Sector}
	if err := h.employmentSectors.Create(c.Request.Context(), &item); err != nil {
		respondRepoError(c, "Failed to create employment sector", err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *Handler) GetAllEmploymentSectorsHandler(c *gin.Context) {
	items, err := h.employmentSectors.GetAll(c.Request.Context(), 0, 0, nil, nil)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "Failed to fetch employment sectors", err)
		return
	}
	c.JSON(http.StatusOK, employmentSectorListResponse{Count: len(items), EmploymentSectors: items})
}

func (h *Handler) GetEmploymentSectorByIDHandler(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	item, err := h.employmentSectors.GetByID(c.Request.Context(), id)
	if err != nil {
		respondRepoError(c, "Failed to fetch employment sector", err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *Handler) UpdateEmploymentSectorHandler(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	var in updateEmploymentSectorInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}
	item, err := h.employmentSectors.Update(c.Request.Context(), id, map[string]interface{}{"sector": in.Sector})
	if err != nil {
		respondRepoError(c, "Failed to update employment sector", err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *Handler) DeleteEmploymentSectorHandler(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	if err := h.employmentSectors.Delete(c.Request.Context(), id); err != nil {
		respondRepoError(c, "Failed to delete employment sector", err)
		return
	}
	c.JSON(http.StatusOK, messageResponse{Message: "Employment sector deleted successfully"})
}
