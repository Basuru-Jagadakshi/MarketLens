package service

import (
	"context"
	"marketlens-go-backend/internal/models"
	"marketlens-go-backend/internal/repositories"
	"time"
)

func NewExperienceService(store *repositories.Store) *CRUDService[models.Experience] {
	return NewCRUDService(
		store.CreateExperience,
		func(ctx context.Context, _, _ int, _, _ *time.Time) ([]models.Experience, error) {
			return store.GetAllExperiences(ctx)
		},
		store.GetExperienceByID,
		store.UpdateExperience,
		store.DeleteExperience,
		requireNonEmptyField("name", func(e *models.Experience) string { return e.Name }),
	)
}
