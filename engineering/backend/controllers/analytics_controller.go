package controllers

import (
	"fmt"
	"marketlens-go-backend/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// This function returns the top hiring employers for a given occupation hierarchy
// level and id, filtered by date range.
func (ctrl *JobController) GetTopHiringEmployersByOccupationLevelHandler(c *gin.Context) {
	level := c.Param("level")
	idStr := c.Param("id")
	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	if !validOccupationLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid level for occupation",
			"valid_levels": []string{
				"major-group", "sub-major-group", "minor-group", "unit-group", "occupation-group",
			},
		})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter, must be a positive integer"})
		return
	}

	if fromDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from-date query parameter is required"})
		return
	}
	if toDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date query parameter is required"})
		return
	}

	const layout = "2006-01-02"

	fromDate, err := time.Parse(layout, fromDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid from-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	toDate, err := time.Parse(layout, toDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid to-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	if toDate.Before(fromDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date must not be before from-date"})
		return
	}

	results, err := ctrl.repo.GetTopHiringEmployersByOccupationLevel(level, uint(id), fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve top hiring employers",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"level":     level,
		"id":        id,
		"from_date": fromDateStr,
		"to_date":   toDateStr,
		"count":     len(results),
		"employers": results,
	})
}

// This function returns all skills (paginated) for a given occupation hierarchy
// level and id, filtered by date range.
func (ctrl *JobController) GetAllSkillsByOccupationLevelHandler(c *gin.Context) {
	level := c.Param("level")
	idStr := c.Param("id")
	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	if !validOccupationLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid level for occupation",
			"valid_levels": []string{
				"major-group", "sub-major-group", "minor-group", "unit-group", "occupation-group",
			},
		})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter, must be a positive integer"})
		return
	}

	if fromDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from-date query parameter is required"})
		return
	}
	if toDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date query parameter is required"})
		return
	}

	const layout = "2006-01-02"

	fromDate, err := time.Parse(layout, fromDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid from-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	toDate, err := time.Parse(layout, toDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid to-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	if toDate.Before(fromDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date must not be before from-date"})
		return
	}

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

	results, total, err := ctrl.repo.GetAllSkillsByOccupationLevel(level, uint(id), fromDate, toDate, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve skills",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"level":     level,
		"id":        id,
		"from_date": fromDateStr,
		"to_date":   toDateStr,
		"count":     len(results),
		"total":     total,
		"limit":     limit,
		"offset":    offset,
		"skills":    results,
	})
}

// This function returns the top 15 skills for a given occupation hierarchy
// level and id, filtered by date range.
func (ctrl *JobController) GetTop15SkillsByOccupationLevelHandler(c *gin.Context) {
	level := c.Param("level")
	idStr := c.Param("id")
	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	if !validOccupationLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid level for occupation",
			"valid_levels": []string{
				"major-group", "sub-major-group", "minor-group", "unit-group", "occupation-group",
			},
		})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter, must be a positive integer"})
		return
	}

	if fromDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from-date query parameter is required"})
		return
	}
	if toDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date query parameter is required"})
		return
	}

	const layout = "2006-01-02"

	fromDate, err := time.Parse(layout, fromDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid from-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	toDate, err := time.Parse(layout, toDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid to-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	if toDate.Before(fromDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date must not be before from-date"})
		return
	}

	results, err := ctrl.repo.GetTop15SkillsByOccupationLevel(level, uint(id), fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve top 15 skills",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"level":     level,
		"id":        id,
		"from_date": fromDateStr,
		"to_date":   toDateStr,
		"count":     len(results),
		"skills":    results,
	})
}

