package repositories

import (
	"context"
	"marketlens-go-backend/internal/models"
)

// Vocational Education CRUD
func (r *Store) CreateVocationalEducation(ctx context.Context, item *models.VocationalEducation) error {
	return newCRUDRepository[models.VocationalEducation](r.db).create(ctx, item)
}

func (r *Store) GetAllVocationalEducations(ctx context.Context) ([]models.VocationalEducation, error) {
	return newCRUDRepository[models.VocationalEducation](r.db).getAll(ctx, 0, 0, nil, nil)
}

func (r *Store) GetVocationalEducationByID(ctx context.Context, id uint) (models.VocationalEducation, error) {
	return newCRUDRepository[models.VocationalEducation](r.db).getByID(ctx, id)
}

func (r *Store) UpdateVocationalEducation(ctx context.Context, id uint, updates map[string]interface{}) (models.VocationalEducation, error) {
	return newCRUDRepository[models.VocationalEducation](r.db).update(ctx, id, updates)
}

func (r *Store) DeleteVocationalEducation(ctx context.Context, id uint) error {
	return newCRUDRepository[models.VocationalEducation](r.db).delete(ctx, id)
}
