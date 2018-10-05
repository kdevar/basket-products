package products

import (
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
	"strconv"
	"github.com/kdevar/basket-products/const"
)

type LivePriceFilter struct {
	productId string
	location  *elastic.GeoPoint
	storeIds  []string
}

func (filter *LivePriceFilter) transform(c *gin.Context) {
	latitude, laterr := strconv.ParseFloat(c.GetHeader(_const.LATITUDEFIELD), 64)
	longitude, longerr := strconv.ParseFloat(c.GetHeader(_const.LONGITUDEFIELD), 64)

	productId := c.Params.ByName(_const.PRODUCTIDFIELD)
	storeIds, _ := c.GetQueryArray(_const.STOREIDFIELD)

	filter.productId = productId
	filter.storeIds = storeIds

	if laterr == nil && longerr == nil {
		filter.location = &elastic.GeoPoint{Lat: latitude, Lon: longitude}
	}
}
