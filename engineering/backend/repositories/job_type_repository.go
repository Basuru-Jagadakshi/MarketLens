package repositories

import "marketlens-go-backend/models"

// Job Type CRUD
func (r *JobRepository) CreateJobType(item *models.JobType) error {
	return r.db.Create(item).Error
}

func (r *JobRepository) GetAllJobTypes() ([]models.JobType, error) {
	var jobTypes []models.JobType
	if err := r.db.Find(&jobTypes).Error; err != nil {
		return nil, err
	}
	return jobTypes, nil
}

func (r *JobRepository) GetJobTypeByID(id uint) (models.JobType, error) {
	var item models.JobType
	err := r.db.First(&item, id).Error
	return item, err
}

func (r *JobRepository) UpdateJobType(id uint, updates map[string]interface{}) (models.JobType, error) {
	var item models.JobType
	if err := r.db.First(&item, id).Error; err != nil {
		return item, err
	}
	if err := r.db.Model(&item).Updates(updates).Error; err != nil {
		return item, err
	}
	return item, nil
}

func (r *JobRepository) DeleteJobType(id uint) error {
	return r.db.Delete(&models.JobType{}, id).Error
}
