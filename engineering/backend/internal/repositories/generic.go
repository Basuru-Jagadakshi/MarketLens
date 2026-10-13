package repositories

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// crudRepository is a generic, GORM-backed repository for the flat-lookup
// and hierarchy-level entities: one table, an auto-increment id, and the
// same create/list/get/update/delete shape. T is a model from
// internal/models (it only needs its own TableName() method, which every
// model there already has). Every Store method for these entities is a
// one-line delegate to one of these, so the actual query logic — ctx
// propagation, deterministic ordering, and mapping "no row matched" to
// gorm.ErrRecordNotFound instead of silently succeeding — lives in exactly
// one place instead of being copy-pasted per entity.
type crudRepository[T any] struct {
	db *gorm.DB
}

func newCRUDRepository[T any](db *gorm.DB) crudRepository[T] {
	return crudRepository[T]{db: db}
}

func (c crudRepository[T]) create(ctx context.Context, item *T) error {
	return c.db.WithContext(ctx).Create(item).Error
}

// getAll lists every row, ordered by id so callers (e.g. dropdowns) get a
// stable order instead of whatever the database happens to return.
// preloads are relation names to eager-load (e.g. "MajorGroup"); pass none
// for a root-level entity with no parent relation.
func (c crudRepository[T]) getAll(ctx context.Context, preloads ...string) ([]T, error) {
	var items []T
	q := c.db.WithContext(ctx)
	for _, p := range preloads {
		q = q.Preload(p)
	}
	err := q.Order("id ASC").Find(&items).Error
	return items, err
}

// getAllPaginated is getAll with an optional limit/offset and a total
// count, for the handful of entities large enough to paginate.
func (c crudRepository[T]) getAllPaginated(ctx context.Context, limit, offset int, preloads ...string) ([]T, int64, error) {
	var zero T
	var total int64
	if err := c.db.WithContext(ctx).Model(&zero).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	q := c.db.WithContext(ctx)
	for _, p := range preloads {
		q = q.Preload(p)
	}
	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}

	var items []T
	err := q.Order("id ASC").Find(&items).Error
	return items, total, err
}

func (c crudRepository[T]) getByID(ctx context.Context, id uint, preloads ...string) (T, error) {
	var item T
	q := c.db.WithContext(ctx)
	for _, p := range preloads {
		q = q.Preload(p)
	}
	err := q.First(&item, id).Error
	return item, err
}

// update applies a partial update atomically (no separate "does it exist"
// read racing against a concurrent delete) and maps "no row matched" to
// gorm.ErrRecordNotFound so the caller can map it to an HTTP 404 instead of
// a 500. The updated row is then reloaded so callers get the current value
// back, not the pre-update snapshot.
func (c crudRepository[T]) update(ctx context.Context, id uint, updates map[string]interface{}) (T, error) {
	var zero T
	res := c.db.WithContext(ctx).Model(&zero).Where("id = ?", id).Updates(updates)
	if res.Error != nil {
		return zero, fmt.Errorf("update %T %d: %w", zero, id, res.Error)
	}
	if res.RowsAffected == 0 {
		return zero, gorm.ErrRecordNotFound
	}

	var item T
	if err := c.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return item, fmt.Errorf("reload %T %d: %w", zero, id, err)
	}
	return item, nil
}

// delete maps "no row matched" (missing or already-deleted id) to
// gorm.ErrRecordNotFound instead of silently reporting success.
func (c crudRepository[T]) delete(ctx context.Context, id uint) error {
	var zero T
	res := c.db.WithContext(ctx).Delete(&zero, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// getByParent lists the rows whose parentColumn (e.g. "major_group_id")
// equals parentID — the child-levels-of-a-parent lookups used by the
// hierarchy (occupation/industry classification) tables.
func (c crudRepository[T]) getByParent(ctx context.Context, parentColumn string, parentID uint) ([]T, error) {
	var items []T
	err := c.db.WithContext(ctx).Where(parentColumn+" = ?", parentID).Find(&items).Error
	return items, err
}

// getByParentForDateRange is getByParent scoped to rows that existed at
// some point during [fromDate, toDate] (Unscoped so soft-deleted rows that
// were still live during the window are included).
func (c crudRepository[T]) getByParentForDateRange(ctx context.Context, parentColumn string, parentID uint, fromDate, toDate time.Time) ([]T, error) {
	var items []T
	endOfToDate := toDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	err := c.db.WithContext(ctx).
		Unscoped().
		Where(parentColumn+" = ?", parentID).
		Where("created_at <= ?", endOfToDate).
		Where("deleted_at IS NULL OR deleted_at >= ?", fromDate).
		Find(&items).Error
	return items, err
}

// getAllForDateRange is getByParentForDateRange without a parent filter,
// for the root level of a hierarchy (e.g. MajorGroup, IndustrySector).
func (c crudRepository[T]) getAllForDateRange(ctx context.Context, fromDate, toDate time.Time) ([]T, error) {
	var items []T
	endOfToDate := toDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	err := c.db.WithContext(ctx).
		Unscoped().
		Where("created_at <= ?", endOfToDate).
		Where("deleted_at IS NULL OR deleted_at >= ?", fromDate).
		Find(&items).Error
	return items, err
}