// This function returns job type breakdown for a given occupation/industry
// hierarchy level and id, filtered by date range.
func (ctrl *JobController) GetJobTypeByLevelHandler(c *gin.Context) {
	standard := c.Param("standard")
	level := c.Param("level")
	idStr := c.Param("id")
	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	if standard != "occupation" && standard != "industry" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid standard, must be 'occupation' or 'industry'"})
		return
	}

	if standard == "occupation" && !validOccupationLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid level for standard 'occupation'",
			"valid_levels": []string{
				"major-group", "sub-major-group", "minor-group", "unit-group", "occupation-group",
			},
		})
		return
	}

	if standard == "industry" && !validIndustryLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid level for standard 'industry'",
			"valid_levels": []string{
				"industry-sector", "industry-division", "industry-group", "industry-class", "industry-subclass",
			},
		})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter, must be a positive integer"})
		return
	}

	if fromDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from-date query parameter is required"})
		return
	}
	if toDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date query parameter is required"})
		return
	}

	const layout = "2006-01-02"

	fromDate, err := time.Parse(layout, fromDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid from-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	toDate, err := time.Parse(layout, toDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid to-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	if toDate.Before(fromDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date must not be before from-date"})
		return
	}

	results, err := ctrl.repo.GetJobTypeByLevel(standard, level, uint(id), fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve job type breakdown",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"standard":  standard,
		"level":     level,
		"id":        id,
		"from_date": fromDateStr,
		"to_date":   toDateStr,
		"count":     len(results),
		"job_types": results,
	})
}

// This function returns remote vs on-site job counts for a given occupation/industry
// hierarchy level and id, filtered by date range.
func (ctrl *JobController) GetRemoteOnSiteHybridByLevelHandler(c *gin.Context) {
	standard := c.Param("standard")
	level := c.Param("level")
	idStr := c.Param("id")
	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	if standard != "occupation" && standard != "industry" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid standard, must be 'occupation' or 'industry'"})
		return
	}

	if standard == "occupation" && !validOccupationLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid level for standard 'occupation'",
			"valid_levels": []string{
				"major-group", "sub-major-group", "minor-group", "unit-group", "occupation-group",
			},
		})
		return
	}

	if standard == "industry" && !validIndustryLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid level for standard 'industry'",
			"valid_levels": []string{
				"industry-sector", "industry-division", "industry-group", "industry-class", "industry-subclass",
			},
		})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter, must be a positive integer"})
		return
	}

	if fromDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from-date query parameter is required"})
		return
	}
	if toDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date query parameter is required"})
		return
	}

	const layout = "2006-01-02"

	fromDate, err := time.Parse(layout, fromDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid from-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	toDate, err := time.Parse(layout, toDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid to-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	if toDate.Before(fromDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date must not be before from-date"})
		return
	}

	result, err := ctrl.repo.GetRemoteOnSiteHybridByLevel(standard, level, uint(id), fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve remote/on-site breakdown",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"standard":                   standard,
		"level":                      level,
		"id":                         id,
		"from_date":                  fromDateStr,
		"to_date":                    toDateStr,
		"remote_vs_onsite_vs_hybrid": result,
	})
}

// This function returns vocational education breakdown for a given occupation/industry
// hierarchy level and id, filtered by date range.
func (ctrl *JobController) GetVocationalEducationByLevelHandler(c *gin.Context) {
	standard := c.Param("standard")
	level := c.Param("level")
	idStr := c.Param("id")
	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	if standard != "occupation" && standard != "industry" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid standard, must be 'occupation' or 'industry'"})
		return
	}

	if standard == "occupation" && !validOccupationLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid level for standard 'occupation'",
			"valid_levels": []string{
				"major-group", "sub-major-group", "minor-group", "unit-group", "occupation-group",
			},
		})
		return
	}

	if standard == "industry" && !validIndustryLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid level for standard 'industry'",
			"valid_levels": []string{
				"industry-sector", "industry-division", "industry-group", "industry-class", "industry-subclass",
			},
		})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter, must be a positive integer"})
		return
	}

	if fromDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from-date query parameter is required"})
		return
	}
	if toDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date query parameter is required"})
		return
	}

	const layout = "2006-01-02"

	fromDate, err := time.Parse(layout, fromDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid from-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	toDate, err := time.Parse(layout, toDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid to-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	if toDate.Before(fromDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date must not be before from-date"})
		return
	}

	results, err := ctrl.repo.GetVocationalEducationByLevel(standard, level, uint(id), fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve vocational education breakdown",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"standard":              standard,
		"level":                 level,
		"id":                    id,
		"from_date":             fromDateStr,
		"to_date":               toDateStr,
		"count":                 len(results),
		"vocational_educations": results,
	})
}

