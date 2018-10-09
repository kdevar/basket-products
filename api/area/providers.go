package area

import (
	"github.com/gin-gonic/gin"
	"github.com/kdevar/basket-products/api/stores"
	"github.com/kdevar/basket-products/config"
	"github.com/kdevar/basket-products/errors"
	"github.com/kdevar/basket-products/util"
	"github.com/olivere/elastic"
)

type AreaFilter struct {
	point *elastic.GeoPoint
}

func (f *AreaFilter) transform(c *gin.Context) {

	if point, ok := util.ConvertHeadersToGeoPoint(c); ok {
		f.point = point
	}

}

type AreaService interface {
	GetAreaInformation(filter AreaFilter) (*Area, *errors.ApiError)
}

func NewAreaController(svc AreaService) *AreaController {
	return &AreaController{
		AreaService: svc,
	}
}

func NewAreaService(config *config.Config, s stores.StoreService) AreaService {
	return &areaServiceImpl{
		Config:       config,
		storeService: s,
	}
}
