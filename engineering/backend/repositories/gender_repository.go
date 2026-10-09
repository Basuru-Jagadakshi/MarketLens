package repositories

import "marketlens-go-backend/models"

// Gender CRUD
func (r *JobRepository) CreateGender(item *models.Gender) error {
	return r.db.Create(item).Error
}

func (r *JobRepository) GetAllGenders() ([]models.Gender, error) {
	var items []models.Gender
	err := r.db.Find(&items).Error
	return items, err
}

func (r *JobRepository) GetGenderByID(id uint) (models.Gender, error) {
	var item models.Gender
	err := r.db.First(&item, id).Error
	return item, err
}

func (r *JobRepository) UpdateGender(id uint, updates map[string]interface{}) (models.Gender, error) {
	var item models.Gender
	if err := r.db.First(&item, id).Error; err != nil {
		return item, err
	}
	if err := r.db.Model(&item).Updates(updates).Error; err != nil {
		return item, err
	}
	return item, nil
}

func (r *JobRepository) DeleteGender(id uint) error {
	return r.db.Delete(&models.Gender{}, id).Error
}
