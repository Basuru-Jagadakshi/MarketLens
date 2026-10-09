package repositories

import (
	"marketlens-go-backend/models"
	"time"
)

// Methods need to multi level filtering in the Crawler
// Occupation levels by parent ids
func (r *JobRepository) GetAllMajorGroupsForDateRange(fromDate, toDate time.Time) ([]models.MajorGroup, error) {
	var items []models.MajorGroup
	endOfToDate := toDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	err := r.db.
		Unscoped().
		Where("created_at <= ?", endOfToDate).
		Where("deleted_at IS NULL OR deleted_at >= ?", fromDate).
		Find(&items).Error
	return items, err
}

func (r *JobRepository) GetSubMajorGroupsByMajorGroup(majorGroupID uint) ([]models.SubMajorGroup, error) {
	var items []models.SubMajorGroup
	err := r.db.Where("major_group_id = ?", majorGroupID).Find(&items).Error
	return items, err
}

func (r *JobRepository) GetSubMajorGroupsByMajorGroupForDateRange(majorGroupID uint, fromDate, toDate time.Time) ([]models.SubMajorGroup, error) {
	var items []models.SubMajorGroup
	endOfToDate := toDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	err := r.db.
		Unscoped().
		Where("major_group_id = ?", majorGroupID).
		Where("created_at <= ?", endOfToDate).
		Where("deleted_at IS NULL OR deleted_at >= ?", fromDate).
		Find(&items).Error
	return items, err
}

func (r *JobRepository) GetMinorGroupsBySubMajorGroup(subMajorGroupID uint) ([]models.MinorGroup, error) {
	var items []models.MinorGroup
	err := r.db.Where("sub_major_group_id = ?", subMajorGroupID).Find(&items).Error
	return items, err
}

func (r *JobRepository) GetMinorGroupsBySubMajorGroupForDateRange(subMajorGroupID uint, fromDate, toDate time.Time) ([]models.MinorGroup, error) {
	var items []models.MinorGroup
	endOfToDate := toDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	err := r.db.
		Unscoped().
		Where("sub_major_group_id = ?", subMajorGroupID).
		Where("created_at <= ?", endOfToDate).
		Where("deleted_at IS NULL OR deleted_at >= ?", fromDate).
		Find(&items).Error
	return items, err
}

func (r *JobRepository) GetUnitGroupsByMinorGroup(minorGroupID uint) ([]models.UnitGroup, error) {
	var items []models.UnitGroup
	err := r.db.Where("minor_group_id = ?", minorGroupID).Find(&items).Error
	return items, err
}

func (r *JobRepository) GetUnitGroupsByMinorGroupForDateRange(minorGroupID uint, fromDate, toDate time.Time) ([]models.UnitGroup, error) {
	var items []models.UnitGroup
	endOfToDate := toDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	err := r.db.
		Unscoped().
		Where("minor_group_id = ?", minorGroupID).
		Where("created_at <= ?", endOfToDate).
		Where("deleted_at IS NULL OR deleted_at >= ?", fromDate).
		Find(&items).Error
	return items, err
}

func (r *JobRepository) GetOccupationGroupsByUnitGroup(unitGroupID uint) ([]models.OccupationGroup, error) {
	var items []models.OccupationGroup
	err := r.db.Where("unit_group_id = ?", unitGroupID).Find(&items).Error
	return items, err
}

func (r *JobRepository) GetOccupationGroupsByUnitGroupForDateRange(unitGroupID uint, fromDate, toDate time.Time) ([]models.OccupationGroup, error) {
	var items []models.OccupationGroup
	endOfToDate := toDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	err := r.db.
		Unscoped().
		Where("unit_group_id = ?", unitGroupID).
		Where("created_at <= ?", endOfToDate).
		Where("deleted_at IS NULL OR deleted_at >= ?", fromDate).
		Find(&items).Error
	return items, err
}

// Major Group CRUD
func (r *JobRepository) CreateMajorGroup(item *models.MajorGroup) error {
	return r.db.Create(item).Error
}

func (r *JobRepository) GetAllMajorGroups() ([]models.MajorGroup, error) {
	var items []models.MajorGroup
	err := r.db.Find(&items).Error
	return items, err
}

func (r *JobRepository) GetMajorGroupByID(id uint) (models.MajorGroup, error) {
	var item models.MajorGroup
	err := r.db.First(&item, id).Error
	return item, err
}

func (r *JobRepository) UpdateMajorGroup(id uint, updates map[string]interface{}) (models.MajorGroup, error) {
	var item models.MajorGroup
	if err := r.db.First(&item, id).Error; err != nil {
		return item, err
	}
	if err := r.db.Model(&item).Updates(updates).Error; err != nil {
		return item, err
	}
	return item, nil
}

