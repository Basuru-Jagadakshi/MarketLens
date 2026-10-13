package service

import (
	"marketlens-go-backend/internal/models"
	"marketlens-go-backend/internal/repositories"
)

func NewEducationLevelService(store *repositories.Store) *CRUDService[models.EducationLevel] {
	return NewCRUDService(
		store.CreateEducationLevel,
		store.GetAllEducationLevels,
		store.GetEducationLevelByID,
		store.UpdateEducationLevel,
		store.DeleteEducationLevel,
		requireNonEmptyField("level", func(e *models.EducationLevel) string { return e.Level }),
	)
}
