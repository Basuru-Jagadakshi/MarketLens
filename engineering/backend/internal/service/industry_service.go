package service

import (
	"context"
	"marketlens-go-backend/internal/models"
	"marketlens-go-backend/internal/repositories"
	"time"

	"gorm.io/gorm"
)

func NewIndustrySectorService(store *repositories.Store) *CRUDService[models.IndustrySector] {
	return NewCRUDService(
		store.CreateIndustrySector,
		func(ctx context.Context, _, _ int, fromDate, toDate *time.Time) ([]models.IndustrySector, error) {
			if fromDate != nil && toDate != nil {
				return store.GetAllIndustrySectorsForDateRange(ctx, *fromDate, *toDate)
			}
			return store.GetAllIndustrySectors(ctx)
		},
		store.GetIndustrySectorByID,
		store.UpdateIndustrySector,
		store.DeleteIndustrySector,
		func(_ context.Context, item *models.IndustrySector) error {
			return validateNameAndCode(item.Name, item.Code)
		},
	)
}

func NewIndustryDivisionService(store *repositories.Store) *CRUDService[models.IndustryDivision] {
	return NewCRUDService(
		store.CreateIndustryDivision,
		func(ctx context.Context, _, _ int, _, _ *time.Time) ([]models.IndustryDivision, error) {
			return store.GetAllIndustryDivisions(ctx)
		},
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
		WithParentLookups(func(ctx context.Context, parentID uint, fromDate, toDate *time.Time) ([]models.IndustryDivision, error) {
			if fromDate != nil && toDate != nil {
				return store.GetIndustryDivisionsByIndustrySectorForDateRange(ctx, parentID, *fromDate, *toDate)
			}
			return store.GetIndustryDivisionsByIndustrySector(ctx, parentID)
		}),
	)
}

func NewIndustryGroupService(store *repositories.Store) *CRUDService[models.IndustryGroup] {
	return NewCRUDService(
		store.CreateIndustryGroup,
		func(ctx context.Context, _, _ int, _, _ *time.Time) ([]models.IndustryGroup, error) {
			return store.GetAllIndustryGroups(ctx)
		},
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
		WithParentLookups(func(ctx context.Context, parentID uint, fromDate, toDate *time.Time) ([]models.IndustryGroup, error) {
			if fromDate != nil && toDate != nil {
				return store.GetIndustryGroupsByIndustryDivisionForDateRange(ctx, parentID, *fromDate, *toDate)
			}
			return store.GetIndustryGroupsByIndustryDivision(ctx, parentID)
		}),
	)
}

func NewIndustryClassService(store *repositories.Store) *CRUDService[models.IndustryClass] {
	return NewCRUDService(
		store.CreateIndustryClass,
		func(ctx context.Context, _, _ int, _, _ *time.Time) ([]models.IndustryClass, error) {
			return store.GetAllIndustryClasses(ctx)
		},
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
		WithParentLookups(func(ctx context.Context, parentID uint, fromDate, toDate *time.Time) ([]models.IndustryClass, error) {
			if fromDate != nil && toDate != nil {
				return store.GetIndustryClassesByIndustryGroupForDateRange(ctx, parentID, *fromDate, *toDate)
			}
			return store.GetIndustryClassesByIndustryGroup(ctx, parentID)
		}),
	)
}

func NewIndustrySubclassService(store *repositories.Store) *CRUDService[models.IndustrySubclass] {
	return NewCRUDService(
		store.CreateIndustrySubclass,
		func(ctx context.Context, limit, offset int, _, _ *time.Time) ([]models.IndustrySubclass, error) {
			return store.GetAllIndustrySubclasses(ctx, limit, offset)
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
		WithParentLookups(func(ctx context.Context, parentID uint, fromDate, toDate *time.Time) ([]models.IndustrySubclass, error) {
			if fromDate != nil && toDate != nil {
				return store.GetIndustrySubclassesByIndustryClassForDateRange(ctx, parentID, *fromDate, *toDate)
			}
			return store.GetIndustrySubclassesByIndustryClass(ctx, parentID)
		}),
	)
}
