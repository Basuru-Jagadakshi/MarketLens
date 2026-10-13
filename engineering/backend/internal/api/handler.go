package api

import (
	"marketlens-go-backend/internal/ingestion"
	"marketlens-go-backend/internal/models"
	"marketlens-go-backend/internal/repositories"
	"marketlens-go-backend/internal/service"
)

type Handler struct {
	repo      *repositories.Store // health checks only — business logic goes through the services below
	ingestion *ingestion.IngestionService

	analytics   *service.AnalyticsService
	jobPosts    *service.JobPostService
	crawlerRuns *service.CrawlerRunService

	educationLevels      *service.CRUDService[models.EducationLevel]
	formalities          *service.CRUDService[models.Formality]
	genders              *service.CRUDService[models.Gender]
	employmentSectors    *service.CRUDService[models.EmploymentSector]
	vocationalEducations *service.CRUDService[models.VocationalEducation]
	experiences          *service.CRUDService[models.Experience]

	majorGroups      *service.CRUDService[models.MajorGroup]
	subMajorGroups   *service.CRUDService[models.SubMajorGroup]
	minorGroups      *service.CRUDService[models.MinorGroup]
	unitGroups       *service.CRUDService[models.UnitGroup]
	occupationGroups *service.CRUDService[models.OccupationGroup]

	industrySectors    *service.CRUDService[models.IndustrySector]
	industryDivisions  *service.CRUDService[models.IndustryDivision]
	industryGroups     *service.CRUDService[models.IndustryGroup]
	industryClasses    *service.CRUDService[models.IndustryClass]
	industrySubclasses *service.CRUDService[models.IndustrySubclass]
}

func NewHandler(repo *repositories.Store, ingestionSvc *ingestion.IngestionService) *Handler {
	return &Handler{
		repo:      repo,
		ingestion: ingestionSvc,

		analytics:   service.NewAnalyticsService(repo),
		jobPosts:    service.NewJobPostService(repo),
		crawlerRuns: service.NewCrawlerRunService(repo),

		educationLevels:      service.NewEducationLevelService(repo),
		formalities:          service.NewFormalityService(repo),
		genders:              service.NewGenderService(repo),
		employmentSectors:    service.NewEmploymentSectorService(repo),
		vocationalEducations: service.NewVocationalEducationService(repo),
		experiences:          service.NewExperienceService(repo),

		majorGroups:      service.NewMajorGroupService(repo),
		subMajorGroups:   service.NewSubMajorGroupService(repo),
		minorGroups:      service.NewMinorGroupService(repo),
		unitGroups:       service.NewUnitGroupService(repo),
		occupationGroups: service.NewOccupationGroupService(repo),

		industrySectors:    service.NewIndustrySectorService(repo),
		industryDivisions:  service.NewIndustryDivisionService(repo),
		industryGroups:     service.NewIndustryGroupService(repo),
		industryClasses:    service.NewIndustryClassService(repo),
		industrySubclasses: service.NewIndustrySubclassService(repo),
	}
}
