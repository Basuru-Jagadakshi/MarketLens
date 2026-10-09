package controllers

import (
	"fmt"
	"marketlens-go-backend/models"
	"net/http"
	"strconv"
	"time"

	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Methods need to multi level filtering in the Crawler
// Occupation levels by parent ids
func (ctrl *JobController) GetSubMajorGroupsByMajorGroupHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid major group id parameter"})
		return
	}

	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	var items []models.SubMajorGroup
	var err error

	if fromDateStr != "" && toDateStr != "" {
		fromDate, parseErr := time.Parse("2006-01-02", fromDateStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from-date format, expected YYYY-MM-DD"})
			return
		}
		toDate, parseErr := time.Parse("2006-01-02", toDateStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to-date format, expected YYYY-MM-DD"})
			return
		}
		if toDate.Before(fromDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "to-date must not be before from-date"})
			return
		}
		items, err = ctrl.repo.GetSubMajorGroupsByMajorGroupForDateRange(id, fromDate, toDate)
	} else {
		items, err = ctrl.repo.GetSubMajorGroupsByMajorGroup(id)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve sub major groups",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"major_group_id":   id,
		"count":            len(items),
		"sub_major_groups": items,
	})
}

func (ctrl *JobController) GetMinorGroupsBySubMajorGroupHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sub major group id parameter"})
		return
	}

	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	var items []models.MinorGroup
	var err error

	if fromDateStr != "" && toDateStr != "" {
		fromDate, parseErr := time.Parse("2006-01-02", fromDateStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from-date format, expected YYYY-MM-DD"})
			return
		}
		toDate, parseErr := time.Parse("2006-01-02", toDateStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to-date format, expected YYYY-MM-DD"})
			return
		}
		if toDate.Before(fromDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "to-date must not be before from-date"})
			return
		}
		items, err = ctrl.repo.GetMinorGroupsBySubMajorGroupForDateRange(id, fromDate, toDate)
	} else {
		items, err = ctrl.repo.GetMinorGroupsBySubMajorGroup(id)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve minor groups",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sub_major_group_id": id,
		"count":              len(items),
		"minor_groups":       items,
	})
}

func (ctrl *JobController) GetUnitGroupsByMinorGroupHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid minor group id parameter"})
		return
	}

	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	var items []models.UnitGroup
	var err error

	if fromDateStr != "" && toDateStr != "" {
		fromDate, parseErr := time.Parse("2006-01-02", fromDateStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from-date format, expected YYYY-MM-DD"})
			return
		}
		toDate, parseErr := time.Parse("2006-01-02", toDateStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to-date format, expected YYYY-MM-DD"})
			return
		}
		if toDate.Before(fromDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "to-date must not be before from-date"})
			return
		}
		items, err = ctrl.repo.GetUnitGroupsByMinorGroupForDateRange(id, fromDate, toDate)
	} else {
		items, err = ctrl.repo.GetUnitGroupsByMinorGroup(id)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve unit groups",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"minor_group_id": id,
		"count":          len(items),
		"unit_groups":    items,
	})
}

func (ctrl *JobController) GetOccupationGroupsByUnitGroupHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid unit group id parameter"})
		return
	}

	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	var items []models.OccupationGroup
	var err error

	if fromDateStr != "" && toDateStr != "" {
		fromDate, parseErr := time.Parse("2006-01-02", fromDateStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from-date format, expected YYYY-MM-DD"})
			return
		}
		toDate, parseErr := time.Parse("2006-01-02", toDateStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to-date format, expected YYYY-MM-DD"})
			return
		}
		if toDate.Before(fromDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "to-date must not be before from-date"})
			return
		}
		items, err = ctrl.repo.GetOccupationGroupsByUnitGroupForDateRange(id, fromDate, toDate)
	} else {
		items, err = ctrl.repo.GetOccupationGroupsByUnitGroup(id)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve occupation groups",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"unit_group_id":     id,
		"count":             len(items),
		"occupation_groups": items,
	})
}

// Occupations CRUDS
// Major Group CRUD
func (ctrl *JobController) CreateMajorGroupHandler(c *gin.Context) {
	var item models.MajorGroup
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}
	if err := ctrl.repo.CreateMajorGroup(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create major group", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (ctrl *JobController) GetAllMajorGroupsHandler(c *gin.Context) {
	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	var items []models.MajorGroup
	var err error

	if fromDateStr != "" && toDateStr != "" {
		fromDate, parseErr := time.Parse("2006-01-02", fromDateStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from-date format, expected YYYY-MM-DD"})
			return
		}
		toDate, parseErr := time.Parse("2006-01-02", toDateStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to-date format, expected YYYY-MM-DD"})
			return
		}
		if toDate.Before(fromDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "to-date must not be before from-date"})
			return
		}
		items, err = ctrl.repo.GetAllMajorGroupsForDateRange(fromDate, toDate)
	} else {
		items, err = ctrl.repo.GetAllMajorGroups()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve major groups",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":        len(items),
		"major_groups": items,
	})
}

