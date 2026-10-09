package repositories

import "marketlens-go-backend/models"

// Vocational Education CRUD
func (r *JobRepository) CreateVocationalEducation(item *models.VocationalEducation) error {
	return r.db.Create(item).Error
}

func (r *JobRepository) GetAllVocationalEducations() ([]models.VocationalEducation, error) {
	var items []models.VocationalEducation
	err := r.db.Find(&items).Error
	return items, err
}

func (r *JobRepository) GetVocationalEducationByID(id uint) (models.VocationalEducation, error) {
	var item models.VocationalEducation
	err := r.db.First(&item, id).Error
	return item, err
}

func (r *JobRepository) UpdateVocationalEducation(id uint, updates map[string]interface{}) (models.VocationalEducation, error) {
	var item models.VocationalEducation
	if err := r.db.First(&item, id).Error; err != nil {
		return item, err
	}
	if err := r.db.Model(&item).Updates(updates).Error; err != nil {
		return item, err
	}
	return item, nil
}

func (r *JobRepository) DeleteVocationalEducation(id uint) error {
	return r.db.Delete(&models.VocationalEducation{}, id).Error
}
