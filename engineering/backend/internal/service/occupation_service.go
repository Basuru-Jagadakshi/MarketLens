package service

import (
	"context"
	"marketlens-go-backend/internal/models"
	"marketlens-go-backend/internal/repositories"

	"gorm.io/gorm"
)

func NewMajorGroupService(store *repositories.Store) *CRUDService[models.MajorGroup] {
	return NewCRUDService(
		store.CreateMajorGroup,
		store.GetAllMajorGroups,
		store.GetMajorGroupByID,
		store.UpdateMajorGroup,
		store.DeleteMajorGroup,
		func(_ context.Context, item *models.MajorGroup) error {
			return validateNameAndCode(item.Name, item.Code)
		},
		WithDateRangeListing(store.GetAllMajorGroupsForDateRange),
	)
}

func NewSubMajorGroupService(store *repositories.Store) *CRUDService[models.SubMajorGroup] {
	return NewCRUDService(
		store.CreateSubMajorGroup,
		store.GetAllSubMajorGroups,
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
		WithParentLookups(store.GetSubMajorGroupsByMajorGroup, store.GetSubMajorGroupsByMajorGroupForDateRange),
	)
}

func NewMinorGroupService(store *repositories.Store) *CRUDService[models.MinorGroup] {
	return NewCRUDService(
		store.CreateMinorGroup,
		store.GetAllMinorGroups,
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
		WithParentLookups(store.GetMinorGroupsBySubMajorGroup, store.GetMinorGroupsBySubMajorGroupForDateRange),
	)
}

func NewUnitGroupService(store *repositories.Store) *CRUDService[models.UnitGroup] {
	return NewCRUDService(
		store.CreateUnitGroup,
		store.GetAllUnitGroups,
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
		WithParentLookups(store.GetUnitGroupsByMinorGroup, store.GetUnitGroupsByMinorGroupForDateRange),
	)
}

func NewOccupationGroupService(store *repositories.Store) *CRUDService[models.OccupationGroup] {
	return NewCRUDService(
		store.CreateOccupationGroup,
		func(ctx context.Context) ([]models.OccupationGroup, error) {
			items, _, err := store.GetAllOccupationGroups(ctx, 0, 0)
			return items, err
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
		WithParentLookups(store.GetOccupationGroupsByUnitGroup, store.GetOccupationGroupsByUnitGroupForDateRange),
		WithPagination(store.GetAllOccupationGroups),
	)
}
