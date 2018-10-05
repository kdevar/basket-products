package area

import (
	"github.com/kdevar/basket-products/errors"
	"github.com/kdevar/basket-products/api/stores"
	"github.com/kdevar/basket-products/config"
	"github.com/olivere/elastic"
)

type AreaService interface {
	GetAreaInformation(point elastic.GeoPoint) (Area, *errors.ApiError)
}

func NewAreaController(svc AreaService) *AreaController {
	return &AreaController{
		AreaService: svc,
	}
}

func NewAreaService(config *config.Config, s stores.StoreService) *areaServiceImpl {
	return &areaServiceImpl{
		Config:       config,
		storeService: s,
	}
}
