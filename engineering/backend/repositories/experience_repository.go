package repositories

import "marketlens-go-backend/models"

// Experience CRUD
func (r *JobRepository) GetAllExperiences() ([]models.Experience, error) {
	var experiences []models.Experience
	if err := r.db.Find(&experiences).Error; err != nil {
		return nil, err
	}
	return experiences, nil
}

func (r *JobRepository) CreateExperience(item *models.Experience) error {
	return r.db.Create(item).Error
}

func (r *JobRepository) GetExperienceByID(id uint) (models.Experience, error) {
	var item models.Experience
	err := r.db.First(&item, id).Error
	return item, err
}

func (r *JobRepository) UpdateExperience(id uint, updates map[string]interface{}) (models.Experience, error) {
	var item models.Experience
	if err := r.db.First(&item, id).Error; err != nil {
		return item, err
	}
	if err := r.db.Model(&item).Updates(updates).Error; err != nil {
		return item, err
	}
	return item, nil
}

func (r *JobRepository) DeleteExperience(id uint) error {
	return r.db.Delete(&models.Experience{}, id).Error
}
