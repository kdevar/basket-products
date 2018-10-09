package stores

import (
	"github.com/kdevar/basket-products/errors"
	"github.com/olivere/elastic"
)

type StoreService interface {
	GetStoresForLocation(point *elastic.GeoPoint) ([]Store, *errors.ApiError)
	GetChainsForLocation(point *elastic.GeoPoint) ([]Chain, *errors.ApiError)
}

func NewStoreService(e *elastic.Client) StoreService {
	return &storesServiceImpl{
		elasticClient: e,
	}
}