// This function returns gender breakdown for a given occupation/industry
// hierarchy level and id, filtered by date range.
func (ctrl *JobController) GetGenderByLevelHandler(c *gin.Context) {
	standard := c.Param("standard")
	level := c.Param("level")
	idStr := c.Param("id")
	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	if standard != "occupation" && standard != "industry" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid standard, must be 'occupation' or 'industry'"})
		return
	}

	if standard == "occupation" && !validOccupationLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid level for standard 'occupation'",
			"valid_levels": []string{
				"major-group", "sub-major-group", "minor-group", "unit-group", "occupation-group",
			},
		})
		return
	}

	if standard == "industry" && !validIndustryLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid level for standard 'industry'",
			"valid_levels": []string{
				"industry-sector", "industry-division", "industry-group", "industry-class", "industry-subclass",
			},
		})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter, must be a positive integer"})
		return
	}

	if fromDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from-date query parameter is required"})
		return
	}
	if toDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date query parameter is required"})
		return
	}

	const layout = "2006-01-02"

	fromDate, err := time.Parse(layout, fromDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid from-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	toDate, err := time.Parse(layout, toDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid to-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	if toDate.Before(fromDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date must not be before from-date"})
		return
	}

	results, err := ctrl.repo.GetGenderByLevel(standard, level, uint(id), fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve gender breakdown",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"standard":  standard,
		"level":     level,
		"id":        id,
		"from_date": fromDateStr,
		"to_date":   toDateStr,
		"count":     len(results),
		"genders":   results,
	})
}

// This function returns formality breakdown for a given occupation/industry
// hierarchy level and id, filtered by date range.
func (ctrl *JobController) GetFormalityByLevelHandler(c *gin.Context) {
	standard := c.Param("standard")
	level := c.Param("level")
	idStr := c.Param("id")
	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	if standard != "occupation" && standard != "industry" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid standard, must be 'occupation' or 'industry'"})
		return
	}

	if standard == "occupation" && !validOccupationLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid level for standard 'occupation'",
			"valid_levels": []string{
				"major-group", "sub-major-group", "minor-group", "unit-group", "occupation-group",
			},
		})
		return
	}

	if standard == "industry" && !validIndustryLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid level for standard 'industry'",
			"valid_levels": []string{
				"industry-sector", "industry-division", "industry-group", "industry-class", "industry-subclass",
			},
		})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter, must be a positive integer"})
		return
	}

	if fromDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from-date query parameter is required"})
		return
	}
	if toDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date query parameter is required"})
		return
	}

	const layout = "2006-01-02"

	fromDate, err := time.Parse(layout, fromDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid from-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	toDate, err := time.Parse(layout, toDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid to-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	if toDate.Before(fromDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date must not be before from-date"})
		return
	}

	results, err := ctrl.repo.GetFormalityByLevel(standard, level, uint(id), fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve formality breakdown",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"standard":    standard,
		"level":       level,
		"id":          id,
		"from_date":   fromDateStr,
		"to_date":     toDateStr,
		"count":       len(results),
		"formalities": results,
	})
}

// This function returns education level breakdown for a given occupation/industry
// hierarchy level and id, filtered by date range.
func (ctrl *JobController) GetEducationLevelByLevelHandler(c *gin.Context) {
	standard := c.Param("standard")
	level := c.Param("level")
	idStr := c.Param("id")
	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	if standard != "occupation" && standard != "industry" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid standard, must be 'occupation' or 'industry'"})
		return
	}

	if standard == "occupation" && !validOccupationLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid level for standard 'occupation'",
			"valid_levels": []string{
				"major-group", "sub-major-group", "minor-group", "unit-group", "occupation-group",
			},
		})
		return
	}

	if standard == "industry" && !validIndustryLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid level for standard 'industry'",
			"valid_levels": []string{
				"industry-sector", "industry-division", "industry-group", "industry-class", "industry-subclass",
			},
		})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter, must be a positive integer"})
		return
	}

	if fromDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from-date query parameter is required"})
		return
	}
	if toDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date query parameter is required"})
		return
	}

	const layout = "2006-01-02"

	fromDate, err := time.Parse(layout, fromDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid from-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	toDate, err := time.Parse(layout, toDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid to-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	if toDate.Before(fromDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date must not be before from-date"})
		return
	}

	results, err := ctrl.repo.GetEducationLevelByLevel(standard, level, uint(id), fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve education level breakdown",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"standard":         standard,
		"level":            level,
		"id":               id,
		"from_date":        fromDateStr,
		"to_date":          toDateStr,
		"count":            len(results),
		"education_levels": results,
	})
}

