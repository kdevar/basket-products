package stores

import (
	"github.com/kdevar/basket-products/api/errors"
	"github.com/olivere/elastic"
)

type StoreService interface {
	GetStoresForLocation(point elastic.GeoPoint) ([]Store, *errors.ApiError)
	GetChainsForLocation(point elastic.GeoPoint) ([]Chain, *errors.ApiError)
}

func NewStoreServie(e *elastic.Client) *storesServiceImpl {
	return &storesServiceImpl{
		elasticClient: e,
	}
}
