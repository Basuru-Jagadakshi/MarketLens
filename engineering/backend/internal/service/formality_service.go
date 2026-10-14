package service

import (
	"context"
	"marketlens-go-backend/internal/models"
	"marketlens-go-backend/internal/repositories"
	"time"
)

func NewFormalityService(store *repositories.Store) *CRUDService[models.Formality] {
	return NewCRUDService(
		store.CreateFormality,
		func(ctx context.Context, _, _ int, _, _ *time.Time) ([]models.Formality, error) {
			return store.GetAllFormalities(ctx)
		},
		store.GetFormalityByID,
		store.UpdateFormality,
		store.DeleteFormality,
		requireNonEmptyField("formality_type", func(f *models.Formality) string { return f.FormalityType }),
	)
}
