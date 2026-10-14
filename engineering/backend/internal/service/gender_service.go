package service

import (
	"context"
	"marketlens-go-backend/internal/models"
	"marketlens-go-backend/internal/repositories"
	"time"
)

func NewGenderService(store *repositories.Store) *CRUDService[models.Gender] {
	return NewCRUDService(
		store.CreateGender,
		func(ctx context.Context, _, _ int, _, _ *time.Time) ([]models.Gender, error) {
			return store.GetAllGenders(ctx)
		},
		store.GetGenderByID,
		store.UpdateGender,
		store.DeleteGender,
		requireNonEmptyField("gender_type", func(g *models.Gender) string { return g.GenderType }),
	)
}
