package api

import (
	"marketlens-go-backend/internal/auth"

	"github.com/gin-gonic/gin"
)

// registerCRUD wires the standard five routes for a flat-lookup or
// hierarchy resource: reads are public, writes need auth + a
// "<scope>:<action>" scope.
func registerCRUD(rg *gin.RouterGroup, path, scope string, create, list, get, update, del gin.HandlerFunc) {
	g := rg.Group(path)
	g.GET("", list)
	g.GET("/:id", get)
	g.POST("", auth.AuthRequired(), auth.RequireScope(scope+":create"), create)
	g.PUT("/:id", auth.AuthRequired(), auth.RequireScope(scope+":update"), update)
	g.DELETE("/:id", auth.AuthRequired(), auth.RequireScope(scope+":delete"), del)
}

// RegisterRoutes attaches every API route and its middleware to r.
func RegisterRoutes(r *gin.Engine, h *Handler) {
	// Kubernetes probes: unversioned, outside /api/v1
	r.GET("/healthz", h.HealthzHandler)
	r.GET("/readyz", h.ReadyzHandler)

	v1 := r.Group("/api/v1")
	{
		// Dashboard endpoints
		v1.GET("/vacancy-trend", h.GetVacancyTrendHandler)
		v1.GET("/vacancy-total", h.GetTotalVacancyCountHandler)
		v1.GET("/occupations/by-date-range", h.GetOccupationJobCountByDateRangeHandler)
		v1.GET("/industries/by-date-range", h.GetIndustryJobCountByDateRangeHandler)
		v1.GET("/:standard/:level/:id/total-job-count", h.GetTotalJobCountByLevelHandler)
		v1.GET("/:standard/:level/:id/children", h.GetLevelChildrenHandler)
		v1.GET("/:standard/:level/:id/employment-sector", h.GetEmploymentSectorByLevelHandler)
		v1.GET("/:standard/:level/:id/experience", h.GetExperienceByLevelHandler)
		v1.GET("/:standard/:level/:id/province", h.GetProvinceByLevelHandler)
		v1.GET("/:standard/:level/:id/education", h.GetEducationLevelByLevelHandler)
		v1.GET("/:standard/:level/:id/formality", h.GetFormalityByLevelHandler)
		v1.GET("/:standard/:level/:id/gender", h.GetGenderByLevelHandler)
		v1.GET("/:standard/:level/:id/vocational-education", h.GetVocationalEducationByLevelHandler)
		v1.GET("/:standard/:level/:id/remote-onsite-hybrid", h.GetRemoteOnSiteHybridByLevelHandler)
		v1.GET("/:standard/:level/:id/job-type", h.GetJobTypeByLevelHandler)
		v1.GET("/occupation/:level/:id/top-15-skills", h.GetTop15SkillsByOccupationLevelHandler)
		v1.GET("/occupation/:level/:id/all-skills", h.GetAllSkillsByOccupationLevelHandler)
		v1.GET("/occupation/:level/:id/top-hiring-employers", h.GetTopHiringEmployersByOccupationLevelHandler)

		// Crawler/ingestion pipeline (authenticated)
		ingestionRoutes := v1.Group("/crawler")
		ingestionRoutes.Use(auth.AuthRequired())
		{
			ingestionRoutes.POST("/runs", auth.RequireScope("crawler:runs"), h.StartCrawlerRunHandler)
			ingestionRoutes.POST("/runs/:id/complete", auth.RequireScope("crawler:complete"), h.CompleteCrawlerRunHandler)
			ingestionRoutes.POST("/jobs/batch-save", auth.RequireScope("crawler:batch-save"), h.BatchSaveJobsHandler)
			ingestionRoutes.POST("/jobs/reconcile", auth.RequireScope("crawler:reconcile"), h.ReconcileStaleVacanciesHandler)
		}

		// Flat lookup tables
		registerCRUD(v1, "/education-levels", "education-levels",
			h.CreateEducationLevelHandler, h.GetAllEducationLevelsHandler, h.GetEducationLevelByIDHandler,
			h.UpdateEducationLevelHandler, h.DeleteEducationLevelHandler)

		registerCRUD(v1, "/formalities", "formalities",
			h.CreateFormalityHandler, h.GetAllFormalitiesHandler, h.GetFormalityByIDHandler,
			h.UpdateFormalityHandler, h.DeleteFormalityHandler)

		registerCRUD(v1, "/genders", "genders",
			h.CreateGenderHandler, h.GetAllGendersHandler, h.GetGenderByIDHandler,
			h.UpdateGenderHandler, h.DeleteGenderHandler)

		registerCRUD(v1, "/employment-sectors", "employment-sectors",
			h.CreateEmploymentSectorHandler, h.GetAllEmploymentSectorsHandler, h.GetEmploymentSectorByIDHandler,
			h.UpdateEmploymentSectorHandler, h.DeleteEmploymentSectorHandler)

		registerCRUD(v1, "/vocational-educations", "vocational-educations",
			h.CreateVocationalEducationHandler, h.GetAllVocationalEducationsHandler, h.GetVocationalEducationByIDHandler,
			h.UpdateVocationalEducationHandler, h.DeleteVocationalEducationHandler)

		registerCRUD(v1, "/experiences", "experiences",
			h.CreateExperienceHandler, h.GetAllExperiencesHandler, h.GetExperienceByIDHandler,
			h.UpdateExperienceHandler, h.DeleteExperienceHandler)

		// Occupation hierarchy
		registerCRUD(v1, "/major-groups", "major-groups",
			h.CreateMajorGroupHandler, h.GetAllMajorGroupsHandler, h.GetMajorGroupByIDHandler,
			h.UpdateMajorGroupHandler, h.DeleteMajorGroupHandler)
		v1.GET("/major-groups/:id/sub-major-groups", h.GetSubMajorGroupsByMajorGroupHandler)

		registerCRUD(v1, "/sub-major-groups", "sub-major-groups",
			h.CreateSubMajorGroupHandler, h.GetAllSubMajorGroupsHandler, h.GetSubMajorGroupByIDHandler,
			h.UpdateSubMajorGroupHandler, h.DeleteSubMajorGroupHandler)
		v1.GET("/sub-major-groups/:id/minor-groups", h.GetMinorGroupsBySubMajorGroupHandler)

		registerCRUD(v1, "/minor-groups", "minor-groups",
			h.CreateMinorGroupHandler, h.GetAllMinorGroupsHandler, h.GetMinorGroupByIDHandler,
			h.UpdateMinorGroupHandler, h.DeleteMinorGroupHandler)
		v1.GET("/minor-groups/:id/unit-groups", h.GetUnitGroupsByMinorGroupHandler)

		registerCRUD(v1, "/unit-groups", "unit-groups",
			h.CreateUnitGroupHandler, h.GetAllUnitGroupsHandler, h.GetUnitGroupByIDHandler,
			h.UpdateUnitGroupHandler, h.DeleteUnitGroupHandler)
		v1.GET("/unit-groups/:id/occupation-groups", h.GetOccupationGroupsByUnitGroupHandler)

		registerCRUD(v1, "/occupation-groups", "occupation-groups",
			h.CreateOccupationGroupHandler, h.GetAllOccupationGroupsHandler, h.GetOccupationGroupByIDHandler,
			h.UpdateOccupationGroupHandler, h.DeleteOccupationGroupHandler)

		// Industry hierarchy
		registerCRUD(v1, "/industry-sectors", "industry-sectors",
			h.CreateIndustrySectorHandler, h.GetAllIndustrySectorsHandler, h.GetIndustrySectorByIDHandler,
			h.UpdateIndustrySectorHandler, h.DeleteIndustrySectorHandler)
		v1.GET("/industry-sectors/:id/industry-divisions", h.GetIndustryDivisionsByIndustrySectorHandler)

		registerCRUD(v1, "/industry-divisions", "industry-divisions",
			h.CreateIndustryDivisionHandler, h.GetAllIndustryDivisionsHandler, h.GetIndustryDivisionByIDHandler,
			h.UpdateIndustryDivisionHandler, h.DeleteIndustryDivisionHandler)
		v1.GET("/industry-divisions/:id/industry-groups", h.GetIndustryGroupsByIndustryDivisionHandler)

		registerCRUD(v1, "/industry-groups", "industry-groups",
			h.CreateIndustryGroupHandler, h.GetAllIndustryGroupsHandler, h.GetIndustryGroupByIDHandler,
			h.UpdateIndustryGroupHandler, h.DeleteIndustryGroupHandler)
		v1.GET("/industry-groups/:id/industry-classes", h.GetIndustryClassesByIndustryGroupHandler)

		registerCRUD(v1, "/industry-classes", "industry-classes",
			h.CreateIndustryClassHandler, h.GetAllIndustryClassesHandler, h.GetIndustryClassByIDHandler,
			h.UpdateIndustryClassHandler, h.DeleteIndustryClassHandler)
		v1.GET("/industry-classes/:id/industry-subclasses", h.GetIndustrySubclassesByIndustryClassHandler)

		// Note: scope stays "industry-sub-classes" (hyphenated) even though
		// the path is "/industry-subclasses" — matches the scope names
		// already granted in the identity provider.
		registerCRUD(v1, "/industry-subclasses", "industry-sub-classes",
			h.CreateIndustrySubclassHandler, h.GetAllIndustrySubclassesHandler, h.GetIndustrySubclassByIDHandler,
			h.UpdateIndustrySubclassHandler, h.DeleteIndustrySubclassHandler)
	}
}
