package service

import (
	"marketlens-go-backend/internal/models"
	"marketlens-go-backend/internal/repositories"
)

func NewGenderService(store *repositories.Store) *CRUDService[models.Gender] {
	return NewCRUDService(
		store.CreateGender,
		store.GetAllGenders,
		store.GetGenderByID,
		store.UpdateGender,
		store.DeleteGender,
		requireNonEmptyField("gender_type", func(g *models.Gender) string { return g.GenderType }),
	)
}
