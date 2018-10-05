package products

import (
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
	"strconv"
)

type LivePriceFilter struct {
	productId string
	location  *elastic.GeoPoint
	storeIds  []string
}

func (filter *LivePriceFilter) transform(c *gin.Context) {
	latitude, laterr := strconv.ParseFloat(c.GetHeader("latitude"), 64)
	longitude, longerr := strconv.ParseFloat(c.GetHeader("longitude"), 64)

	productId := c.Params.ByName("productid")
	storeIds, _ := c.GetQueryArray("storeId")

	filter.productId = productId
	filter.storeIds = storeIds

	if laterr == nil && longerr == nil {
		filter.location = &elastic.GeoPoint{Lat: latitude, Lon: longitude}
	}
}
