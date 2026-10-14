package repositories

import (
	"context"
	"marketlens-go-backend/internal/models"
)

// Geo data CRUD
func (r *Store) CreateGeoData(ctx context.Context, item *models.GeoData) error {
	return newCRUDRepository[models.GeoData](r.db).create(ctx, item)
}

func (r *Store) GetAllGeoData(ctx context.Context) ([]models.GeoData, error) {
	return newCRUDRepository[models.GeoData](r.db).getAll(ctx, 0, 0, nil, nil)
}

func (r *Store) GetGeoDataByID(ctx context.Context, id uint) (models.GeoData, error) {
	return newCRUDRepository[models.GeoData](r.db).getByID(ctx, id)
}

func (r *Store) UpdateGeoData(ctx context.Context, id uint, updates map[string]interface{}) (models.GeoData, error) {
	return newCRUDRepository[models.GeoData](r.db).update(ctx, id, updates)
}

func (r *Store) DeleteGeoData(ctx context.Context, id uint) error {
	return newCRUDRepository[models.GeoData](r.db).delete(ctx, id)
}
