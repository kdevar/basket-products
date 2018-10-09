package products

import (
	"github.com/gin-gonic/gin"
	"github.com/kdevar/basket-products/const"
	"github.com/kdevar/basket-products/util"
	"github.com/olivere/elastic"
	"strings"
)

type EstimatedPriceFilter struct {
	ProductId   string
	ChainId     []string
	MetroAreaId string
	CityId      string
	ZipCodeId   string
	Location    *elastic.GeoPoint
}

func (filter *EstimatedPriceFilter) GetLatLongString() string {
	return util.ConvertGeoPointToString(filter.Location)
}

func (filter *EstimatedPriceFilter) Transform(c *gin.Context) {

	if cityId, ok := c.GetQuery(_const.CITYIDFIELD); ok {
		filter.CityId = cityId
	}

	if chainIds, ok := c.GetQuery(_const.CHAINIDFIELD); ok {
		filter.ChainId = strings.Split(chainIds, ",")
	}

	if metroAreaId, ok := c.GetQuery(_const.METROAREAIDFIELD); ok {
		filter.MetroAreaId = metroAreaId
	}

	if zipCodeId, ok := c.GetQuery(_const.ZIPIDFIELD); ok {
		filter.ZipCodeId = zipCodeId
	}

	if point, ok := util.ConvertHeadersToGeoPoint(c); ok {
		filter.Location = point
	}

	productId := c.Params.ByName(strings.ToLower(_const.PRODUCTIDFIELD))
	filter.ProductId = productId
}
