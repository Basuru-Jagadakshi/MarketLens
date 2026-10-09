package repositories

import "marketlens-go-backend/models"

// Education level CRUD
func (r *JobRepository) CreateEducationLevel(item *models.EducationLevel) error {
	return r.db.Create(item).Error
}

func (r *JobRepository) GetAllEducationLevels() ([]models.EducationLevel, error) {
	var items []models.EducationLevel
	err := r.db.Find(&items).Error
	return items, err
}

func (r *JobRepository) GetEducationLevelByID(id uint) (models.EducationLevel, error) {
	var item models.EducationLevel
	err := r.db.First(&item, id).Error
	return item, err
}

func (r *JobRepository) UpdateEducationLevel(id uint, updates map[string]interface{}) (models.EducationLevel, error) {
	var item models.EducationLevel
	if err := r.db.First(&item, id).Error; err != nil {
		return item, err
	}
	if err := r.db.Model(&item).Updates(updates).Error; err != nil {
		return item, err
	}
	return item, nil
}

func (r *JobRepository) DeleteEducationLevel(id uint) error {
	return r.db.Delete(&models.EducationLevel{}, id).Error
}
