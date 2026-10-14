package repositories

import (
	"context"
	"marketlens-go-backend/internal/models"
)

// Education level CRUD
func (r *Store) CreateEducationLevel(ctx context.Context, item *models.EducationLevel) error {
	return newCRUDRepository[models.EducationLevel](r.db).create(ctx, item)
}

func (r *Store) GetAllEducationLevels(ctx context.Context) ([]models.EducationLevel, error) {
	return newCRUDRepository[models.EducationLevel](r.db).getAll(ctx, 0, 0, nil, nil)
}

func (r *Store) GetEducationLevelByID(ctx context.Context, id uint) (models.EducationLevel, error) {
	return newCRUDRepository[models.EducationLevel](r.db).getByID(ctx, id)
}

func (r *Store) UpdateEducationLevel(ctx context.Context, id uint, updates map[string]interface{}) (models.EducationLevel, error) {
	return newCRUDRepository[models.EducationLevel](r.db).update(ctx, id, updates)
}

func (r *Store) DeleteEducationLevel(ctx context.Context, id uint) error {
	return newCRUDRepository[models.EducationLevel](r.db).delete(ctx, id)
}