// This function returns province breakdown for a given occupation/industry
// hierarchy level and id, filtered by date range.
func (ctrl *JobController) GetProvinceByLevelHandler(c *gin.Context) {
	standard := c.Param("standard")
	level := c.Param("level")
	idStr := c.Param("id")
	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	if standard != "occupation" && standard != "industry" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid standard, must be 'occupation' or 'industry'"})
		return
	}

	if standard == "occupation" && !validOccupationLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid level for standard 'occupation'",
			"valid_levels": []string{
				"major-group", "sub-major-group", "minor-group", "unit-group", "occupation-group",
			},
		})
		return
	}

	if standard == "industry" && !validIndustryLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid level for standard 'industry'",
			"valid_levels": []string{
				"industry-sector", "industry-division", "industry-group", "industry-class", "industry-subclass",
			},
		})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter, must be a positive integer"})
		return
	}

	if fromDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from-date query parameter is required"})
		return
	}
	if toDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date query parameter is required"})
		return
	}

	const layout = "2006-01-02"

	fromDate, err := time.Parse(layout, fromDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid from-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	toDate, err := time.Parse(layout, toDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid to-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	if toDate.Before(fromDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date must not be before from-date"})
		return
	}

	results, err := ctrl.repo.GetProvinceByLevel(standard, level, uint(id), fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve province breakdown",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"standard":  standard,
		"level":     level,
		"id":        id,
		"from_date": fromDateStr,
		"to_date":   toDateStr,
		"count":     len(results),
		"provinces": results,
	})
}

// This function returns experience-level breakdown for a given occupation/industry
// hierarchy level and id, filtered by date range.
func (ctrl *JobController) GetExperienceByLevelHandler(c *gin.Context) {
	standard := c.Param("standard")
	level := c.Param("level")
	idStr := c.Param("id")
	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	if standard != "occupation" && standard != "industry" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid standard, must be 'occupation' or 'industry'"})
		return
	}

	if standard == "occupation" && !validOccupationLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid level for standard 'occupation'",
			"valid_levels": []string{
				"major-group", "sub-major-group", "minor-group", "unit-group", "occupation-group",
			},
		})
		return
	}

	if standard == "industry" && !validIndustryLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid level for standard 'industry'",
			"valid_levels": []string{
				"industry-sector", "industry-division", "industry-group", "industry-class", "industry-subclass",
			},
		})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter, must be a positive integer"})
		return
	}

	if fromDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from-date query parameter is required"})
		return
	}
	if toDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date query parameter is required"})
		return
	}

	const layout = "2006-01-02"

	fromDate, err := time.Parse(layout, fromDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid from-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	toDate, err := time.Parse(layout, toDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid to-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	if toDate.Before(fromDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date must not be before from-date"})
		return
	}

	results, err := ctrl.repo.GetExperienceByLevel(standard, level, uint(id), fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve experience breakdown",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"standard":    standard,
		"level":       level,
		"id":          id,
		"from_date":   fromDateStr,
		"to_date":     toDateStr,
		"count":       len(results),
		"experiences": results,
	})
}

