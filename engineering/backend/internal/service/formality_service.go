package service

import (
	"marketlens-go-backend/internal/models"
	"marketlens-go-backend/internal/repositories"
)

func NewFormalityService(store *repositories.Store) *CRUDService[models.Formality] {
	return NewCRUDService(
		store.CreateFormality,
		store.GetAllFormalities,
		store.GetFormalityByID,
		store.UpdateFormality,
		store.DeleteFormality,
		requireNonEmptyField("formality_type", func(f *models.Formality) string { return f.FormalityType }),
	)
}
