package main

import (
	"marketlens-go-backend/config"
	"marketlens-go-backend/controllers"
	"marketlens-go-backend/repositories"

	"github.com/gin-gonic/gin"
)


func main() { 
	
	config.ConnectDatabase()

	r := gin.Default()

	jobRepo := repositories.NewJobRepository(config.DB)
	jobCtrl := controllers.NewJobController(jobRepo)

	v1 := r.Group("/api/v1")
	{
		v1.DELETE("/jobs/:id", jobCtrl.DeleteJobHandler)
		v1.GET("/jobs", jobCtrl.GetActiveJobsHandler)
		v1.GET("/industries", jobCtrl.GetAllIndustriesHandler)
		v1.GET("/experiences", jobCtrl.GetAllExperiencesHandler)
		v1.GET("/provinces", jobCtrl.GetAllProvincesHandler)
		v1.GET("/job-types", jobCtrl.GetAllJobTypesHandler)
		v1.GET("/industries/skills/count",      jobCtrl.GetUniqueSkillsCountByIndustryHandler)
		v1.GET("/industries/skills/top-demand", jobCtrl.GetMostDemandingSkillByIndustryHandler)
		v1.GET("/industries/skills/top15",      jobCtrl.GetTop15SkillsByIndustryHandler)
		v1.GET("/industries/skills",            jobCtrl.GetAllSkillsByIndustryHandler)
		v1.GET("/industries/employers",         jobCtrl.GetTopHiringEmployersByIndustryHandler)
		v1.GET("/crawler/last-job-count", jobCtrl.GetLastCrawledJobCountHandler)
		v1.GET("/crawler/time-gap",       jobCtrl.GetTimeSinceLastCrawlHandler)
		v1.GET("/crawler/runs",           jobCtrl.GetAllCrawlerRunsHandler)
		v1.GET("/sources",                jobCtrl.GetSourcesWithActiveJobCountHandler)
		v1.GET("/stats/active-jobs",       jobCtrl.GetActiveJobCountWithTrendHandler)
		v1.GET("/stats/by-occupation",     jobCtrl.GetActiveJobCountByOccupationHandler)
		v1.GET("/stats/by-industry",       jobCtrl.GetActiveJobCountByIndustryHandler)
		v1.GET("/stats/by-experience",     jobCtrl.GetActiveJobCountByExperienceHandler)
		v1.GET("/stats/by-education",      jobCtrl.GetActiveJobCountByEducationLevelHandler)
		v1.GET("/stats/remote-vs-onsite",  jobCtrl.GetRemoteVsOnSiteCountHandler)
		v1.GET("/stats/by-job-type",       jobCtrl.GetActiveJobCountByJobTypeHandler)
		v1.GET("/occupations/yearly-trend", jobCtrl.GetYearlyJobTrendByOccupationHandler)
		v1.GET("/occupations/top-job-roles", jobCtrl.GetTop3JobRolesByOccupationHandler)
		v1.GET("/industries/yearly-trend",        jobCtrl.GetYearlyJobTrendByIndustryHandler)
		v1.GET("/industries/by-experience",       jobCtrl.GetJobCountByExperienceForIndustryAndYearHandler)
		v1.GET("/industries/by-province",         jobCtrl.GetProvinceWiseJobCountForIndustryAndYearHandler)
		v1.GET("/industries/by-education",        jobCtrl.GetJobCountByEducationLevelForIndustryAndYearHandler)
		v1.GET("/industries/top-employers",       jobCtrl.GetTopHiringEmployersForIndustryAndYearHandler)

		crawler := v1.Group("/crawler")
		{
			crawler.POST("/runs", jobCtrl.StartCrawlerRunHandler)
			crawler.POST("/runs/:id/complete", jobCtrl.CompleteCrawlerRunHandler)

			crawler.POST("/radar/lookup", jobCtrl.GetJobsByBucketKeysHandler)

			crawler.POST("/jobs/batch-save", jobCtrl.BatchSaveJobsHandler)
			crawler.POST("/jobs/batch-update", jobCtrl.BatchUpdateDuplicatesHandler)
			crawler.POST("/jobs/reconcile", jobCtrl.ReconcileStaleVacanciesHandler)
		}
	}

	r.Run(":8080")
}
