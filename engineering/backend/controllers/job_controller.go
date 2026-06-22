package controllers

import (
	"fmt"
	"marketlens-go-backend/models"
	"marketlens-go-backend/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)



type JobController struct {
	repo *repositories.JobRepository
}

func NewJobController(repo *repositories.JobRepository) *JobController {
	return &JobController{repo: repo}
}

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

func (ctrl *JobController) GetJobsByBucketKeysHandler(c *gin.Context) {
	var req struct {
		BucketKeys []string `json:"bucket_keys" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request structure payload mapping",
			"details": err.Error(),
		})
		return
	}

	matchingJobs, err := ctrl.repo.GetJobsByBucketKeys(req.BucketKeys)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "LSH bucket radar execution sequence encountered a processing exception",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, matchingJobs)
}

func (ctrl *JobController) BatchSaveJobsHandler(c *gin.Context) {
	var payload models.BatchSavePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body structural mapping",
			"details": err.Error(),
		})
		return
	}

	if len(payload.NewJobs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The new_jobs collection buffer cannot be empty"})
		return
	}

	err := ctrl.repo.BatchSaveNewJobs(payload.NewJobs, payload.LshIndexes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Bulk insertion transaction routine failed execution",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "Successfully persisted unique job block chunk",
		"inserted_records": len(payload.NewJobs),
	})
}

func (ctrl *JobController) BatchUpdateDuplicatesHandler(c *gin.Context) {
	var payload models.BatchUpdatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload configuration mapping",
			"details": err.Error(),
		})
		return
	}

	if len(payload.Duplicates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The duplicates reference buffer cannot be empty"})
		return
	}

	err := ctrl.repo.BatchUpdateDuplicateJobs(payload.Duplicates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Bulk update checkpoint modifications failed execution",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "Successfully refreshed duplicate job keep-alive markers",
		"refreshed_records": len(payload.Duplicates),
	})
}

func (ctrl *JobController) ReconcileStaleVacanciesHandler(c *gin.Context) {
	var payload models.ReconciliationPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing active execution tracking sequence identifier criteria",
			"details": err.Error(),
		})
		return
	}

	closedCount, err := ctrl.repo.ReconcileStaleVacancies(payload.CrawlerRunID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Global snapshot reconciliation sweeping routine failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":               "System-wide stale vacancy cleanup sweep finalized",
		"reconciled_stale_jobs": closedCount,
	})
}


func (ctrl *JobController) DeleteJobHandler(c *gin.Context) {
	
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID format parameter"})
		return
	}

	id, err := ctrl.repo.DeleteJob(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to execute deletion on targeted job profile",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Job post with ID %d has been successfully deleted", id),
	})
}

func (ctrl *JobController) GetActiveJobsHandler(c *gin.Context) {

	// Helper to parse optional uint query params
	parseUintParam := func(key string) (*uint, error) {
		val := c.Query(key)
		if val == "" {
			return nil, nil
		}
		var parsed uint
		if _, err := fmt.Sscanf(val, "%d", &parsed); err != nil {
			return nil, fmt.Errorf("invalid value for query parameter '%s'", key)
		}
		return &parsed, nil
	}

	industryID, err := parseUintParam("industry_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	geoDataID, err := parseUintParam("geo_data_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	jobTypeID, err := parseUintParam("job_type_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	experienceID, err := parseUintParam("experience_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jobs, err := ctrl.repo.GetActiveJobs(industryID, geoDataID, jobTypeID, experienceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve active job listings",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(jobs),
		"jobs":  jobs,
	})
}

func (ctrl *JobController) GetAllIndustriesHandler(c *gin.Context) {
	industries, err := ctrl.repo.GetAllIndustries()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve industries",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"count":      len(industries),
		"industries": industries,
	})
}

func (ctrl *JobController) GetAllExperiencesHandler(c *gin.Context) {
	experiences, err := ctrl.repo.GetAllExperiences()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve experience levels",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"count":       len(experiences),
		"experiences": experiences,
	})
}

func (ctrl *JobController) GetAllProvincesHandler(c *gin.Context) {
	provinces, err := ctrl.repo.GetAllProvinces()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve provinces",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"count":     len(provinces),
		"provinces": provinces,
	})
}

func (ctrl *JobController) GetAllJobTypesHandler(c *gin.Context) {
	jobTypes, err := ctrl.repo.GetAllJobTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve job types",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"count":     len(jobTypes),
		"job_types": jobTypes,
	})
}

func (ctrl *JobController) GetUniqueSkillsCountByIndustryHandler(c *gin.Context) {
	industryID, err := parseIndustryID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count, err := ctrl.repo.GetUniqueSkillsCountByIndustry(industryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve unique skills count",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"industry_id":         industryID,
		"unique_skills_count": count,
	})
}

func (ctrl *JobController) GetMostDemandingSkillByIndustryHandler(c *gin.Context) {
	industryID, err := parseIndustryID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	skill, err := ctrl.repo.GetMostDemandingSkillByIndustry(industryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve most demanding skill",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"industry_id":       industryID,
		"most_in_demand_skill": skill,
	})
}

func (ctrl *JobController) GetTop15SkillsByIndustryHandler(c *gin.Context) {
	industryID, err := parseIndustryID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	skills, err := ctrl.repo.GetTop15SkillsByIndustry(industryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve top 15 skills",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"industry_id": industryID,
		"count":       len(skills),
		"skills":      skills,
	})
}

func (ctrl *JobController) GetAllSkillsByIndustryHandler(c *gin.Context) {
	industryID, err := parseIndustryID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	skills, err := ctrl.repo.GetAllSkillsByIndustry(industryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve all skills",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"industry_id": industryID,
		"count":       len(skills),
		"skills":      skills,
	})
}

func (ctrl *JobController) GetTopHiringEmployersByIndustryHandler(c *gin.Context) {
	industryID, err := parseIndustryID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employers, err := ctrl.repo.GetTopHiringEmployersByIndustry(industryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve top hiring employers",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"industry_id": industryID,
		"count":       len(employers),
		"employers":   employers,
	})
}

// parseIndustryID is a shared helper for all industry-scoped handlers
func parseIndustryID(c *gin.Context) (uint, error) {
	val := c.Query("industry_id")
	if val == "" {
		return 0, fmt.Errorf("industry_id query parameter is required")
	}
	var id uint
	if _, err := fmt.Sscanf(val, "%d", &id); err != nil {
		return 0, fmt.Errorf("invalid industry_id value")
	}
	return id, nil
}

func (ctrl *JobController) GetLastCrawledJobCountHandler(c *gin.Context) {
	count, err := ctrl.repo.GetLastCrawledJobCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve last crawled job count",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"last_crawl_job_count": count,
	})
}

func (ctrl *JobController) GetTimeSinceLastCrawlHandler(c *gin.Context) {
	gap, err := ctrl.repo.GetTimeSinceLastCrawl()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve last crawl time gap",
			"details": err.Error(),
		})
		return
	}

	if gap.LastCrawledAt == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "No completed crawler runs found",
		})
		return
	}

	c.JSON(http.StatusOK, gap)
}

func (ctrl *JobController) GetSourcesWithActiveJobCountHandler(c *gin.Context) {
	sources, err := ctrl.repo.GetSourcesWithActiveJobCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve sources with active job counts",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":   len(sources),
		"sources": sources,
	})
}

func (ctrl *JobController) GetAllCrawlerRunsHandler(c *gin.Context) {
	runs, err := ctrl.repo.GetAllCrawlerRuns()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve crawler runs",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(runs),
		"runs":  runs,
	})
}

func (ctrl *JobController) GetActiveJobCountWithTrendHandler(c *gin.Context) {
	result, err := ctrl.repo.GetActiveJobCountWithTrend()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve active job count with trend",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (ctrl *JobController) GetActiveJobCountByOccupationHandler(c *gin.Context) {
	results, err := ctrl.repo.GetActiveJobCountByOccupation()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve job counts by occupation",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"count":       len(results),
		"occupations": results,
	})
}

func (ctrl *JobController) GetActiveJobCountByIndustryHandler(c *gin.Context) {
	results, err := ctrl.repo.GetActiveJobCountByIndustry()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve job counts by industry",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"count":      len(results),
		"industries": results,
	})
}

func (ctrl *JobController) GetActiveJobCountByExperienceHandler(c *gin.Context) {
	results, err := ctrl.repo.GetActiveJobCountByExperience()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve job counts by experience",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"count":       len(results),
		"experiences": results,
	})
}

func (ctrl *JobController) GetActiveJobCountByEducationLevelHandler(c *gin.Context) {
	results, err := ctrl.repo.GetActiveJobCountByEducationLevel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve job counts by education level",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"count":            len(results),
		"education_levels": results,
	})
}

func (ctrl *JobController) GetRemoteVsOnSiteCountHandler(c *gin.Context) {
	result, err := ctrl.repo.GetRemoteVsOnSiteCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve remote vs on-site counts",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (ctrl *JobController) GetActiveJobCountByJobTypeHandler(c *gin.Context) {
	results, err := ctrl.repo.GetActiveJobCountByJobType()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve job counts by job type",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"count":     len(results),
		"job_types": results,
	})
}

func parseOccupationID(c *gin.Context) (uint, error) {
	val := c.Query("occupation_id")
	if val == "" {
		return 0, fmt.Errorf("occupation_id query parameter is required")
	}
	var id uint
	if _, err := fmt.Sscanf(val, "%d", &id); err != nil {
		return 0, fmt.Errorf("invalid occupation_id value")
	}
	return id, nil
}

func (ctrl *JobController) GetYearlyJobTrendByOccupationHandler(c *gin.Context) {
	occupationID, err := parseOccupationID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := ctrl.repo.GetYearlyJobTrendByOccupation(occupationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve yearly job trend by occupation",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"occupation_id": occupationID,
		"count":         len(results),
		"yearly_trend":  results,
	})
}

func (ctrl *JobController) GetTop3JobRolesByOccupationHandler(c *gin.Context) {
	occupationID, err := parseOccupationID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := ctrl.repo.GetTop3JobRolesByOccupation(occupationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve top 3 job roles by occupation",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"occupation_id": occupationID,
		"top_job_roles": results,
	})
}

// parseIndustryAndYear is a shared helper for industry + year scoped handlers
func parseIndustryAndYear(c *gin.Context) (uint, int, error) {
	industryID, err := parseIndustryID(c)
	if err != nil {
		return 0, 0, err
	}

	yearStr := c.Query("year")
	if yearStr == "" {
		return 0, 0, fmt.Errorf("year query parameter is required")
	}
	var year int
	if _, err := fmt.Sscanf(yearStr, "%d", &year); err != nil {
		return 0, 0, fmt.Errorf("invalid year value")
	}

	return industryID, year, nil
}

func (ctrl *JobController) GetYearlyJobTrendByIndustryHandler(c *gin.Context) {
	industryID, err := parseIndustryID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := ctrl.repo.GetYearlyJobTrendByIndustry(industryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve yearly job trend by industry",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"industry_id":  industryID,
		"count":        len(results),
		"yearly_trend": results,
	})
}

func (ctrl *JobController) GetJobCountByExperienceForIndustryAndYearHandler(c *gin.Context) {
	industryID, year, err := parseIndustryAndYear(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := ctrl.repo.GetJobCountByExperienceForIndustryAndYear(industryID, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve job counts by experience",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"industry_id": industryID,
		"year":        year,
		"count":       len(results),
		"experiences": results,
	})
}

func (ctrl *JobController) GetProvinceWiseJobCountForIndustryAndYearHandler(c *gin.Context) {
	industryID, year, err := parseIndustryAndYear(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := ctrl.repo.GetProvinceWiseJobCountForIndustryAndYear(industryID, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve province wise job counts",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"industry_id": industryID,
		"year":        year,
		"count":       len(results),
		"provinces":   results,
	})
}

func (ctrl *JobController) GetJobCountByEducationLevelForIndustryAndYearHandler(c *gin.Context) {
	industryID, year, err := parseIndustryAndYear(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := ctrl.repo.GetJobCountByEducationLevelForIndustryAndYear(industryID, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve job counts by education level",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"industry_id":      industryID,
		"year":             year,
		"count":            len(results),
		"education_levels": results,
	})
}

func (ctrl *JobController) GetTopHiringEmployersForIndustryAndYearHandler(c *gin.Context) {
	industryID, year, err := parseIndustryAndYear(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := ctrl.repo.GetTopHiringEmployersForIndustryAndYear(industryID, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve top hiring employers",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"industry_id": industryID,
		"year":        year,
		"count":       len(results),
		"employers":   results,
	})
}