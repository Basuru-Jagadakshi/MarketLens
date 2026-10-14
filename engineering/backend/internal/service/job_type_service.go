package service

import (
	"context"
	"marketlens-go-backend/internal/models"
	"marketlens-go-backend/internal/repositories"
	"time"
)

func NewJobTypeService(store *repositories.Store) *CRUDService[models.JobType] {
	return NewCRUDService(
		store.CreateJobType,
		func(ctx context.Context, _, _ int, _, _ *time.Time) ([]models.JobType, error) {
			return store.GetAllJobTypes(ctx)
		},
		store.GetJobTypeByID,
		store.UpdateJobType,
		store.DeleteJobType,
		requireNonEmptyField("type", func(j *models.JobType) string { return j.Type }),
	)
}
