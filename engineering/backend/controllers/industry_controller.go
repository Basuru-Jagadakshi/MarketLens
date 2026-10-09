package controllers

import (
	"errors"
	"fmt"
	"marketlens-go-backend/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Industry levels by parent ids
// This function returns industry divisions under an industry sector, currently-active by default,
// or active during a given date range if from-date and to-date query params are provided
func (ctrl *JobController) GetIndustryDivisionsByIndustrySectorHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid industry sector id parameter"})
		return
	}

	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	var items []models.IndustryDivision
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
		items, err = ctrl.repo.GetIndustryDivisionsByIndustrySectorForDateRange(id, fromDate, toDate)
	} else {
		items, err = ctrl.repo.GetIndustryDivisionsByIndustrySector(id)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve industry divisions",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"industry_sector_id": id,
		"count":              len(items),
		"industry_divisions": items,
	})
}

// This function returns industry groups under an industry division, currently-active by default,
// or active during a given date range if from-date and to-date query params are provided
func (ctrl *JobController) GetIndustryGroupsByIndustryDivisionHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid industry division id parameter"})
		return
	}

	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	var items []models.IndustryGroup
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
		items, err = ctrl.repo.GetIndustryGroupsByIndustryDivisionForDateRange(id, fromDate, toDate)
	} else {
		items, err = ctrl.repo.GetIndustryGroupsByIndustryDivision(id)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve industry groups",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"industry_division_id": id,
		"count":                len(items),
		"industry_groups":      items,
	})
}

// This function returns industry classes under an industry group, currently-active by default,
// or active during a given date range if from-date and to-date query params are provided
func (ctrl *JobController) GetIndustryClassesByIndustryGroupHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid industry group id parameter"})
		return
	}

	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	var items []models.IndustryClass
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
		items, err = ctrl.repo.GetIndustryClassesByIndustryGroupForDateRange(id, fromDate, toDate)
	} else {
		items, err = ctrl.repo.GetIndustryClassesByIndustryGroup(id)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve industry classes",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"industry_group_id": id,
		"count":             len(items),
		"industry_classes":  items,
	})
}

// This function returns industry subclasses under an industry class, currently-active by default,
// or active during a given date range if from-date and to-date query params are provided
func (ctrl *JobController) GetIndustrySubclassesByIndustryClassHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid industry class id parameter"})
		return
	}

	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	var items []models.IndustrySubclass
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
		items, err = ctrl.repo.GetIndustrySubclassesByIndustryClassForDateRange(id, fromDate, toDate)
	} else {
		items, err = ctrl.repo.GetIndustrySubclassesByIndustryClass(id)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve industry subclasses",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"industry_class_id":   id,
		"count":               len(items),
		"industry_subclasses": items,
	})
}

// Industry Sectors CRUDS
// Industry Sctor CRUD
func (ctrl *JobController) CreateIndustrySectorHandler(c *gin.Context) {
	var item models.IndustrySector
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}
	if err := ctrl.repo.CreateIndustrySector(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create industry sector", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// This function returns all currently-active industry sectors, or ones active during a
// given date range if from-date and to-date query params are provided
func (ctrl *JobController) GetAllIndustrySectorsHandler(c *gin.Context) {
	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	var items []models.IndustrySector
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
		items, err = ctrl.repo.GetAllIndustrySectorsForDateRange(fromDate, toDate)
	} else {
		items, err = ctrl.repo.GetAllIndustrySectors()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch industry sectors", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": len(items), "industry_sectors": items})
}

func (ctrl *JobController) GetIndustrySectorByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	item, err := ctrl.repo.GetIndustrySectorByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Industry sector not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch industry sector", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) UpdateIndustrySectorHandler(c *gin.Context) {
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
	item, err := ctrl.repo.UpdateIndustrySector(id, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update industry sector", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) DeleteIndustrySectorHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	if err := ctrl.repo.DeleteIndustrySector(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete industry sector", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Industry sector deleted successfully"})
}

