package repositories

import (
	"context"
	"marketlens-go-backend/internal/models"
)

// Job Type CRUD
func (r *Store) CreateJobType(ctx context.Context, item *models.JobType) error {
	return newCRUDRepository[models.JobType](r.db).create(ctx, item)
}

func (r *Store) GetAllJobTypes(ctx context.Context) ([]models.JobType, error) {
	return newCRUDRepository[models.JobType](r.db).getAll(ctx)
}

func (r *Store) GetJobTypeByID(ctx context.Context, id uint) (models.JobType, error) {
	return newCRUDRepository[models.JobType](r.db).getByID(ctx, id)
}

func (r *Store) UpdateJobType(ctx context.Context, id uint, updates map[string]interface{}) (models.JobType, error) {
	return newCRUDRepository[models.JobType](r.db).update(ctx, id, updates)
}

func (r *Store) DeleteJobType(ctx context.Context, id uint) error {
	return newCRUDRepository[models.JobType](r.db).delete(ctx, id)
}
