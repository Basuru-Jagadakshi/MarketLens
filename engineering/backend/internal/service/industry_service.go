package service

import (
	"context"
	"marketlens-go-backend/internal/models"
	"marketlens-go-backend/internal/repositories"

	"gorm.io/gorm"
)

func NewIndustrySectorService(store *repositories.Store) *CRUDService[models.IndustrySector] {
	return NewCRUDService(
		store.CreateIndustrySector,
		store.GetAllIndustrySectors,
		store.GetIndustrySectorByID,
		store.UpdateIndustrySector,
		store.DeleteIndustrySector,
		func(_ context.Context, item *models.IndustrySector) error {
			return validateNameAndCode(item.Name, item.Code)
		},
		WithDateRangeListing(store.GetAllIndustrySectorsForDateRange),
	)
}

func NewIndustryDivisionService(store *repositories.Store) *CRUDService[models.IndustryDivision] {
	return NewCRUDService(
		store.CreateIndustryDivision,
		store.GetAllIndustryDivisions,
		store.GetIndustryDivisionByID,
		store.UpdateIndustryDivision,
		store.DeleteIndustryDivision,
		func(ctx context.Context, item *models.IndustryDivision) error {
			if err := validateNameAndCode(item.Name, item.Code); err != nil {
				return err
			}
			_, err := store.GetIndustrySectorByID(ctx, item.IndustrySectorID)
			return validateParentExists(err, gorm.ErrRecordNotFound, "industry_sector_id", item.IndustrySectorID)
		},
		WithParentLookups(store.GetIndustryDivisionsByIndustrySector, store.GetIndustryDivisionsByIndustrySectorForDateRange),
	)
}

func NewIndustryGroupService(store *repositories.Store) *CRUDService[models.IndustryGroup] {
	return NewCRUDService(
		store.CreateIndustryGroup,
		store.GetAllIndustryGroups,
		store.GetIndustryGroupByID,
		store.UpdateIndustryGroup,
		store.DeleteIndustryGroup,
		func(ctx context.Context, item *models.IndustryGroup) error {
			if err := validateNameAndCode(item.Name, item.Code); err != nil {
				return err
			}
			_, err := store.GetIndustryDivisionByID(ctx, item.IndustryDivisionID)
			return validateParentExists(err, gorm.ErrRecordNotFound, "industry_division_id", item.IndustryDivisionID)
		},
		WithParentLookups(store.GetIndustryGroupsByIndustryDivision, store.GetIndustryGroupsByIndustryDivisionForDateRange),
	)
}

func NewIndustryClassService(store *repositories.Store) *CRUDService[models.IndustryClass] {
	return NewCRUDService(
		store.CreateIndustryClass,
		store.GetAllIndustryClasses,
		store.GetIndustryClassByID,
		store.UpdateIndustryClass,
		store.DeleteIndustryClass,
		func(ctx context.Context, item *models.IndustryClass) error {
			if err := validateNameAndCode(item.Name, item.Code); err != nil {
				return err
			}
			_, err := store.GetIndustryGroupByID(ctx, item.IndustryGroupID)
			return validateParentExists(err, gorm.ErrRecordNotFound, "industry_group_id", item.IndustryGroupID)
		},
		WithParentLookups(store.GetIndustryClassesByIndustryGroup, store.GetIndustryClassesByIndustryGroupForDateRange),
	)
}

func NewIndustrySubclassService(store *repositories.Store) *CRUDService[models.IndustrySubclass] {
	return NewCRUDService(
		store.CreateIndustrySubclass,
		func(ctx context.Context) ([]models.IndustrySubclass, error) {
			items, _, err := store.GetAllIndustrySubclasses(ctx, 0, 0)
			return items, err
		},
		store.GetIndustrySubclassByID,
		store.UpdateIndustrySubclass,
		store.DeleteIndustrySubclass,
		func(ctx context.Context, item *models.IndustrySubclass) error {
			if err := validateNameAndCode(item.Name, item.Code); err != nil {
				return err
			}
			_, err := store.GetIndustryClassByID(ctx, item.IndustryClassID)
			return validateParentExists(err, gorm.ErrRecordNotFound, "industry_class_id", item.IndustryClassID)
		},
		WithParentLookups(store.GetIndustrySubclassesByIndustryClass, store.GetIndustrySubclassesByIndustryClassForDateRange),
		WithPagination(store.GetAllIndustrySubclasses),
	)
}
