package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"marketlens-go-backend/internal/api"
	"marketlens-go-backend/internal/config"
	"marketlens-go-backend/internal/ingestion"
	"marketlens-go-backend/internal/repositories"
	"marketlens-go-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// ginMode decides gin's run mode from APP_ENV (falling back to GIN_MODE,
// then defaulting to release). gin.Default()'s debug mode logs every
// registered route and request — fine for local development, but a minor
// information leak in production, so it's opt-in rather than the default.
func ginMode() string {
	if m := os.Getenv("GIN_MODE"); m == gin.DebugMode || m == gin.TestMode {
		return m
	}
	return gin.ReleaseMode
}

func main() {
	config.ConnectDatabase()
	gin.SetMode(ginMode())

	store := repositories.NewStore(config.DB)

	llmClient := ingestion.NewDeepSeekClient(http.DefaultClient)
	industryClassifier := ingestion.NewIndustryClassifier(
		service.NewIndustrySectorService(store),
		service.NewIndustryDivisionService(store),
		service.NewIndustryGroupService(store),
		service.NewIndustryClassService(store),
		service.NewIndustrySubclassService(store),
		llmClient,
	)
	occupationClassifier := ingestion.NewOccupationClassifier(
		service.NewMajorGroupService(store),
		service.NewSubMajorGroupService(store),
		service.NewMinorGroupService(store),
		service.NewUnitGroupService(store),
		service.NewOccupationGroupService(store),
		llmClient,
	)
	metadataBuilder := ingestion.NewMetadataBuilder(
		service.NewFormalityService(store),
		service.NewGenderService(store),
		service.NewVocationalEducationService(store),
		service.NewEmploymentSectorService(store),
		service.NewEducationLevelService(store),
		service.NewExperienceService(store),
		30*time.Minute,
	)
	ingestionSvc := ingestion.NewIngestionService(
		service.NewJobPostService(store),
		llmClient,
		industryClassifier,
		occupationClassifier,
		metadataBuilder,
		0.65,
	)

	h := api.NewHandler(store, ingestionSvc)

	r := gin.Default()
	api.RegisterRoutes(r, h)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("API server failed: %v", err)
	}
}
