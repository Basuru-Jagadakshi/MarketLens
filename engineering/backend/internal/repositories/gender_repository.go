package repositories

import (
	"context"
	"marketlens-go-backend/internal/models"
)

// Gender CRUD
func (r *Store) CreateGender(ctx context.Context, item *models.Gender) error {
	return newCRUDRepository[models.Gender](r.db).create(ctx, item)
}

func (r *Store) GetAllGenders(ctx context.Context) ([]models.Gender, error) {
	return newCRUDRepository[models.Gender](r.db).getAll(ctx, 0, 0, nil, nil)
}

func (r *Store) GetGenderByID(ctx context.Context, id uint) (models.Gender, error) {
	return newCRUDRepository[models.Gender](r.db).getByID(ctx, id)
}

func (r *Store) UpdateGender(ctx context.Context, id uint, updates map[string]interface{}) (models.Gender, error) {
	return newCRUDRepository[models.Gender](r.db).update(ctx, id, updates)
}

func (r *Store) DeleteGender(ctx context.Context, id uint) error {
	return newCRUDRepository[models.Gender](r.db).delete(ctx, id)
}