// Industry Division CRUD
func (ctrl *JobController) CreateIndustryDivisionHandler(c *gin.Context) {
	var item models.IndustryDivision
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}
	if err := ctrl.repo.CreateIndustryDivision(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create industry division", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (ctrl *JobController) GetAllIndustryDivisionsHandler(c *gin.Context) {
	items, err := ctrl.repo.GetAllIndustryDivisions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch industry divisions", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": len(items), "industry_divisions": items})
}

func (ctrl *JobController) GetIndustryDivisionByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	item, err := ctrl.repo.GetIndustryDivisionByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Industry division not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch industry division", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) UpdateIndustryDivisionHandler(c *gin.Context) {
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
	item, err := ctrl.repo.UpdateIndustryDivision(id, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update industry division", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) DeleteIndustryDivisionHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	if err := ctrl.repo.DeleteIndustryDivision(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete industry division", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Industry division deleted successfully"})
}

// Industry Group CRUD
func (ctrl *JobController) CreateIndustryGroupHandler(c *gin.Context) {
	var item models.IndustryGroup
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}
	if err := ctrl.repo.CreateIndustryGroup(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create industry group", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (ctrl *JobController) GetAllIndustryGroupsHandler(c *gin.Context) {
	items, err := ctrl.repo.GetAllIndustryGroups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch industry groups", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": len(items), "industry_groups": items})
}

func (ctrl *JobController) GetIndustryGroupByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	item, err := ctrl.repo.GetIndustryGroupByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Industry group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch industry group", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) UpdateIndustryGroupHandler(c *gin.Context) {
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
	item, err := ctrl.repo.UpdateIndustryGroup(id, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update industry group", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) DeleteIndustryGroupHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	if err := ctrl.repo.DeleteIndustryGroup(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete industry group", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Industry group deleted successfully"})
}

// Industry Class CRUD
func (ctrl *JobController) CreateIndustryClassHandler(c *gin.Context) {
	var item models.IndustryClass
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}
	if err := ctrl.repo.CreateIndustryClass(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create industry class", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (ctrl *JobController) GetAllIndustryClassesHandler(c *gin.Context) {
	items, err := ctrl.repo.GetAllIndustryClasses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch industry classes", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": len(items), "industry_classes": items})
}

func (ctrl *JobController) GetIndustryClassByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	item, err := ctrl.repo.GetIndustryClassByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Industry class not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch industry class", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) UpdateIndustryClassHandler(c *gin.Context) {
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
	item, err := ctrl.repo.UpdateIndustryClass(id, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update industry class", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) DeleteIndustryClassHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	if err := ctrl.repo.DeleteIndustryClass(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete industry class", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Industry class deleted successfully"})
}

// Industry Sub Class CRUD
func (ctrl *JobController) CreateIndustrySubclassHandler(c *gin.Context) {
	var item models.IndustrySubclass
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}
	if err := ctrl.repo.CreateIndustrySubclass(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create industry subclass", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (ctrl *JobController) GetAllIndustrySubclassesHandler(c *gin.Context) {
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

	items, total, err := ctrl.repo.GetAllIndustrySubclasses(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch industry subclasses",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":               len(items),
		"total":               total,
		"limit":               limit,
		"offset":              offset,
		"industry_subclasses": items,
	})
}

func (ctrl *JobController) GetIndustrySubclassByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	item, err := ctrl.repo.GetIndustrySubclassByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Industry subclass not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch industry subclass", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) UpdateIndustrySubclassHandler(c *gin.Context) {
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
	item, err := ctrl.repo.UpdateIndustrySubclass(id, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update industry subclass", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *JobController) DeleteIndustrySubclassHandler(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}
	if err := ctrl.repo.DeleteIndustrySubclass(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete industry subclass", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Industry subclass deleted successfully"})
}
