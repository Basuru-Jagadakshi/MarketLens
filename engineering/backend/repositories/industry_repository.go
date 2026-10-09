package repositories

import (
	"marketlens-go-backend/models"
	"time"
)

// Industry levels by parent ids
func (r *JobRepository) GetAllIndustrySectorsForDateRange(fromDate, toDate time.Time) ([]models.IndustrySector, error) {
	var items []models.IndustrySector
	endOfToDate := toDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	err := r.db.
		Unscoped().
		Where("created_at <= ?", endOfToDate).
		Where("deleted_at IS NULL OR deleted_at::date >= ?", fromDate.Format("2006-01-02")).
		Find(&items).Error
	return items, err
}

func (r *JobRepository) GetIndustryDivisionsByIndustrySector(industrySectorID uint) ([]models.IndustryDivision, error) {
	var items []models.IndustryDivision
	err := r.db.Where("industry_sector_id = ?", industrySectorID).Find(&items).Error
	return items, err
}

func (r *JobRepository) GetIndustryDivisionsByIndustrySectorForDateRange(industrySectorID uint, fromDate, toDate time.Time) ([]models.IndustryDivision, error) {
	var items []models.IndustryDivision
	endOfToDate := toDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	err := r.db.
		Unscoped().
		Where("industry_sector_id = ?", industrySectorID).
		Where("created_at <= ?", endOfToDate).
		Where("deleted_at IS NULL OR deleted_at::date >= ?", fromDate.Format("2006-01-02")).
		Find(&items).Error
	return items, err
}

func (r *JobRepository) GetIndustryGroupsByIndustryDivision(industryDivisionID uint) ([]models.IndustryGroup, error) {
	var items []models.IndustryGroup
	err := r.db.Where("industry_division_id = ?", industryDivisionID).Find(&items).Error
	return items, err
}

func (r *JobRepository) GetIndustryGroupsByIndustryDivisionForDateRange(industryDivisionID uint, fromDate, toDate time.Time) ([]models.IndustryGroup, error) {
	var items []models.IndustryGroup
	endOfToDate := toDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	err := r.db.
		Unscoped().
		Where("industry_division_id = ?", industryDivisionID).
		Where("created_at <= ?", endOfToDate).
		Where("deleted_at IS NULL OR deleted_at::date >= ?", fromDate.Format("2006-01-02")).
		Find(&items).Error
	return items, err
}

func (r *JobRepository) GetIndustryClassesByIndustryGroup(industryGroupID uint) ([]models.IndustryClass, error) {
	var items []models.IndustryClass
	err := r.db.Where("industry_group_id = ?", industryGroupID).Find(&items).Error
	return items, err
}

func (r *JobRepository) GetIndustryClassesByIndustryGroupForDateRange(industryGroupID uint, fromDate, toDate time.Time) ([]models.IndustryClass, error) {
	var items []models.IndustryClass
	endOfToDate := toDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	err := r.db.
		Unscoped().
		Where("industry_group_id = ?", industryGroupID).
		Where("created_at <= ?", endOfToDate).
		Where("deleted_at IS NULL OR deleted_at::date >= ?", fromDate.Format("2006-01-02")).
		Find(&items).Error
	return items, err
}

func (r *JobRepository) GetIndustrySubclassesByIndustryClass(industryClassID uint) ([]models.IndustrySubclass, error) {
	var items []models.IndustrySubclass
	err := r.db.Where("industry_class_id = ?", industryClassID).Find(&items).Error
	return items, err
}

func (r *JobRepository) GetIndustrySubclassesByIndustryClassForDateRange(industryClassID uint, fromDate, toDate time.Time) ([]models.IndustrySubclass, error) {
	var items []models.IndustrySubclass
	endOfToDate := toDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	err := r.db.
		Unscoped().
		Where("industry_class_id = ?", industryClassID).
		Where("created_at <= ?", endOfToDate).
		Where("deleted_at IS NULL OR deleted_at::date >= ?", fromDate.Format("2006-01-02")).
		Find(&items).Error
	return items, err
}

// Industry sector CRUD
func (r *JobRepository) CreateIndustrySector(item *models.IndustrySector) error {
	return r.db.Create(item).Error
}

func (r *JobRepository) GetAllIndustrySectors() ([]models.IndustrySector, error) {
	var items []models.IndustrySector
	err := r.db.Find(&items).Error
	return items, err
}

func (r *JobRepository) GetIndustrySectorByID(id uint) (models.IndustrySector, error) {
	var item models.IndustrySector
	err := r.db.First(&item, id).Error
	return item, err
}

func (r *JobRepository) UpdateIndustrySector(id uint, updates map[string]interface{}) (models.IndustrySector, error) {
	var item models.IndustrySector
	if err := r.db.First(&item, id).Error; err != nil {
		return item, err
	}
	if err := r.db.Model(&item).Updates(updates).Error; err != nil {
		return item, err
	}
	return item, nil
}

