package products

import (
	"github.com/gin-gonic/gin"
	"github.com/kdevar/basket-products/const"
	"github.com/kdevar/basket-products/util"
	"github.com/olivere/elastic"
	"strings"
)

type LivePriceFilter struct {
	productId string
	location  *elastic.GeoPoint
	storeIds  []string
}

func (filter *LivePriceFilter) transform(c *gin.Context) {
	if point, ok := util.ConvertHeadersToGeoPoint(c); ok {
		filter.location = point
	}

	productId := c.Params.ByName(strings.ToLower(_const.PRODUCTIDFIELD))
	storeIds, _ := c.GetQueryArray(_const.STOREIDFIELD)

	filter.productId = productId
	filter.storeIds = storeIds

}
