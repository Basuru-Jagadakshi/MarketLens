package service

import (
	"context"
	"marketlens-go-backend/internal/models"
	"marketlens-go-backend/internal/repositories"
	"time"
)

func NewEducationLevelService(store *repositories.Store) *CRUDService[models.EducationLevel] {
	return NewCRUDService(
		store.CreateEducationLevel,
		func(ctx context.Context, _, _ int, _, _ *time.Time) ([]models.EducationLevel, error) {
			return store.GetAllEducationLevels(ctx)
		},
		store.GetEducationLevelByID,
		store.UpdateEducationLevel,
		store.DeleteEducationLevel,
		requireNonEmptyField("level", func(e *models.EducationLevel) string { return e.Level }),
	)
}
