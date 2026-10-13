package repositories

import (
	"context"
	"marketlens-go-backend/internal/models"
	"time"
)

// Industry levels by parent ids
func (r *Store) GetAllIndustrySectorsForDateRange(ctx context.Context, fromDate, toDate time.Time) ([]models.IndustrySector, error) {
	return newCRUDRepository[models.IndustrySector](r.db).getAllForDateRange(ctx, fromDate, toDate)
}

func (r *Store) GetIndustryDivisionsByIndustrySector(ctx context.Context, industrySectorID uint) ([]models.IndustryDivision, error) {
	return newCRUDRepository[models.IndustryDivision](r.db).getByParent(ctx, "industry_sector_id", industrySectorID)
}

func (r *Store) GetIndustryDivisionsByIndustrySectorForDateRange(ctx context.Context, industrySectorID uint, fromDate, toDate time.Time) ([]models.IndustryDivision, error) {
	return newCRUDRepository[models.IndustryDivision](r.db).getByParentForDateRange(ctx, "industry_sector_id", industrySectorID, fromDate, toDate)
}

func (r *Store) GetIndustryGroupsByIndustryDivision(ctx context.Context, industryDivisionID uint) ([]models.IndustryGroup, error) {
	return newCRUDRepository[models.IndustryGroup](r.db).getByParent(ctx, "industry_division_id", industryDivisionID)
}

func (r *Store) GetIndustryGroupsByIndustryDivisionForDateRange(ctx context.Context, industryDivisionID uint, fromDate, toDate time.Time) ([]models.IndustryGroup, error) {
	return newCRUDRepository[models.IndustryGroup](r.db).getByParentForDateRange(ctx, "industry_division_id", industryDivisionID, fromDate, toDate)
}

func (r *Store) GetIndustryClassesByIndustryGroup(ctx context.Context, industryGroupID uint) ([]models.IndustryClass, error) {
	return newCRUDRepository[models.IndustryClass](r.db).getByParent(ctx, "industry_group_id", industryGroupID)
}

func (r *Store) GetIndustryClassesByIndustryGroupForDateRange(ctx context.Context, industryGroupID uint, fromDate, toDate time.Time) ([]models.IndustryClass, error) {
	return newCRUDRepository[models.IndustryClass](r.db).getByParentForDateRange(ctx, "industry_group_id", industryGroupID, fromDate, toDate)
}

func (r *Store) GetIndustrySubclassesByIndustryClass(ctx context.Context, industryClassID uint) ([]models.IndustrySubclass, error) {
	return newCRUDRepository[models.IndustrySubclass](r.db).getByParent(ctx, "industry_class_id", industryClassID)
}

func (r *Store) GetIndustrySubclassesByIndustryClassForDateRange(ctx context.Context, industryClassID uint, fromDate, toDate time.Time) ([]models.IndustrySubclass, error) {
	return newCRUDRepository[models.IndustrySubclass](r.db).getByParentForDateRange(ctx, "industry_class_id", industryClassID, fromDate, toDate)
}

// Industry sector CRUD
func (r *Store) CreateIndustrySector(ctx context.Context, item *models.IndustrySector) error {
	return newCRUDRepository[models.IndustrySector](r.db).create(ctx, item)
}

func (r *Store) GetAllIndustrySectors(ctx context.Context) ([]models.IndustrySector, error) {
	return newCRUDRepository[models.IndustrySector](r.db).getAll(ctx)
}

func (r *Store) GetIndustrySectorByID(ctx context.Context, id uint) (models.IndustrySector, error) {
	return newCRUDRepository[models.IndustrySector](r.db).getByID(ctx, id)
}

func (r *Store) UpdateIndustrySector(ctx context.Context, id uint, updates map[string]interface{}) (models.IndustrySector, error) {
	return newCRUDRepository[models.IndustrySector](r.db).update(ctx, id, updates)
}

func (r *Store) DeleteIndustrySector(ctx context.Context, id uint) error {
	return newCRUDRepository[models.IndustrySector](r.db).delete(ctx, id)
}

// Industry division CRUD
func (r *Store) CreateIndustryDivision(ctx context.Context, item *models.IndustryDivision) error {
	return newCRUDRepository[models.IndustryDivision](r.db).create(ctx, item)
}

