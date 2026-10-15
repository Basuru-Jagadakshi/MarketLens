package service

import (
	"context"
	"errors"
	"time"
)

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
