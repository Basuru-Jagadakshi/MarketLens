package service

import (
	"context"
	"marketlens-go-backend/internal/models"
	"marketlens-go-backend/internal/repositories"
	"time"
)

func NewEmploymentSectorService(store *repositories.Store) *CRUDService[models.EmploymentSector] {
	return NewCRUDService(
		store.CreateEmploymentSector,
		func(ctx context.Context, _, _ int, _, _ *time.Time) ([]models.EmploymentSector, error) {
			return store.GetAllEmploymentSectors(ctx)
		},
		store.GetEmploymentSectorByID,
		store.UpdateEmploymentSector,
		store.DeleteEmploymentSector,
		requireNonEmptyField("sector", func(e *models.EmploymentSector) string { return e.Sector }),
	)
}
