package service

import (
	"context"
	"marketlens-go-backend/internal/models"
	"marketlens-go-backend/internal/repositories"
	"time"
)

func NewVocationalEducationService(store *repositories.Store) *CRUDService[models.VocationalEducation] {
	return NewCRUDService(
		store.CreateVocationalEducation,
		func(ctx context.Context, _, _ int, _, _ *time.Time) ([]models.VocationalEducation, error) {
			return store.GetAllVocationalEducations(ctx)
		},
		store.GetVocationalEducationByID,
		store.UpdateVocationalEducation,
		store.DeleteVocationalEducation,
		requireNonEmptyField("level", func(v *models.VocationalEducation) string { return v.Level }),
	)
}