func (ctrl *JobController) GetMajorGroupByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	item, err := ctrl.repo.GetMajorGroupByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Major group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch major group", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) UpdateMajorGroupHandler(c *gin.Context) {
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
	item, err := ctrl.repo.UpdateMajorGroup(id, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update major group", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) DeleteMajorGroupHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	if err := ctrl.repo.DeleteMajorGroup(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete major group", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Major group deleted successfully"})
}

// Sub Major Group CRUD
func (ctrl *JobController) CreateSubMajorGroupHandler(c *gin.Context) {
	var item models.SubMajorGroup
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}
	if err := ctrl.repo.CreateSubMajorGroup(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create sub major group", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (ctrl *JobController) GetAllSubMajorGroupsHandler(c *gin.Context) {
	items, err := ctrl.repo.GetAllSubMajorGroups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sub major groups", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": len(items), "sub_major_groups": items})
}

func (ctrl *JobController) GetSubMajorGroupByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	item, err := ctrl.repo.GetSubMajorGroupByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Sub major group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sub major group", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) UpdateSubMajorGroupHandler(c *gin.Context) {
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
	item, err := ctrl.repo.UpdateSubMajorGroup(id, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update sub major group", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) DeleteSubMajorGroupHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	if err := ctrl.repo.DeleteSubMajorGroup(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete sub major group", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sub major group deleted successfully"})
}

// Minor Group CRUD
func (ctrl *JobController) CreateMinorGroupHandler(c *gin.Context) {
	var item models.MinorGroup
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}
	if err := ctrl.repo.CreateMinorGroup(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create minor group", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (ctrl *JobController) GetAllMinorGroupsHandler(c *gin.Context) {
	items, err := ctrl.repo.GetAllMinorGroups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch minor groups", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": len(items), "minor_groups": items})
}

func (ctrl *JobController) GetMinorGroupByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	item, err := ctrl.repo.GetMinorGroupByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Minor group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch minor group", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) UpdateMinorGroupHandler(c *gin.Context) {
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
	item, err := ctrl.repo.UpdateMinorGroup(id, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update minor group", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) DeleteMinorGroupHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	if err := ctrl.repo.DeleteMinorGroup(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete minor group", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Minor group deleted successfully"})
}

// Unit Group CRUD
func (ctrl *JobController) CreateUnitGroupHandler(c *gin.Context) {
	var item models.UnitGroup
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}
	if err := ctrl.repo.CreateUnitGroup(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create unit group", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (ctrl *JobController) GetAllUnitGroupsHandler(c *gin.Context) {
	items, err := ctrl.repo.GetAllUnitGroups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch unit groups", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": len(items), "unit_groups": items})
}

func (ctrl *JobController) GetUnitGroupByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	item, err := ctrl.repo.GetUnitGroupByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Unit group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch unit group", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) UpdateUnitGroupHandler(c *gin.Context) {
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
	item, err := ctrl.repo.UpdateUnitGroup(id, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update unit group", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) DeleteUnitGroupHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	if err := ctrl.repo.DeleteUnitGroup(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete unit group", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Unit group deleted successfully"})
}

// Occupation Group CRUD
func (ctrl *JobController) CreateOccupationGroupHandler(c *gin.Context) {
	var item models.OccupationGroup
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}
	if err := ctrl.repo.CreateOccupationGroup(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create occupation group", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (ctrl *JobController) GetAllOccupationGroupsHandler(c *gin.Context) {
	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			if v > 100 {
				v = 100
			}
			limit = v
		}
	}
	if o := c.Query("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	items, total, err := ctrl.repo.GetAllOccupationGroups(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch industry subclasses",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":             len(items),
		"total":             total,
		"limit":             limit,
		"offset":            offset,
		"occupation_groups": items,
	})
}

func (ctrl *JobController) GetOccupationGroupByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	item, err := ctrl.repo.GetOccupationGroupByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Occupation group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch occupation group", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) UpdateOccupationGroupHandler(c *gin.Context) {
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
	item, err := ctrl.repo.UpdateOccupationGroup(id, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update occupation group", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) DeleteOccupationGroupHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	if err := ctrl.repo.DeleteOccupationGroup(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete occupation group", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Occupation group deleted successfully"})
}
