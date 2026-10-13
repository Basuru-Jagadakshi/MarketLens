package service

import (
	"marketlens-go-backend/internal/models"
	"marketlens-go-backend/internal/repositories"
)

func NewExperienceService(store *repositories.Store) *CRUDService[models.Experience] {
	return NewCRUDService(
		store.CreateExperience,
		store.GetAllExperiences,
		store.GetExperienceByID,
		store.UpdateExperience,
		store.DeleteExperience,
		requireNonEmptyField("name", func(e *models.Experience) string { return e.Name }),
	)
}
