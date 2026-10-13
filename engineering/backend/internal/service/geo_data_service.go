package service

import (
	"marketlens-go-backend/internal/models"
	"marketlens-go-backend/internal/repositories"
)

func NewGeoDataService(store *repositories.Store) *CRUDService[models.GeoData] {
	return NewCRUDService(
		store.CreateGeoData,
		store.GetAllGeoData,
		store.GetGeoDataByID,
		store.UpdateGeoData,
		store.DeleteGeoData,
		requireNonEmptyField("province", func(g *models.GeoData) string { return g.Province }),
	)
}
