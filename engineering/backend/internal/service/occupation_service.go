package service

import (
	"context"
	"marketlens-go-backend/internal/models"
	"marketlens-go-backend/internal/repositories"
	"time"

	"gorm.io/gorm"
)

func NewMajorGroupService(store *repositories.Store) *CRUDService[models.MajorGroup] {
	return NewCRUDService(
		store.CreateMajorGroup,
		func(ctx context.Context, _, _ int, fromDate, toDate *time.Time) ([]models.MajorGroup, error) {
			if fromDate != nil && toDate != nil {
				return store.GetAllMajorGroupsForDateRange(ctx, *fromDate, *toDate)
			}
			return store.GetAllMajorGroups(ctx)
		},
		store.GetMajorGroupByID,
		store.UpdateMajorGroup,
		store.DeleteMajorGroup,
		func(_ context.Context, item *models.MajorGroup) error {
			return validateNameAndCode(item.Name, item.Code)
		},
	)
}

func NewSubMajorGroupService(store *repositories.Store) *CRUDService[models.SubMajorGroup] {
	return NewCRUDService(
		store.CreateSubMajorGroup,
		func(ctx context.Context, _, _ int, _, _ *time.Time) ([]models.SubMajorGroup, error) {
			return store.GetAllSubMajorGroups(ctx)
		},
		store.GetSubMajorGroupByID,
		store.UpdateSubMajorGroup,
		store.DeleteSubMajorGroup,
		func(ctx context.Context, item *models.SubMajorGroup) error {
			if err := validateNameAndCode(item.Name, item.Code); err != nil {
				return err
			}
			_, err := store.GetMajorGroupByID(ctx, item.MajorGroupID)
			return validateParentExists(err, gorm.ErrRecordNotFound, "major_group_id", item.MajorGroupID)
		},
		WithParentLookups(func(ctx context.Context, parentID uint, fromDate, toDate *time.Time) ([]models.SubMajorGroup, error) {
			if fromDate != nil && toDate != nil {
				return store.GetSubMajorGroupsByMajorGroupForDateRange(ctx, parentID, *fromDate, *toDate)
			}
			return store.GetSubMajorGroupsByMajorGroup(ctx, parentID)
		}),
	)
}

func NewMinorGroupService(store *repositories.Store) *CRUDService[models.MinorGroup] {
	return NewCRUDService(
		store.CreateMinorGroup,
		func(ctx context.Context, _, _ int, _, _ *time.Time) ([]models.MinorGroup, error) {
			return store.GetAllMinorGroups(ctx)
		},
		store.GetMinorGroupByID,
		store.UpdateMinorGroup,
		store.DeleteMinorGroup,
		func(ctx context.Context, item *models.MinorGroup) error {
			if err := validateNameAndCode(item.Name, item.Code); err != nil {
				return err
			}
			_, err := store.GetSubMajorGroupByID(ctx, item.SubMajorGroupID)
			return validateParentExists(err, gorm.ErrRecordNotFound, "sub_major_group_id", item.SubMajorGroupID)
		},
		WithParentLookups(func(ctx context.Context, parentID uint, fromDate, toDate *time.Time) ([]models.MinorGroup, error) {
			if fromDate != nil && toDate != nil {
				return store.GetMinorGroupsBySubMajorGroupForDateRange(ctx, parentID, *fromDate, *toDate)
			}
			return store.GetMinorGroupsBySubMajorGroup(ctx, parentID)
		}),
	)
}

func NewUnitGroupService(store *repositories.Store) *CRUDService[models.UnitGroup] {
	return NewCRUDService(
		store.CreateUnitGroup,
		func(ctx context.Context, _, _ int, _, _ *time.Time) ([]models.UnitGroup, error) {
			return store.GetAllUnitGroups(ctx)
		},
		store.GetUnitGroupByID,
		store.UpdateUnitGroup,
		store.DeleteUnitGroup,
		func(ctx context.Context, item *models.UnitGroup) error {
			if err := validateNameAndCode(item.Name, item.Code); err != nil {
				return err
			}
			_, err := store.GetMinorGroupByID(ctx, item.MinorGroupID)
			return validateParentExists(err, gorm.ErrRecordNotFound, "minor_group_id", item.MinorGroupID)
		},
		WithParentLookups(func(ctx context.Context, parentID uint, fromDate, toDate *time.Time) ([]models.UnitGroup, error) {
			if fromDate != nil && toDate != nil {
				return store.GetUnitGroupsByMinorGroupForDateRange(ctx, parentID, *fromDate, *toDate)
			}
			return store.GetUnitGroupsByMinorGroup(ctx, parentID)
		}),
	)
}

func NewOccupationGroupService(store *repositories.Store) *CRUDService[models.OccupationGroup] {
	return NewCRUDService(
		store.CreateOccupationGroup,
		func(ctx context.Context, limit, offset int, _, _ *time.Time) ([]models.OccupationGroup, error) {
			return store.GetAllOccupationGroups(ctx, limit, offset)
		},
		store.GetOccupationGroupByID,
		store.UpdateOccupationGroup,
		store.DeleteOccupationGroup,
		func(ctx context.Context, item *models.OccupationGroup) error {
			if err := validateNameAndCode(item.Name, item.Code); err != nil {
				return err
			}
			_, err := store.GetUnitGroupByID(ctx, item.UnitGroupID)
			return validateParentExists(err, gorm.ErrRecordNotFound, "unit_group_id", item.UnitGroupID)
		},
		WithParentLookups(func(ctx context.Context, parentID uint, fromDate, toDate *time.Time) ([]models.OccupationGroup, error) {
			if fromDate != nil && toDate != nil {
				return store.GetOccupationGroupsByUnitGroupForDateRange(ctx, parentID, *fromDate, *toDate)
			}
			return store.GetOccupationGroupsByUnitGroup(ctx, parentID)
		}),
	)
}