func (r *Store) GetAllIndustryDivisions(ctx context.Context) ([]models.IndustryDivision, error) {
	return newCRUDRepository[models.IndustryDivision](r.db).getAll(ctx, "IndustrySector")
}

func (r *Store) GetIndustryDivisionByID(ctx context.Context, id uint) (models.IndustryDivision, error) {
	return newCRUDRepository[models.IndustryDivision](r.db).getByID(ctx, id, "IndustrySector")
}

func (r *Store) UpdateIndustryDivision(ctx context.Context, id uint, updates map[string]interface{}) (models.IndustryDivision, error) {
	return newCRUDRepository[models.IndustryDivision](r.db).update(ctx, id, updates)
}

func (r *Store) DeleteIndustryDivision(ctx context.Context, id uint) error {
	return newCRUDRepository[models.IndustryDivision](r.db).delete(ctx, id)
}

// Industry group CRUD
func (r *Store) CreateIndustryGroup(ctx context.Context, item *models.IndustryGroup) error {
	return newCRUDRepository[models.IndustryGroup](r.db).create(ctx, item)
}

func (r *Store) GetAllIndustryGroups(ctx context.Context) ([]models.IndustryGroup, error) {
	return newCRUDRepository[models.IndustryGroup](r.db).getAll(ctx, "IndustryDivision")
}

func (r *Store) GetIndustryGroupByID(ctx context.Context, id uint) (models.IndustryGroup, error) {
	return newCRUDRepository[models.IndustryGroup](r.db).getByID(ctx, id, "IndustryDivision")
}

func (r *Store) UpdateIndustryGroup(ctx context.Context, id uint, updates map[string]interface{}) (models.IndustryGroup, error) {
	return newCRUDRepository[models.IndustryGroup](r.db).update(ctx, id, updates)
}

func (r *Store) DeleteIndustryGroup(ctx context.Context, id uint) error {
	return newCRUDRepository[models.IndustryGroup](r.db).delete(ctx, id)
}

// Industry class CRUD
func (r *Store) CreateIndustryClass(ctx context.Context, item *models.IndustryClass) error {
	return newCRUDRepository[models.IndustryClass](r.db).create(ctx, item)
}

func (r *Store) GetAllIndustryClasses(ctx context.Context) ([]models.IndustryClass, error) {
	return newCRUDRepository[models.IndustryClass](r.db).getAll(ctx, "IndustryGroup")
}

func (r *Store) GetIndustryClassByID(ctx context.Context, id uint) (models.IndustryClass, error) {
	return newCRUDRepository[models.IndustryClass](r.db).getByID(ctx, id, "IndustryGroup")
}

func (r *Store) UpdateIndustryClass(ctx context.Context, id uint, updates map[string]interface{}) (models.IndustryClass, error) {
	return newCRUDRepository[models.IndustryClass](r.db).update(ctx, id, updates)
}

func (r *Store) DeleteIndustryClass(ctx context.Context, id uint) error {
	return newCRUDRepository[models.IndustryClass](r.db).delete(ctx, id)
}

// Industry sub class CRUD
func (r *Store) CreateIndustrySubclass(ctx context.Context, item *models.IndustrySubclass) error {
	return newCRUDRepository[models.IndustrySubclass](r.db).create(ctx, item)
}

func (r *Store) GetAllIndustrySubclasses(ctx context.Context, limit, offset int) ([]models.IndustrySubclass, int64, error) {
	return newCRUDRepository[models.IndustrySubclass](r.db).getAllPaginated(ctx, limit, offset, "IndustryClass")
}

func (r *Store) GetIndustrySubclassByID(ctx context.Context, id uint) (models.IndustrySubclass, error) {
	return newCRUDRepository[models.IndustrySubclass](r.db).getByID(ctx, id, "IndustryClass")
}

func (r *Store) UpdateIndustrySubclass(ctx context.Context, id uint, updates map[string]interface{}) (models.IndustrySubclass, error) {
	return newCRUDRepository[models.IndustrySubclass](r.db).update(ctx, id, updates)
}

func (r *Store) DeleteIndustrySubclass(ctx context.Context, id uint) error {
	return newCRUDRepository[models.IndustrySubclass](r.db).delete(ctx, id)
}
