package service

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrValidation marks a request as rejected by a business rule (a required
// field was empty, a referenced parent id doesn't exist, ...) rather than
// failing because of the database. Handlers map it to 400 Bad Request
// instead of 500.
var ErrValidation = errors.New("validation failed")

// ValidationErrorf builds an error that errors.Is(err, ErrValidation) finds.
func ValidationErrorf(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrValidation, fmt.Sprintf(format, args...))
}

// CRUDService fronts a Store's generated methods for one flat-lookup or
// hierarchy entity. Every field is a Store method value; the optional ones
// (getByParent) are nil unless the entity actually has that shape (see each
// entity's constructor in *_service.go) and their accessor returns a clear
// error if called on an entity that doesn't support them, instead of
// panicking.
type CRUDService[T any] struct {
	create func(ctx context.Context, item *T) error
	getAll  func(ctx context.Context, limit, offset int, fromDate, toDate *time.Time) ([]T, error)
	getByID func(ctx context.Context, id uint) (T, error)
	update  func(ctx context.Context, id uint, updates map[string]interface{}) (T, error)
	delete  func(ctx context.Context, id uint) error
	getByParent func(ctx context.Context, parentID uint, fromDate, toDate *time.Time) ([]T, error)

	validate func(ctx context.Context, item *T) error
}

type CRUDOption[T any] func(*CRUDService[T])

// WithParentLookups gives a hierarchy child-level entity its
// "children of a parent id" reads.
func WithParentLookups[T any](getByParent func(ctx context.Context, parentID uint, fromDate, toDate *time.Time) ([]T, error)) CRUDOption[T] {
	return func(s *CRUDService[T]) { s.getByParent = getByParent }
}

func NewCRUDService[T any](
	create func(ctx context.Context, item *T) error,
	getAll func(ctx context.Context, limit, offset int, fromDate, toDate *time.Time) ([]T, error),
	getByID func(ctx context.Context, id uint) (T, error),
	update func(ctx context.Context, id uint, updates map[string]interface{}) (T, error),
	del func(ctx context.Context, id uint) error,
	validate func(ctx context.Context, item *T) error,
	opts ...CRUDOption[T],
) *CRUDService[T] {
	s := &CRUDService[T]{create: create, getAll: getAll, getByID: getByID, update: update, delete: del, validate: validate}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *CRUDService[T]) Create(ctx context.Context, item *T) error {
	if s.validate != nil {
		if err := s.validate(ctx, item); err != nil {
			return err
		}
	}
	return s.create(ctx, item)
}


func (s *CRUDService[T]) GetAll(ctx context.Context, limit, offset int, fromDate, toDate *time.Time) ([]T, error) {
	if fromDate != nil && toDate != nil {
		if err := validateDateRange(*fromDate, *toDate); err != nil {
			return nil, err
		}
	}
	return s.getAll(ctx, limit, offset, fromDate, toDate)
}

func (s *CRUDService[T]) GetByID(ctx context.Context, id uint) (T, error) {
	return s.getByID(ctx, id)
}

func (s *CRUDService[T]) Update(ctx context.Context, id uint, updates map[string]interface{}) (T, error) {
	return s.update(ctx, id, updates)
}

func (s *CRUDService[T]) Delete(ctx context.Context, id uint) error {
	return s.delete(ctx, id)
}

// GetByParent lists the children of a parent id, optionally scoped to a
// date range (fromDate/toDate nil means no filter). Only valid for a
// hierarchy child-level entity constructed with WithParentLookups.
func (s *CRUDService[T]) GetByParent(ctx context.Context, parentID uint, fromDate, toDate *time.Time) ([]T, error) {
	if s.getByParent == nil {
		return nil, fmt.Errorf("GetByParent is not supported for this entity")
	}
	if fromDate != nil && toDate != nil {
		if err := validateDateRange(*fromDate, *toDate); err != nil {
			return nil, err
		}
	}
	return s.getByParent(ctx, parentID, fromDate, toDate)
}

// validateDateRange is the one shared "to must not be before from" rule
// used by every date-range query across the service package.
func validateDateRange(fromDate, toDate time.Time) error {
	if toDate.Before(fromDate) {
		return ValidationErrorf("to-date must not be before from-date")
	}
	return nil
}

// requireNonEmptyField builds a Create validator for the common flat-lookup
// shape: one required display string. get reads that field off the item
// (e.g. func(e *models.EmploymentSector) string { return e.Sector }).
func requireNonEmptyField[T any](fieldName string, get func(*T) string) func(context.Context, *T) error {
	return func(_ context.Context, item *T) error {
		if get(item) == "" {
			return ValidationErrorf("%s is required", fieldName)
		}
		return nil
	}
}

// validateParentExists turns a parent lookup's "not found" into a
// validation error naming the offending field; any other error (a real DB
// problem) passes through unchanged so it surfaces as a 500, not a
// misleading "doesn't exist". It takes the already-performed lookup's
// error rather than doing the lookup itself, since each entity's parent
// type differs.
func validateParentExists(lookupErr error, notFound error, field string, id uint) error {
	if errors.Is(lookupErr, notFound) {
		return ValidationErrorf("%s %d does not exist", field, id)
	}
	return lookupErr
}

// validateNameAndCode is the shared rule for every occupation/industry
// classification level: both fields are required.
func validateNameAndCode(name, code string) error {
	if name == "" {
		return ValidationErrorf("name is required")
	}
	if code == "" {
		return ValidationErrorf("code is required")
	}
	return nil
}
