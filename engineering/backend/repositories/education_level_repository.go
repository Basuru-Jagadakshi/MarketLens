package repositories

import "marketlens-go-backend/models"

// IMPORTANT:
// there are 16 resources x 5 CRUD routes = 80 near-identical handlers + repository methods.
// a generic CRUD (repository + service + handler) would replace most of this code:
// - 7 flat lookups (education levels, genders, ...) have one string field each
// - 10 hierarchy tables (occupation/industry levels) also need a parent id and a check that the parent exists
// validation currently lives in the HTTP handlers, so MCP tools (which call the repository directly) bypass it.
// move validation and rules into a service layer (context.Context + plain Go types, no gin) that both
// the HTTP handlers and the MCP tools call. finish the backend with that design first, then switch MCP over.

// Education level CRUD
func (r *JobRepository) CreateEducationLevel(item *models.EducationLevel) error {
	return r.db.Create(item).Error
}

// repository methods don't accept a context.Context, so queries keep running after the client
// disconnects or the request times out. add ctx as the first parameter (all repository methods)
// and run queries with r.db.WithContext(ctx). handlers should pass c.Request.Context(), not c itself.
// ex: func (r *JobRepository) GetAllEducationLevels(ctx context.Context) ([]models.EducationLevel, error)
func (r *JobRepository) GetAllEducationLevels() ([]models.EducationLevel, error) {
	var items []models.EducationLevel
	// currently the gorm query doesnt care about the order of the rows.
	// but this leads the frontend to show the items in different order. (ex: dropdowns)
	// better if we can order the list items

	// here is a more cleaner way to do this
	// if err := r.db.Find(&items).Error; err != nil {
	//return nil, fmt.Errorf("get all education levels: %w", err)
	// }
	// return items, nil
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

// here for the missing or already deleted IDs the route returns 200 "deleted successfully"
// this is a bug
// better we use rows affected here.
func (r *JobRepository) DeleteEducationLevel(id uint) error {
	return r.db.Delete(&models.EducationLevel{}, id).Error
}
