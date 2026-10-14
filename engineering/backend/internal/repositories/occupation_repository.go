package repositories

import (
	"context"
	"marketlens-go-backend/internal/models"
	"time"
)

// Methods needed for multi-level filtering in the crawler
// Occupation levels by parent ids
func (r *Store) GetAllMajorGroupsForDateRange(ctx context.Context, fromDate, toDate time.Time) ([]models.MajorGroup, error) {
	return newCRUDRepository[models.MajorGroup](r.db).getAll(ctx, 0, 0, &fromDate, &toDate)
}

func (r *Store) GetSubMajorGroupsByMajorGroup(ctx context.Context, majorGroupID uint) ([]models.SubMajorGroup, error) {
	return newCRUDRepository[models.SubMajorGroup](r.db).getByParent(ctx, "major_group_id", majorGroupID, nil, nil)
}

func (r *Store) GetSubMajorGroupsByMajorGroupForDateRange(ctx context.Context, majorGroupID uint, fromDate, toDate time.Time) ([]models.SubMajorGroup, error) {
	return newCRUDRepository[models.SubMajorGroup](r.db).getByParent(ctx, "major_group_id", majorGroupID, &fromDate, &toDate)
}

func (r *Store) GetMinorGroupsBySubMajorGroup(ctx context.Context, subMajorGroupID uint) ([]models.MinorGroup, error) {
	return newCRUDRepository[models.MinorGroup](r.db).getByParent(ctx, "sub_major_group_id", subMajorGroupID, nil, nil)
}

func (r *Store) GetMinorGroupsBySubMajorGroupForDateRange(ctx context.Context, subMajorGroupID uint, fromDate, toDate time.Time) ([]models.MinorGroup, error) {
	return newCRUDRepository[models.MinorGroup](r.db).getByParent(ctx, "sub_major_group_id", subMajorGroupID, &fromDate, &toDate)
}

func (r *Store) GetUnitGroupsByMinorGroup(ctx context.Context, minorGroupID uint) ([]models.UnitGroup, error) {
	return newCRUDRepository[models.UnitGroup](r.db).getByParent(ctx, "minor_group_id", minorGroupID, nil, nil)
}

func (r *Store) GetUnitGroupsByMinorGroupForDateRange(ctx context.Context, minorGroupID uint, fromDate, toDate time.Time) ([]models.UnitGroup, error) {
	return newCRUDRepository[models.UnitGroup](r.db).getByParent(ctx, "minor_group_id", minorGroupID, &fromDate, &toDate)
}

func (r *Store) GetOccupationGroupsByUnitGroup(ctx context.Context, unitGroupID uint) ([]models.OccupationGroup, error) {
	return newCRUDRepository[models.OccupationGroup](r.db).getByParent(ctx, "unit_group_id", unitGroupID, nil, nil)
}

func (r *Store) GetOccupationGroupsByUnitGroupForDateRange(ctx context.Context, unitGroupID uint, fromDate, toDate time.Time) ([]models.OccupationGroup, error) {
	return newCRUDRepository[models.OccupationGroup](r.db).getByParent(ctx, "unit_group_id", unitGroupID, &fromDate, &toDate)
}

// Major Group CRUD
func (r *Store) CreateMajorGroup(ctx context.Context, item *models.MajorGroup) error {
	return newCRUDRepository[models.MajorGroup](r.db).create(ctx, item)
}

func (r *Store) GetAllMajorGroups(ctx context.Context) ([]models.MajorGroup, error) {
	return newCRUDRepository[models.MajorGroup](r.db).getAll(ctx, 0, 0, nil, nil)
}

func (r *Store) GetMajorGroupByID(ctx context.Context, id uint) (models.MajorGroup, error) {
	return newCRUDRepository[models.MajorGroup](r.db).getByID(ctx, id)
}

func (r *Store) UpdateMajorGroup(ctx context.Context, id uint, updates map[string]interface{}) (models.MajorGroup, error) {
	return newCRUDRepository[models.MajorGroup](r.db).update(ctx, id, updates)
}

func (r *Store) DeleteMajorGroup(ctx context.Context, id uint) error {
	return newCRUDRepository[models.MajorGroup](r.db).delete(ctx, id)
}

// Sub Major Group CRUD
func (r *Store) CreateSubMajorGroup(ctx context.Context, item *models.SubMajorGroup) error {
	return newCRUDRepository[models.SubMajorGroup](r.db).create(ctx, item)
}

