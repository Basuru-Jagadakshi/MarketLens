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
		// no need to return the item in the error block, return the zero value instead.
		// return models.EducationLevel{}, err
		// returns gorm.ErrRecordNotFound when the id doesn't exist (handler should map it to 404),
		// but other errors (ex: lost connection, timeout) can also come from here.
		return item, err
	}
	if err := r.db.Model(&item).Updates(updates).Error; err != nil {
		// no need to return the item in the error block because here the item contains the old value
		// return models.EducationLevel{}, fmt.Errorf("update education level %d: %w", id, err)
		// wrap with %w to add context while keeping the original error inspectable (errors.Is/As)
		return item, err
	}

	// alternatively, replace First + Updates with a single query (avoids a race between them):
	// res := r.db.Model(&models.EducationLevel{}).Where("id = ?", id).Updates(updates)
	// if res.Error != nil { return models.EducationLevel{}, res.Error }
	// if res.RowsAffected == 0 { return models.EducationLevel{}, gorm.ErrRecordNotFound } // no row matched (missing or soft-deleted)
	// then reload the row with First if you need to return it

	return item, nil
}

func (r *JobRepository) DeleteEducationLevel(id uint) error {
	return r.db.Delete(&models.EducationLevel{}, id).Error
}
