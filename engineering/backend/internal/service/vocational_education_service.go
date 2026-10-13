package service

import (
	"marketlens-go-backend/internal/models"
	"marketlens-go-backend/internal/repositories"
)

func NewVocationalEducationService(store *repositories.Store) *CRUDService[models.VocationalEducation] {
	return NewCRUDService(
		store.CreateVocationalEducation,
		store.GetAllVocationalEducations,
		store.GetVocationalEducationByID,
		store.UpdateVocationalEducation,
		store.DeleteVocationalEducation,
		requireNonEmptyField("level", func(v *models.VocationalEducation) string { return v.Level }),
	)
}
