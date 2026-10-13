package repositories

import (
	"context"
	"marketlens-go-backend/internal/models"
)

// Formality CRUD
func (r *Store) CreateFormality(ctx context.Context, item *models.Formality) error {
	return newCRUDRepository[models.Formality](r.db).create(ctx, item)
}

func (r *Store) GetAllFormalities(ctx context.Context) ([]models.Formality, error) {
	return newCRUDRepository[models.Formality](r.db).getAll(ctx)
}

func (r *Store) GetFormalityByID(ctx context.Context, id uint) (models.Formality, error) {
	return newCRUDRepository[models.Formality](r.db).getByID(ctx, id)
}

func (r *Store) UpdateFormality(ctx context.Context, id uint, updates map[string]interface{}) (models.Formality, error) {
	return newCRUDRepository[models.Formality](r.db).update(ctx, id, updates)
}

func (r *Store) DeleteFormality(ctx context.Context, id uint) error {
	return newCRUDRepository[models.Formality](r.db).delete(ctx, id)
}
