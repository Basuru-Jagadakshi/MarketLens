package service

import (
	"marketlens-go-backend/internal/models"
	"marketlens-go-backend/internal/repositories"
)

func NewJobTypeService(store *repositories.Store) *CRUDService[models.JobType] {
	return NewCRUDService(
		store.CreateJobType,
		store.GetAllJobTypes,
		store.GetJobTypeByID,
		store.UpdateJobType,
		store.DeleteJobType,
		requireNonEmptyField("type", func(j *models.JobType) string { return j.Type }),
	)
}
