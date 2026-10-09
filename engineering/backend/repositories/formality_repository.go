package repositories

import "marketlens-go-backend/models"

// Formality CRUD
func (r *JobRepository) CreateFormality(item *models.Formality) error {
	return r.db.Create(item).Error
}

func (r *JobRepository) GetAllFormalities() ([]models.Formality, error) {
	var items []models.Formality
	err := r.db.Find(&items).Error
	return items, err
}

func (r *JobRepository) GetFormalityByID(id uint) (models.Formality, error) {
	var item models.Formality
	err := r.db.First(&item, id).Error
	return item, err
}

func (r *JobRepository) UpdateFormality(id uint, updates map[string]interface{}) (models.Formality, error) {
	var item models.Formality
	if err := r.db.First(&item, id).Error; err != nil {
		return item, err
	}
	if err := r.db.Model(&item).Updates(updates).Error; err != nil {
		return item, err
	}
	return item, nil
}

func (r *JobRepository) DeleteFormality(id uint) error {
	return r.db.Delete(&models.Formality{}, id).Error
}
