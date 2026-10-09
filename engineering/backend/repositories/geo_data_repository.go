package repositories

import "marketlens-go-backend/models"

// CRUD for database entities
// Geo data CRUD
func (r *JobRepository) CreateGeoData(item *models.GeoData) error {
	return r.db.Create(item).Error
}

func (r *JobRepository) GetAllGeoData() ([]models.GeoData, error) {
	var items []models.GeoData
	err := r.db.Find(&items).Error
	return items, err
}

func (r *JobRepository) GetGeoDataByID(id uint) (models.GeoData, error) {
	var item models.GeoData
	err := r.db.First(&item, id).Error
	return item, err
}

func (r *JobRepository) UpdateGeoData(id uint, updates map[string]interface{}) (models.GeoData, error) {
	var item models.GeoData
	if err := r.db.First(&item, id).Error; err != nil {
		return item, err
	}
	if err := r.db.Model(&item).Updates(updates).Error; err != nil {
		return item, err
	}
	return item, nil
}

func (r *JobRepository) DeleteGeoData(id uint) error {
	return r.db.Delete(&models.GeoData{}, id).Error
}