func (r *JobRepository) DeleteMajorGroup(id uint) error {
	return r.db.Delete(&models.MajorGroup{}, id).Error
}

// Sub Major Group CRUD
func (r *JobRepository) CreateSubMajorGroup(item *models.SubMajorGroup) error {
	return r.db.Create(item).Error
}

func (r *JobRepository) GetAllSubMajorGroups() ([]models.SubMajorGroup, error) {
	var items []models.SubMajorGroup
	err := r.db.Preload("MajorGroup").Find(&items).Error
	return items, err
}

func (r *JobRepository) GetSubMajorGroupByID(id uint) (models.SubMajorGroup, error) {
	var item models.SubMajorGroup
	err := r.db.Preload("MajorGroup").First(&item, id).Error
	return item, err
}

func (r *JobRepository) UpdateSubMajorGroup(id uint, updates map[string]interface{}) (models.SubMajorGroup, error) {
	var item models.SubMajorGroup
	if err := r.db.First(&item, id).Error; err != nil {
		return item, err
	}
	if err := r.db.Model(&item).Updates(updates).Error; err != nil {
		return item, err
	}
	return item, nil
}

func (r *JobRepository) DeleteSubMajorGroup(id uint) error {
	return r.db.Delete(&models.SubMajorGroup{}, id).Error
}

// Minor Group CRUD
func (r *JobRepository) CreateMinorGroup(item *models.MinorGroup) error {
	return r.db.Create(item).Error
}

func (r *JobRepository) GetAllMinorGroups() ([]models.MinorGroup, error) {
	var items []models.MinorGroup
	err := r.db.Preload("SubMajorGroup").Find(&items).Error
	return items, err
}

func (r *JobRepository) GetMinorGroupByID(id uint) (models.MinorGroup, error) {
	var item models.MinorGroup
	err := r.db.Preload("SubMajorGroup").First(&item, id).Error
	return item, err
}

func (r *JobRepository) UpdateMinorGroup(id uint, updates map[string]interface{}) (models.MinorGroup, error) {
	var item models.MinorGroup
	if err := r.db.First(&item, id).Error; err != nil {
		return item, err
	}
	if err := r.db.Model(&item).Updates(updates).Error; err != nil {
		return item, err
	}
	return item, nil
}

func (r *JobRepository) DeleteMinorGroup(id uint) error {
	return r.db.Delete(&models.MinorGroup{}, id).Error
}

// Unit Group CRUD
func (r *JobRepository) CreateUnitGroup(item *models.UnitGroup) error {
	return r.db.Create(item).Error
}

func (r *JobRepository) GetAllUnitGroups() ([]models.UnitGroup, error) {
	var items []models.UnitGroup
	err := r.db.Preload("MinorGroup").Find(&items).Error
	return items, err
}

func (r *JobRepository) GetUnitGroupByID(id uint) (models.UnitGroup, error) {
	var item models.UnitGroup
	err := r.db.Preload("MinorGroup").First(&item, id).Error
	return item, err
}

func (r *JobRepository) UpdateUnitGroup(id uint, updates map[string]interface{}) (models.UnitGroup, error) {
	var item models.UnitGroup
	if err := r.db.First(&item, id).Error; err != nil {
		return item, err
	}
	if err := r.db.Model(&item).Updates(updates).Error; err != nil {
		return item, err
	}
	return item, nil
}

func (r *JobRepository) DeleteUnitGroup(id uint) error {
	return r.db.Delete(&models.UnitGroup{}, id).Error
}

// Ocation Group CRUD
func (r *JobRepository) CreateOccupationGroup(item *models.OccupationGroup) error {
	return r.db.Create(item).Error
}

func (r *JobRepository) GetAllOccupationGroups(limit, offset int) ([]models.OccupationGroup, int64, error) {
	var items []models.OccupationGroup
	var total int64

	if err := r.db.Model(&models.OccupationGroup{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.Preload("UnitGroup")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("id ASC").Find(&items).Error

	return items, total, err
}

func (r *JobRepository) GetOccupationGroupByID(id uint) (models.OccupationGroup, error) {
	var item models.OccupationGroup
	err := r.db.Preload("UnitGroup").First(&item, id).Error
	return item, err
}

func (r *JobRepository) UpdateOccupationGroup(id uint, updates map[string]interface{}) (models.OccupationGroup, error) {
	var item models.OccupationGroup
	if err := r.db.First(&item, id).Error; err != nil {
		return item, err
	}
	if err := r.db.Model(&item).Updates(updates).Error; err != nil {
		return item, err
	}
	return item, nil
}

func (r *JobRepository) DeleteOccupationGroup(id uint) error {
	return r.db.Delete(&models.OccupationGroup{}, id).Error
}
