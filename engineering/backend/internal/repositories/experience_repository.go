package repositories

import (
	"context"
	"marketlens-go-backend/internal/models"
)

// Experience CRUD
func (r *Store) GetAllExperiences(ctx context.Context) ([]models.Experience, error) {
	return newCRUDRepository[models.Experience](r.db).getAll(ctx)
}

func (r *Store) CreateExperience(ctx context.Context, item *models.Experience) error {
	return newCRUDRepository[models.Experience](r.db).create(ctx, item)
}

func (r *Store) GetExperienceByID(ctx context.Context, id uint) (models.Experience, error) {
	return newCRUDRepository[models.Experience](r.db).getByID(ctx, id)
}

func (r *Store) UpdateExperience(ctx context.Context, id uint, updates map[string]interface{}) (models.Experience, error) {
	return newCRUDRepository[models.Experience](r.db).update(ctx, id, updates)
}

func (r *Store) DeleteExperience(ctx context.Context, id uint) error {
	return newCRUDRepository[models.Experience](r.db).delete(ctx, id)
}