func (r *JobRepository) DeleteIndustrySector(id uint) error {
	return r.db.Delete(&models.IndustrySector{}, id).Error
}

// Industry division CRUD
func (r *JobRepository) CreateIndustryDivision(item *models.IndustryDivision) error {
	return r.db.Create(item).Error
}

func (r *JobRepository) GetAllIndustryDivisions() ([]models.IndustryDivision, error) {
	var items []models.IndustryDivision
	err := r.db.Preload("IndustrySector").Find(&items).Error
	return items, err
}

func (r *JobRepository) GetIndustryDivisionByID(id uint) (models.IndustryDivision, error) {
	var item models.IndustryDivision
	err := r.db.Preload("IndustrySector").First(&item, id).Error
	return item, err
}

func (r *JobRepository) UpdateIndustryDivision(id uint, updates map[string]interface{}) (models.IndustryDivision, error) {
	var item models.IndustryDivision
	if err := r.db.First(&item, id).Error; err != nil {
		return item, err
	}
	if err := r.db.Model(&item).Updates(updates).Error; err != nil {
		return item, err
	}
	return item, nil
}

func (r *JobRepository) DeleteIndustryDivision(id uint) error {
	return r.db.Delete(&models.IndustryDivision{}, id).Error
}

// Industry group CRUD
func (r *JobRepository) CreateIndustryGroup(item *models.IndustryGroup) error {
	return r.db.Create(item).Error
}

func (r *JobRepository) GetAllIndustryGroups() ([]models.IndustryGroup, error) {
	var items []models.IndustryGroup
	err := r.db.Preload("IndustryDivision").Find(&items).Error
	return items, err
}

func (r *JobRepository) GetIndustryGroupByID(id uint) (models.IndustryGroup, error) {
	var item models.IndustryGroup
	err := r.db.Preload("IndustryDivision").First(&item, id).Error
	return item, err
}

func (r *JobRepository) UpdateIndustryGroup(id uint, updates map[string]interface{}) (models.IndustryGroup, error) {
	var item models.IndustryGroup
	if err := r.db.First(&item, id).Error; err != nil {
		return item, err
	}
	if err := r.db.Model(&item).Updates(updates).Error; err != nil {
		return item, err
	}
	return item, nil
}

func (r *JobRepository) DeleteIndustryGroup(id uint) error {
	return r.db.Delete(&models.IndustryGroup{}, id).Error
}

// Industry class CRUD
func (r *JobRepository) CreateIndustryClass(item *models.IndustryClass) error {
	return r.db.Create(item).Error
}

func (r *JobRepository) GetAllIndustryClasses() ([]models.IndustryClass, error) {
	var items []models.IndustryClass
	err := r.db.Preload("IndustryGroup").Find(&items).Error
	return items, err
}

func (r *JobRepository) GetIndustryClassByID(id uint) (models.IndustryClass, error) {
	var item models.IndustryClass
	err := r.db.Preload("IndustryGroup").First(&item, id).Error
	return item, err
}

func (r *JobRepository) UpdateIndustryClass(id uint, updates map[string]interface{}) (models.IndustryClass, error) {
	var item models.IndustryClass
	if err := r.db.First(&item, id).Error; err != nil {
		return item, err
	}
	if err := r.db.Model(&item).Updates(updates).Error; err != nil {
		return item, err
	}
	return item, nil
}

func (r *JobRepository) DeleteIndustryClass(id uint) error {
	return r.db.Delete(&models.IndustryClass{}, id).Error
}

// Industry sub class CRUD
func (r *JobRepository) CreateIndustrySubclass(item *models.IndustrySubclass) error {
	return r.db.Create(item).Error
}

func (r *JobRepository) GetAllIndustrySubclasses(limit, offset int) ([]models.IndustrySubclass, int64, error) {
	var items []models.IndustrySubclass
	var total int64

	if err := r.db.Model(&models.IndustrySubclass{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.Preload("IndustryClass")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("id ASC").Find(&items).Error
	return items, total, err
}

func (r *JobRepository) GetIndustrySubclassByID(id uint) (models.IndustrySubclass, error) {
	var item models.IndustrySubclass
	err := r.db.Preload("IndustryClass").First(&item, id).Error
	return item, err
}

func (r *JobRepository) UpdateIndustrySubclass(id uint, updates map[string]interface{}) (models.IndustrySubclass, error) {
	var item models.IndustrySubclass
	if err := r.db.First(&item, id).Error; err != nil {
		return item, err
	}
	if err := r.db.Model(&item).Updates(updates).Error; err != nil {
		return item, err
	}
	return item, nil
}

func (r *JobRepository) DeleteIndustrySubclass(id uint) error {
	return r.db.Delete(&models.IndustrySubclass{}, id).Error
}
