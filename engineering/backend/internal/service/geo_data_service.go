package service

import (
	"context"
	"marketlens-go-backend/internal/models"
	"marketlens-go-backend/internal/repositories"
	"time"
)

func NewGeoDataService(store *repositories.Store) *CRUDService[models.GeoData] {
	return NewCRUDService(
		store.CreateGeoData,
		func(ctx context.Context, _, _ int, _, _ *time.Time) ([]models.GeoData, error) {
			return store.GetAllGeoData(ctx)
		},
		store.GetGeoDataByID,
		store.UpdateGeoData,
		store.DeleteGeoData,
		requireNonEmptyField("province", func(g *models.GeoData) string { return g.Province }),
	)
}