// This function returns employment sector breakdown for a given occupation/industry
// hierarchy level and id, filtered by date range.
func (ctrl *JobController) GetEmploymentSectorByLevelHandler(c *gin.Context) {
	standard := c.Param("standard")
	level := c.Param("level")
	idStr := c.Param("id")
	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	if standard != "occupation" && standard != "industry" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid standard, must be 'occupation' or 'industry'"})
		return
	}

	if standard == "occupation" && !validOccupationLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid level for standard 'occupation'",
			"valid_levels": []string{
				"major-group", "sub-major-group", "minor-group", "unit-group", "occupation-group",
			},
		})
		return
	}

	if standard == "industry" && !validIndustryLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid level for standard 'industry'",
			"valid_levels": []string{
				"industry-sector", "industry-division", "industry-group", "industry-class", "industry-subclass",
			},
		})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter, must be a positive integer"})
		return
	}

	if fromDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from-date query parameter is required"})
		return
	}
	if toDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date query parameter is required"})
		return
	}

	const layout = "2006-01-02"

	fromDate, err := time.Parse(layout, fromDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid from-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	toDate, err := time.Parse(layout, toDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid to-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	if toDate.Before(fromDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date must not be before from-date"})
		return
	}

	results, err := ctrl.repo.GetEmploymentSectorByLevel(standard, level, uint(id), fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve employment sector breakdown",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"standard":           standard,
		"level":              level,
		"id":                 id,
		"from_date":          fromDateStr,
		"to_date":            toDateStr,
		"count":              len(results),
		"employment_sectors": results,
	})
}

// leafLevels marks the bottom of each hierarchy, which has no children to return
var leafLevels = map[string]bool{
	"occupation-group":  true,
	"industry-subclass": true,
}

// This function returns the immediate children of a given occupation/industry hierarchy
// level and id, each with its aggregated job count for a given date range.
func (ctrl *JobController) GetLevelChildrenHandler(c *gin.Context) {
	standard := c.Param("standard")
	level := c.Param("level")
	idStr := c.Param("id")
	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	if standard != "occupation" && standard != "industry" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid standard, must be 'occupation' or 'industry'"})
		return
	}

	if standard == "occupation" && !validOccupationLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid level for standard 'occupation'",
			"valid_levels": []string{
				"major-group", "sub-major-group", "minor-group", "unit-group", "occupation-group",
			},
		})
		return
	}

	if standard == "industry" && !validIndustryLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid level for standard 'industry'",
			"valid_levels": []string{
				"industry-sector", "industry-division", "industry-group", "industry-class", "industry-subclass",
			},
		})
		return
	}

	if leafLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("'%s' is a leaf level and has no children", level)})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter, must be a positive integer"})
		return
	}

	if fromDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from-date query parameter is required"})
		return
	}
	if toDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date query parameter is required"})
		return
	}

	const layout = "2006-01-02"

	fromDate, err := time.Parse(layout, fromDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid from-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	toDate, err := time.Parse(layout, toDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid to-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	if toDate.Before(fromDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date must not be before from-date"})
		return
	}

	results, childLevel, err := ctrl.repo.GetLevelChildren(standard, level, uint(id), fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve child level data",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"standard":    standard,
		"level":       level,
		"id":          id,
		"child_level": childLevel,
		"from_date":   fromDateStr,
		"to_date":     toDateStr,
		"count":       len(results),
		"children":    results,
	})
}

var validOccupationLevels = map[string]bool{
	"major-group":      true,
	"sub-major-group":  true,
	"minor-group":      true,
	"unit-group":       true,
	"occupation-group": true,
}

var validIndustryLevels = map[string]bool{
	"industry-sector":   true,
	"industry-division": true,
	"industry-group":    true,
	"industry-class":    true,
	"industry-subclass": true,
}

// This function returns the total job count for a given occupation/industry hierarchy level and id,
// filtered by an optional date range.
func (ctrl *JobController) GetTotalJobCountByLevelHandler(c *gin.Context) {
	standard := c.Param("standard")
	level := c.Param("level")
	idStr := c.Param("id")
	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	if standard != "occupation" && standard != "industry" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid standard, must be 'occupation' or 'industry'"})
		return
	}

	if standard == "occupation" && !validOccupationLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid level for standard 'occupation'",
			"valid_levels": []string{
				"major-group", "sub-major-group", "minor-group", "unit-group", "occupation-group",
			},
		})
		return
	}

	if standard == "industry" && !validIndustryLevels[level] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid level for standard 'industry'",
			"valid_levels": []string{
				"industry-sector", "industry-division", "industry-group", "industry-class", "industry-subclass",
			},
		})
		return
	}

	var id uint64
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter, must be a positive integer"})
		return
	}

	if fromDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from-date query parameter is required"})
		return
	}
	if toDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date query parameter is required"})
		return
	}

	const layout = "2006-01-02"

	fromDate, err := time.Parse(layout, fromDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid from-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	toDate, err := time.Parse(layout, toDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid to-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	if toDate.Before(fromDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date must not be before from-date"})
		return
	}

	total, err := ctrl.repo.GetTotalJobCountByLevel(standard, level, uint(id), fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve total job count",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"standard":        standard,
		"level":           level,
		"id":              id,
		"from_date":       fromDateStr,
		"to_date":         toDateStr,
		"total_job_count": total,
	})
}

