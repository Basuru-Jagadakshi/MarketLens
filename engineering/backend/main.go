package main

import (
	"marketlens-go-backend/auth"
	"marketlens-go-backend/config"
	"marketlens-go-backend/controllers"
	"marketlens-go-backend/crawler"
	mcpserver "marketlens-go-backend/mcp"
	"marketlens-go-backend/repositories"

	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// in this main package there are 2 servers running on. (original backend and the mcpserver)
// recommended way to do this is move these 2 services to /cmd under 2 different subfolders.
// this way we can separate the services and decouple the services internaly.
// We need to change this in the dockerfiles as well.

// backend/
// ├── cmd/
// │   ├── api/
// │   │   └── main.go          # Gin server: connect DB → build repo/services → routes → r.Run(":8080")
// │   └── mcp/
// │       └── main.go          # MCP server: connect DB → build repo → mcpserver.New → StartHTTP(":9090")
// ├── internal/
// │   ├── config/              # database.go (move from ./config)
// │   ├── auth/                # JWT/scope middleware (used by api)
// │   ├── models/              # GORM models + payloads (shared by both)
// │   ├── repositories/        # data access (shared by both)
// │   ├── api/                 # Gin handlers (api only) renamed to api/ because controller name comes from spring frameworks according to MVC.
// │   ├── ingestion/           # was crawler/: dedup, classification, metadata, ProcessBatch
// │   ├── llm/                 # (optional) DeepSeek client
// │   └── mcpserver/           # MCP tool registration + auth
// ├── sql/                     # schema.sql, seed*.sql (not Go, stays at root)
// ├── tests/                   # integration tests (can still import internal/: same module)
// ├── openapi.yaml
// ├── go.mod
// └── Dockerfile.backend

// use 'mcpserver' name instead of mcp because it clash with the official mcp package. thats why you have a separate label in the import 'mcpserver'

func main() {

	config.ConnectDatabase()
	// catch the database and the error here to throw the exceptions

	// gin.Default() is the gin.DebugMode where it logs the internal route patterns.
	// using gin.Default() is a security risk. Better we get an environment variable and set the mode here. (ex: APP_MODE == "production")
	// recommended way for stg & prod is gin.setMode(gin.ReleaseMode)
	r := gin.Default()

	// r.Use(gin.Recovery()) Built in recovery middleware catches any panic that happens while the service hanlding a request, writes a 500 response and keep the server running.
	// if you want to do some other custom handling you can use r.Use(gin.CustomRecovery())
	// ex: if you want to save the logs to dbs.
	// refer: https://gin-gonic.com/en/docs/middleware/custom-recovery/

	jobRepo := repositories.NewJobRepository(config.DB)
	llmClient := crawler.NewDeepSeekClient(http.DefaultClient)
	metaBuilder := crawler.NewMetadataBuilder(jobRepo, 30*time.Minute)
	ingestionSvc := crawler.NewIngestionService(jobRepo, llmClient, metaBuilder, 0.65)
	jobCtrl := controllers.NewJobController(jobRepo, ingestionSvc)

	// move all the routes to routes.go inside the /internal/api/
	// Kubernetes liveness/readiness probes - unversioned, sit outside /api/v1
	r.GET("/healthz", jobCtrl.HealthzHandler)
	r.GET("/readyz", jobCtrl.ReadyzHandler)

	// Consider splitting the MCP server into its own binary (cmd/mcp) alongside cmd/api.
	// Right now they share a process, so a fatal error in the MCP goroutine (log.Fatalf) also stops the API,
	// and the two can't be scaled or deployed independently. Trade-off: a second Deployment and DB pool to manage.
	// If we keep one binary for now, at minimum the MCP goroutine shouldn't call log.Fatalf. It should log the error,
	// or shut the whole server down gracefully on purpose.
	mcpServer := mcpserver.New(jobRepo)
	go func() {
		log.Println("MCP server listening on :9090/mcp")
		if err := mcpserver.StartHTTP(mcpServer, ":9090"); err != nil {
			log.Fatalf("MCP server failed: %v", err)
		}
	}()

	// move all the routes to routes.go inside the /internal/api/
	// create a separate function to createCRUD and reuse inside
	// then register the routes in the main file
	// ex:
	// 		h := api.NewHandler(store, ingestionSvc)
	// 		api.RegisterRoutes(r, h)

	// keep both cmd/api/main.go and cmd/mcp/main.go minimal since they are the entry points

	// here is the sample template for routes.go

	// package api

	// import (
	// 	"marketlens-go-backend/internal/auth"

	// 	"github.com/gin-gonic/gin"
	// )

	// // RegisterRoutes attaches all API routes and their middleware to r.
	// func RegisterRoutes(r *gin.Engine, h *Handler) {
	// 	// Kubernetes probes: unversioned, outside /api/v1
	// 	r.GET("/healthz", h.Healthz)
	// 	r.GET("/readyz", h.Readyz)

	// 	v1 := r.Group("/api/v1")

	// 	// Dashboard (public, read-only)
	// 	v1.GET("/vacancy-trend", h.GetVacancyTrend)
	// 	v1.GET("/vacancy-total", h.GetTotalVacancyCount)
	// 	// ...

	// 	// Crawler (authenticated)
	// 	crawler := v1.Group("/crawler", auth.AuthRequired())
	// 	crawler.POST("/runs", auth.RequireScope("crawler:runs"), h.StartCrawlerRun)
	// 	crawler.POST("/jobs/batch-save", auth.RequireScope("crawler:batch-save"), h.BatchSaveJobs)
	// 	// ...

	// 	// Lookup tables
	// 	registerCRUD(v1, "/genders", "genders", h.CreateGender, h.GetAllGenders, h.GetGenderByID, h.UpdateGender, h.DeleteGender)
	// 	registerCRUD(v1, "/formalities", "formalities", h.CreateFormality, h.GetAllFormalities, h.GetFormalityByID, h.UpdateFormality, h.DeleteFormality)
	// 	// ...
	// }

	// // registerCRUD wires the standard five routes: reads are public, writes need auth + scope.
	// func registerCRUD(rg *gin.RouterGroup, path, scope string, create, list, get, update, del gin.HandlerFunc) {
	// 	g := rg.Group(path)
	// 	g.GET("", list)
	// 	g.GET("/:id", get)
	// 	g.POST("", auth.AuthRequired(), auth.RequireScope(scope+":create"), create)
	// 	g.PUT("/:id", auth.AuthRequired(), auth.RequireScope(scope+":update"), update)
	// 	g.DELETE("/:id", auth.AuthRequired(), auth.RequireScope(scope+":delete"), del)
	// }

	v1 := r.Group("/api/v1")
	{
		// Dashboard endpoints
		v1.GET("/vacancy-trend", jobCtrl.GetVacancyTrendHandler)
		v1.GET("/vacancy-total", jobCtrl.GetTotalVacancyCountHandler)
		v1.GET("/occupations/by-date-range", jobCtrl.GetOccupationJobCountByDateRangeHandler)
		v1.GET("/industries/by-date-range", jobCtrl.GetIndustryJobCountByDateRangeHandler)
		v1.GET("/:standard/:level/:id/total-job-count", jobCtrl.GetTotalJobCountByLevelHandler)
		v1.GET("/:standard/:level/:id/children", jobCtrl.GetLevelChildrenHandler)
		v1.GET("/:standard/:level/:id/employment-sector", jobCtrl.GetEmploymentSectorByLevelHandler)
		v1.GET("/:standard/:level/:id/experience", jobCtrl.GetExperienceByLevelHandler)
		v1.GET("/:standard/:level/:id/province", jobCtrl.GetProvinceByLevelHandler)
		v1.GET("/:standard/:level/:id/education", jobCtrl.GetEducationLevelByLevelHandler)
		v1.GET("/:standard/:level/:id/formality", jobCtrl.GetFormalityByLevelHandler)
		v1.GET("/:standard/:level/:id/gender", jobCtrl.GetGenderByLevelHandler)
		v1.GET("/:standard/:level/:id/vocational-education", jobCtrl.GetVocationalEducationByLevelHandler)
		v1.GET("/:standard/:level/:id/remote-onsite-hybrid", jobCtrl.GetRemoteOnSiteHybridByLevelHandler)
		v1.GET("/:standard/:level/:id/job-type", jobCtrl.GetJobTypeByLevelHandler)
		v1.GET("/occupation/:level/:id/top-15-skills", jobCtrl.GetTop15SkillsByOccupationLevelHandler)
		v1.GET("/occupation/:level/:id/all-skills", jobCtrl.GetAllSkillsByOccupationLevelHandler)
		v1.GET("/occupation/:level/:id/top-hiring-employers", jobCtrl.GetTopHiringEmployersByOccupationLevelHandler)

		// think about this crawler naming. because we have another python crawler. use a better name
		crawler := v1.Group("/crawler")
		crawler.Use(auth.AuthRequired())
		{
			crawler.POST("/runs", auth.RequireScope("crawler:runs"), jobCtrl.StartCrawlerRunHandler)
			// does this always complete the run? dont you need a way to change any other status (ex: failed/ completed/ running etc?)
			crawler.POST("/runs/:id/complete", auth.RequireScope("crawler:complete"), jobCtrl.CompleteCrawlerRunHandler)
			crawler.POST("/jobs/batch-save", auth.RequireScope("crawler:batch-save"), jobCtrl.BatchSaveJobsHandler)
			crawler.POST("/jobs/reconcile", auth.RequireScope("crawler:reconcile"), jobCtrl.ReconcileStaleVacanciesHandler)
		}

		educationLevels := v1.Group("/education-levels")
		{
			educationLevels.POST("", auth.AuthRequired(), auth.RequireScope("education-levels:create"), jobCtrl.CreateEducationLevelHandler)
			educationLevels.GET("", jobCtrl.GetAllEducationLevelsHandler)
			educationLevels.GET("/:id", jobCtrl.GetEducationLevelByIDHandler)
			educationLevels.PUT("/:id", auth.AuthRequired(), auth.RequireScope("education-levels:update"), jobCtrl.UpdateEducationLevelHandler)
			educationLevels.DELETE("/:id", auth.AuthRequired(), auth.RequireScope("education-levels:delete"), jobCtrl.DeleteEducationLevelHandler)
		}

		formalities := v1.Group("/formalities")
		{
			formalities.POST("", auth.AuthRequired(), auth.RequireScope("formalities:create"), jobCtrl.CreateFormalityHandler)
			formalities.GET("", jobCtrl.GetAllFormalitiesHandler)
			formalities.GET("/:id", jobCtrl.GetFormalityByIDHandler)
			formalities.PUT("/:id", auth.AuthRequired(), auth.RequireScope("formalities:update"), jobCtrl.UpdateFormalityHandler)
			formalities.DELETE("/:id", auth.AuthRequired(), auth.RequireScope("formalities:delete"), jobCtrl.DeleteFormalityHandler)
		}

		genders := v1.Group("/genders")
		{
			genders.POST("", auth.AuthRequired(), auth.RequireScope("genders:create"), jobCtrl.CreateGenderHandler)
			genders.GET("", jobCtrl.GetAllGendersHandler)
			genders.GET("/:id", jobCtrl.GetGenderByIDHandler)
			genders.PUT("/:id", auth.AuthRequired(), auth.RequireScope("genders:update"), jobCtrl.UpdateGenderHandler)
			genders.DELETE("/:id", auth.AuthRequired(), auth.RequireScope("genders:delete"), jobCtrl.DeleteGenderHandler)
		}

		employmentSectors := v1.Group("/employment-sectors")
		{
			employmentSectors.POST("", auth.AuthRequired(), auth.RequireScope("employment-sectors:create"), jobCtrl.CreateEmploymentSectorHandler)
			employmentSectors.GET("", jobCtrl.GetAllEmploymentSectorsHandler)
			employmentSectors.GET("/:id", jobCtrl.GetEmploymentSectorByIDHandler)
			employmentSectors.PUT("/:id", auth.AuthRequired(), auth.RequireScope("employment-sectors:update"), jobCtrl.UpdateEmploymentSectorHandler)
			employmentSectors.DELETE("/:id", auth.AuthRequired(), auth.RequireScope("employment-sectors:delete"), jobCtrl.DeleteEmploymentSectorHandler)
		}

		vocationalEducations := v1.Group("/vocational-educations")
		{
			vocationalEducations.POST("", auth.AuthRequired(), auth.RequireScope("vocational-educations:create"), jobCtrl.CreateVocationalEducationHandler)
			vocationalEducations.GET("", jobCtrl.GetAllVocationalEducationsHandler)
			vocationalEducations.GET("/:id", jobCtrl.GetVocationalEducationByIDHandler)
			vocationalEducations.PUT("/:id", auth.AuthRequired(), auth.RequireScope("vocational-educations:update"), jobCtrl.UpdateVocationalEducationHandler)
			vocationalEducations.DELETE("/:id", auth.AuthRequired(), auth.RequireScope("vocational-educations:delete"), jobCtrl.DeleteVocationalEducationHandler)
		}

		experiences := v1.Group("/experiences")
		{
			experiences.GET("", jobCtrl.GetAllExperiencesHandler)
			experiences.POST("", auth.AuthRequired(), auth.RequireScope("experiences:create"), jobCtrl.CreateExperienceHandler)
			experiences.GET("/:id", jobCtrl.GetExperienceByIDHandler)
			experiences.PUT("/:id", auth.AuthRequired(), auth.RequireScope("experiences:update"), jobCtrl.UpdateExperienceHandler)
			experiences.DELETE("/:id", auth.AuthRequired(), auth.RequireScope("experiences:delete"), jobCtrl.DeleteExperienceHandler)
		}

		majorGroups := v1.Group("/major-groups")
		{
			majorGroups.POST("", auth.AuthRequired(), auth.RequireScope("major-groups:create"), jobCtrl.CreateMajorGroupHandler)
			majorGroups.GET("", jobCtrl.GetAllMajorGroupsHandler)
			majorGroups.GET("/:id", jobCtrl.GetMajorGroupByIDHandler)
			majorGroups.PUT("/:id", auth.AuthRequired(), auth.RequireScope("major-groups:update"), jobCtrl.UpdateMajorGroupHandler)
			majorGroups.DELETE("/:id", auth.AuthRequired(), auth.RequireScope("major-groups:delete"), jobCtrl.DeleteMajorGroupHandler)
			majorGroups.GET("/:id/sub-major-groups", jobCtrl.GetSubMajorGroupsByMajorGroupHandler)
		}

		subMajorGroups := v1.Group("/sub-major-groups")
		{
			subMajorGroups.POST("", auth.AuthRequired(), auth.RequireScope("sub-major-groups:create"), jobCtrl.CreateSubMajorGroupHandler)
			subMajorGroups.GET("", jobCtrl.GetAllSubMajorGroupsHandler)
			subMajorGroups.GET("/:id", jobCtrl.GetSubMajorGroupByIDHandler)
			subMajorGroups.PUT("/:id", auth.AuthRequired(), auth.RequireScope("sub-major-groups:update"), jobCtrl.UpdateSubMajorGroupHandler)
			subMajorGroups.DELETE("/:id", auth.AuthRequired(), auth.RequireScope("sub-major-groups:delete"), jobCtrl.DeleteSubMajorGroupHandler)
			subMajorGroups.GET("/:id/minor-groups", jobCtrl.GetMinorGroupsBySubMajorGroupHandler)
		}

		minorGroups := v1.Group("/minor-groups")
		{
			minorGroups.POST("", auth.AuthRequired(), auth.RequireScope("minor-groups:create"), jobCtrl.CreateMinorGroupHandler)
			minorGroups.GET("", jobCtrl.GetAllMinorGroupsHandler)
			minorGroups.GET("/:id", jobCtrl.GetMinorGroupByIDHandler)
			minorGroups.PUT("/:id", auth.AuthRequired(), auth.RequireScope("minor-groups:update"), jobCtrl.UpdateMinorGroupHandler)
			minorGroups.DELETE("/:id", auth.AuthRequired(), auth.RequireScope("minor-groups:delete"), jobCtrl.DeleteMinorGroupHandler)
			minorGroups.GET("/:id/unit-groups", jobCtrl.GetUnitGroupsByMinorGroupHandler)
		}

		unitGroups := v1.Group("/unit-groups")
		{
			unitGroups.POST("", auth.AuthRequired(), auth.RequireScope("unit-groups:create"), jobCtrl.CreateUnitGroupHandler)
			unitGroups.GET("", jobCtrl.GetAllUnitGroupsHandler)
			unitGroups.GET("/:id", jobCtrl.GetUnitGroupByIDHandler)
			unitGroups.PUT("/:id", auth.AuthRequired(), auth.RequireScope("unit-groups:update"), jobCtrl.UpdateUnitGroupHandler)
			unitGroups.DELETE("/:id", auth.AuthRequired(), auth.RequireScope("unit-groups:delete"), jobCtrl.DeleteUnitGroupHandler)
			unitGroups.GET("/:id/occupation-groups", jobCtrl.GetOccupationGroupsByUnitGroupHandler)
		}

		occupationGroups := v1.Group("/occupation-groups")
		{
			occupationGroups.POST("", auth.AuthRequired(), auth.RequireScope("occupation-groups:create"), jobCtrl.CreateOccupationGroupHandler)
			occupationGroups.GET("", jobCtrl.GetAllOccupationGroupsHandler)
			occupationGroups.GET("/:id", jobCtrl.GetOccupationGroupByIDHandler)
			occupationGroups.PUT("/:id", auth.AuthRequired(), auth.RequireScope("occupation-groups:update"), jobCtrl.UpdateOccupationGroupHandler)
			occupationGroups.DELETE("/:id", auth.AuthRequired(), auth.RequireScope("occupation-groups:delete"), jobCtrl.DeleteOccupationGroupHandler)
		}

		industrySectors := v1.Group("/industry-sectors")
		{
			industrySectors.POST("", auth.AuthRequired(), auth.RequireScope("industry-sectors:create"), jobCtrl.CreateIndustrySectorHandler)
			industrySectors.GET("", jobCtrl.GetAllIndustrySectorsHandler)
			industrySectors.GET("/:id", jobCtrl.GetIndustrySectorByIDHandler)
			industrySectors.PUT("/:id", auth.AuthRequired(), auth.RequireScope("industry-sectors:update"), jobCtrl.UpdateIndustrySectorHandler)
			industrySectors.DELETE("/:id", auth.AuthRequired(), auth.RequireScope("industry-sectors:delete"), jobCtrl.DeleteIndustrySectorHandler)
			industrySectors.GET("/:id/industry-divisions", jobCtrl.GetIndustryDivisionsByIndustrySectorHandler)
		}

		industryDivisions := v1.Group("/industry-divisions")
		{
			industryDivisions.POST("", auth.AuthRequired(), auth.RequireScope("industry-divisions:create"), jobCtrl.CreateIndustryDivisionHandler)
			industryDivisions.GET("", jobCtrl.GetAllIndustryDivisionsHandler)
			industryDivisions.GET("/:id", jobCtrl.GetIndustryDivisionByIDHandler)
			industryDivisions.PUT("/:id", auth.AuthRequired(), auth.RequireScope("industry-divisions:update"), jobCtrl.UpdateIndustryDivisionHandler)
			industryDivisions.DELETE("/:id", auth.AuthRequired(), auth.RequireScope("industry-divisions:delete"), jobCtrl.DeleteIndustryDivisionHandler)
			industryDivisions.GET("/:id/industry-groups", jobCtrl.GetIndustryGroupsByIndustryDivisionHandler)
		}

		industryGroups := v1.Group("/industry-groups")
		{
			industryGroups.POST("", auth.AuthRequired(), auth.RequireScope("industry-groups:create"), jobCtrl.CreateIndustryGroupHandler)
			industryGroups.GET("", jobCtrl.GetAllIndustryGroupsHandler)
			industryGroups.GET("/:id", jobCtrl.GetIndustryGroupByIDHandler)
			industryGroups.PUT("/:id", auth.AuthRequired(), auth.RequireScope("industry-groups:update"), jobCtrl.UpdateIndustryGroupHandler)
			industryGroups.DELETE("/:id", auth.AuthRequired(), auth.RequireScope("industry-groups:delete"), jobCtrl.DeleteIndustryGroupHandler)
			industryGroups.GET("/:id/industry-classes", jobCtrl.GetIndustryClassesByIndustryGroupHandler)
		}

		industryClasses := v1.Group("/industry-classes")
		{
			industryClasses.POST("", auth.AuthRequired(), auth.RequireScope("industry-classes:create"), jobCtrl.CreateIndustryClassHandler)
			industryClasses.GET("", jobCtrl.GetAllIndustryClassesHandler)
			industryClasses.GET("/:id", jobCtrl.GetIndustryClassByIDHandler)
			industryClasses.PUT("/:id", auth.AuthRequired(), auth.RequireScope("industry-classes:update"), jobCtrl.UpdateIndustryClassHandler)
			industryClasses.DELETE("/:id", auth.AuthRequired(), auth.RequireScope("industry-classes:delete"), jobCtrl.DeleteIndustryClassHandler)
			industryClasses.GET("/:id/industry-subclasses", jobCtrl.GetIndustrySubclassesByIndustryClassHandler)
		}

		industrySubclasses := v1.Group("/industry-subclasses")
		{
			industrySubclasses.POST("", auth.AuthRequired(), auth.RequireScope("industry-sub-classes:create"), jobCtrl.CreateIndustrySubclassHandler)
			industrySubclasses.GET("", jobCtrl.GetAllIndustrySubclassesHandler)
			industrySubclasses.GET("/:id", jobCtrl.GetIndustrySubclassByIDHandler)
			industrySubclasses.PUT("/:id", auth.AuthRequired(), auth.RequireScope("industry-sub-classes:update"), jobCtrl.UpdateIndustrySubclassHandler)
			industrySubclasses.DELETE("/:id", auth.AuthRequired(), auth.RequireScope("industry-sub-classes:delete"), jobCtrl.DeleteIndustrySubclassHandler)
		}
	}

	r.Run(":8080")
}
