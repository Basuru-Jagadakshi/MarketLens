package repositories

import (
	"context"
	"marketlens-go-backend/internal/models"
)

// Employment Sector CRUD
func (r *Store) CreateEmploymentSector(ctx context.Context, item *models.EmploymentSector) error {
	return newCRUDRepository[models.EmploymentSector](r.db).create(ctx, item)
}

func (r *Store) GetAllEmploymentSectors(ctx context.Context) ([]models.EmploymentSector, error) {
	return newCRUDRepository[models.EmploymentSector](r.db).getAll(ctx)
}

func (r *Store) GetEmploymentSectorByID(ctx context.Context, id uint) (models.EmploymentSector, error) {
	return newCRUDRepository[models.EmploymentSector](r.db).getByID(ctx, id)
}

func (r *Store) UpdateEmploymentSector(ctx context.Context, id uint, updates map[string]interface{}) (models.EmploymentSector, error) {
	return newCRUDRepository[models.EmploymentSector](r.db).update(ctx, id, updates)
}

func (r *Store) DeleteEmploymentSector(ctx context.Context, id uint) error {
	return newCRUDRepository[models.EmploymentSector](r.db).delete(ctx, id)
}