func (r *Store) GetAllSubMajorGroups(ctx context.Context) ([]models.SubMajorGroup, error) {
	return newCRUDRepository[models.SubMajorGroup](r.db).getAll(ctx, 0, 0, nil, nil, "MajorGroup")
}

func (r *Store) GetSubMajorGroupByID(ctx context.Context, id uint) (models.SubMajorGroup, error) {
	return newCRUDRepository[models.SubMajorGroup](r.db).getByID(ctx, id, "MajorGroup")
}

func (r *Store) UpdateSubMajorGroup(ctx context.Context, id uint, updates map[string]interface{}) (models.SubMajorGroup, error) {
	return newCRUDRepository[models.SubMajorGroup](r.db).update(ctx, id, updates)
}

func (r *Store) DeleteSubMajorGroup(ctx context.Context, id uint) error {
	return newCRUDRepository[models.SubMajorGroup](r.db).delete(ctx, id)
}

// Minor Group CRUD
func (r *Store) CreateMinorGroup(ctx context.Context, item *models.MinorGroup) error {
	return newCRUDRepository[models.MinorGroup](r.db).create(ctx, item)
}

func (r *Store) GetAllMinorGroups(ctx context.Context) ([]models.MinorGroup, error) {
	return newCRUDRepository[models.MinorGroup](r.db).getAll(ctx, 0, 0, nil, nil, "SubMajorGroup")
}

func (r *Store) GetMinorGroupByID(ctx context.Context, id uint) (models.MinorGroup, error) {
	return newCRUDRepository[models.MinorGroup](r.db).getByID(ctx, id, "SubMajorGroup")
}

func (r *Store) UpdateMinorGroup(ctx context.Context, id uint, updates map[string]interface{}) (models.MinorGroup, error) {
	return newCRUDRepository[models.MinorGroup](r.db).update(ctx, id, updates)
}

func (r *Store) DeleteMinorGroup(ctx context.Context, id uint) error {
	return newCRUDRepository[models.MinorGroup](r.db).delete(ctx, id)
}

// Unit Group CRUD
func (r *Store) CreateUnitGroup(ctx context.Context, item *models.UnitGroup) error {
	return newCRUDRepository[models.UnitGroup](r.db).create(ctx, item)
}

func (r *Store) GetAllUnitGroups(ctx context.Context) ([]models.UnitGroup, error) {
	return newCRUDRepository[models.UnitGroup](r.db).getAll(ctx, 0, 0, nil, nil, "MinorGroup")
}

func (r *Store) GetUnitGroupByID(ctx context.Context, id uint) (models.UnitGroup, error) {
	return newCRUDRepository[models.UnitGroup](r.db).getByID(ctx, id, "MinorGroup")
}

func (r *Store) UpdateUnitGroup(ctx context.Context, id uint, updates map[string]interface{}) (models.UnitGroup, error) {
	return newCRUDRepository[models.UnitGroup](r.db).update(ctx, id, updates)
}

func (r *Store) DeleteUnitGroup(ctx context.Context, id uint) error {
	return newCRUDRepository[models.UnitGroup](r.db).delete(ctx, id)
}

// Occupation Group CRUD
func (r *Store) CreateOccupationGroup(ctx context.Context, item *models.OccupationGroup) error {
	return newCRUDRepository[models.OccupationGroup](r.db).create(ctx, item)
}

func (r *Store) GetAllOccupationGroups(ctx context.Context, limit, offset int) ([]models.OccupationGroup, error) {
	return newCRUDRepository[models.OccupationGroup](r.db).getAll(ctx, limit, offset, nil, nil, "UnitGroup")
}

func (r *Store) GetOccupationGroupByID(ctx context.Context, id uint) (models.OccupationGroup, error) {
	return newCRUDRepository[models.OccupationGroup](r.db).getByID(ctx, id, "UnitGroup")
}

func (r *Store) UpdateOccupationGroup(ctx context.Context, id uint, updates map[string]interface{}) (models.OccupationGroup, error) {
	return newCRUDRepository[models.OccupationGroup](r.db).update(ctx, id, updates)
}

func (r *Store) DeleteOccupationGroup(ctx context.Context, id uint) error {
	return newCRUDRepository[models.OccupationGroup](r.db).delete(ctx, id)
}
