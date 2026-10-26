package ingestion

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"marketlens-go-backend/internal/models"
	"marketlens-go-backend/internal/service"
)

// MetadataBuilder fetches the six reference lists (formality, gender,
// vocational education, employment sector, education level, experience)
// and assembles the extraction schema + instruction text.
//
// Reference lists are cached for ttl and only refetched once that
// expires, since this data changes rarely.
type MetadataBuilder struct {
	formalities          *service.CRUDService[models.Formality]
	genders              *service.CRUDService[models.Gender]
	vocationalEducations *service.CRUDService[models.VocationalEducation]
	employmentSectors    *service.CRUDService[models.EmploymentSector]
	educationLevels      *service.CRUDService[models.EducationLevel]
	experiences          *service.CRUDService[models.Experience]
	ttl                  time.Duration

	mu       sync.RWMutex
	cached   *optionSet
	cachedAt time.Time
}

type optionSet struct {
	formalities          string
	genders              string
	vocationalEducations string
	employmentSectors    string
	educationLevels      string
	experiences          string
}

// NewMetadataBuilder creates a schema builder that refreshes its cached
// reference data at most once per ttl.
func NewMetadataBuilder(
	formalities *service.CRUDService[models.Formality],
	genders *service.CRUDService[models.Gender],
	vocationalEducations *service.CRUDService[models.VocationalEducation],
	employmentSectors *service.CRUDService[models.EmploymentSector],
	educationLevels *service.CRUDService[models.EducationLevel],
	experiences *service.CRUDService[models.Experience],
	ttl time.Duration,
) *MetadataBuilder {
	return &MetadataBuilder{
		formalities:          formalities,
		genders:              genders,
		vocationalEducations: vocationalEducations,
		employmentSectors:    employmentSectors,
		educationLevels:      educationLevels,
		experiences:          experiences,
		ttl:                  ttl,
	}
}

// Build returns the static extraction schema and the instruction text
// with the current reference option lists filled in.
func (b *MetadataBuilder) Build(ctx context.Context) (map[string]any, string, error) {
	options, err := b.getOptions(ctx)
	if err != nil {
		return nil, "", err
	}

	instruction := fmt.Sprintf(
		instructionTemplate,
		options.formalities,
		options.genders,
		options.vocationalEducations,
		options.employmentSectors,
		options.educationLevels,
		options.experiences,
	)

	return Schema, instruction, nil
}

func (b *MetadataBuilder) getOptions(ctx context.Context) (*optionSet, error) {
	b.mu.RLock()
	if b.cached != nil && time.Since(b.cachedAt) < b.ttl {
		cached := b.cached
		b.mu.RUnlock()
		return cached, nil
	}
	b.mu.RUnlock()

	return b.refresh(ctx)
}

func (b *MetadataBuilder) refresh(ctx context.Context) (*optionSet, error) {
	formalities, _, err := b.formalities.GetAll(ctx, 0, 0, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: fetching formalities: %v", ErrTransientExtractionFailure, err)
	}
	genders, _, err := b.genders.GetAll(ctx, 0, 0, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: fetching genders: %v", ErrTransientExtractionFailure, err)
	}
	vocationalEducations, _, err := b.vocationalEducations.GetAll(ctx, 0, 0, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: fetching vocational educations: %v", ErrTransientExtractionFailure, err)
	}
	employmentSectors, _, err := b.employmentSectors.GetAll(ctx, 0, 0, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: fetching employment sectors: %v", ErrTransientExtractionFailure, err)
	}
	educationLevels, _, err := b.educationLevels.GetAll(ctx, 0, 0, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: fetching education levels: %v", ErrTransientExtractionFailure, err)
	}
	experiences, _, err := b.experiences.GetAll(ctx, 0, 0, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: fetching experiences: %v", ErrTransientExtractionFailure, err)
	}

	options := &optionSet{
		formalities: formatOptions(formalities, func(f models.Formality) (uint, string) {
			return f.ID, f.FormalityType
		}),
		genders: formatOptions(genders, func(g models.Gender) (uint, string) {
			return g.ID, g.GenderType
		}),
		vocationalEducations: formatOptions(vocationalEducations, func(v models.VocationalEducation) (uint, string) {
			return v.ID, v.Level
		}),
		employmentSectors: formatOptions(employmentSectors, func(e models.EmploymentSector) (uint, string) {
			return e.ID, e.Sector
		}),
		educationLevels: formatOptions(educationLevels, func(e models.EducationLevel) (uint, string) {
			return e.ID, e.Level
		}),
		experiences: formatOptions(experiences, func(e models.Experience) (uint, string) {
			return e.ID, e.Name
		}),
	}

	b.mu.Lock()
	b.cached = options
	b.cachedAt = time.Now()
	b.mu.Unlock()

	return options, nil
}

// formatOptions converts a list of items into "id: label" lines, one per
// line, decoupled from any specific model's field name via extract.
func formatOptions[T any](items []T, extract func(T) (uint, string)) string {
	lines := make([]string, len(items))
	for i, item := range items {
		id, label := extract(item)
		lines[i] = fmt.Sprintf("%d: %s", id, label)
	}
	return strings.Join(lines, "\n")
}