// This function returns vacancy counts grouped by industry sector, for a given date range
func (ctrl *JobController) GetIndustryJobCountByDateRangeHandler(c *gin.Context) {
	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	if fromDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from-date query parameter is required"})
		return
	}
	if toDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date query parameter is required"})
		return
	}

	const layout = "2006-01-02"

	fromDate, err := time.Parse(layout, fromDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid from-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	toDate, err := time.Parse(layout, toDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid to-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	if toDate.Before(fromDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date must not be before from-date"})
		return
	}

	results, err := ctrl.repo.GetIndustryJobCountByDateRange(fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve industry job counts",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"from_date":  fromDateStr,
		"to_date":    toDateStr,
		"count":      len(results),
		"industries": results,
	})
}

// This function returns vacancy counts grouped by major group (occupation), for a given date range
func (ctrl *JobController) GetOccupationJobCountByDateRangeHandler(c *gin.Context) {
	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	if fromDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from-date query parameter is required"})
		return
	}
	if toDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date query parameter is required"})
		return
	}

	const layout = "2006-01-02"

	fromDate, err := time.Parse(layout, fromDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid from-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	toDate, err := time.Parse(layout, toDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid to-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	if toDate.Before(fromDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date must not be before from-date"})
		return
	}

	results, err := ctrl.repo.GetOccupationJobCountByDateRange(fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve occupation job counts",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"from_date":   fromDateStr,
		"to_date":     toDateStr,
		"count":       len(results),
		"occupations": results,
	})
}

// This function returns the total number of vacancies posted within a given date range
func (ctrl *JobController) GetTotalVacancyCountHandler(c *gin.Context) {
	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	if fromDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from-date query parameter is required"})
		return
	}
	if toDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date query parameter is required"})
		return
	}

	const layout = "2006-01-02"

	fromDate, err := time.Parse(layout, fromDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid from-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	toDate, err := time.Parse(layout, toDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid to-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	if toDate.Before(fromDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date must not be before from-date"})
		return
	}

	total, err := ctrl.repo.GetTotalVacancyCount(fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve total vacancy count",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.TotalVacancyCount{
		FromDate:       fromDateStr,
		ToDate:         toDateStr,
		TotalVacancies: total,
	})
}

// This function returns vacancy counts over a date range with granularity that adapts
// to the range's length: 7-day set under 60 days, monthly otherwise
func (ctrl *JobController) GetVacancyTrendHandler(c *gin.Context) {
	fromDateStr := c.Query("from-date")
	toDateStr := c.Query("to-date")

	if fromDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from-date query parameter is required"})
		return
	}
	if toDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date query parameter is required"})
		return
	}

	const layout = "2006-01-02"

	fromDate, err := time.Parse(layout, fromDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid from-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	toDate, err := time.Parse(layout, toDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid to-date format, expected YYYY-MM-DD",
			"details": err.Error(),
		})
		return
	}

	if toDate.Before(fromDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to-date must not be before from-date"})
		return
	}

	diffDays := int(toDate.Sub(fromDate).Hours() / 24)

	var (
		results     []models.VacancyTrendPoint
		granularity string
		repoErr     error
	)

	if diffDays < 60 {
		granularity = "weekly"
		results, repoErr = ctrl.repo.GetVacancyTrendDaily(fromDate, toDate)
	} else {
		granularity = "monthly"
		results, repoErr = ctrl.repo.GetVacancyTrendMonthly(fromDate, toDate)
	}

	if repoErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve vacancy trend data",
			"details": repoErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"from_date":   fromDateStr,
		"to_date":     toDateStr,
		"granularity": granularity,
		"count":       len(results),
		"data":        results,
	})
}
