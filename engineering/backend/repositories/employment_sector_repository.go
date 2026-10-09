package repositories

import "marketlens-go-backend/models"

// Employment Sector CRUD
func (r *JobRepository) CreateEmploymentSector(item *models.EmploymentSector) error {
	return r.db.Create(item).Error
}

func (r *JobRepository) GetAllEmploymentSectors() ([]models.EmploymentSector, error) {
	var items []models.EmploymentSector
	err := r.db.Find(&items).Error
	return items, err
}

func (r *JobRepository) GetEmploymentSectorByID(id uint) (models.EmploymentSector, error) {
	var item models.EmploymentSector
	err := r.db.First(&item, id).Error
	return item, err
}

func (r *JobRepository) UpdateEmploymentSector(id uint, updates map[string]interface{}) (models.EmploymentSector, error) {
	var item models.EmploymentSector
	if err := r.db.First(&item, id).Error; err != nil {
		return item, err
	}
	if err := r.db.Model(&item).Updates(updates).Error; err != nil {
		return item, err
	}
	return item, nil
}

func (r *JobRepository) DeleteEmploymentSector(id uint) error {
	return r.db.Delete(&models.EmploymentSector{}, id).Error
}
