package products

import (
	"github.com/gin-gonic/gin"
	"github.com/kdevar/basket-products/const"
	"github.com/kdevar/basket-products/util"
	"github.com/olivere/elastic"
)

type SearchFilter struct {
	keyword      string
	location     *elastic.GeoPoint
	productIds   []int
	categoryId   *string
	brandId      *string
	typeId       *string
	categoryDesc *string
	from         *int
	to           *int
}

func (f *SearchFilter) transform(c *gin.Context) {
	if point, ok := util.ConvertHeadersToGeoPoint(c); ok {
		f.location = point
	}
	if keyword, ok := c.GetQuery("keyword"); ok {
		if productIds, ok := util.ConvertStringToIntArr(keyword); ok {
			f.productIds = productIds
		} else {
			f.keyword = keyword
		}
	}

	if categoryId, ok := c.GetQuery(_const.CATEGORYIDFIELD); ok {
		f.categoryId = &categoryId
	}

	if brandId, ok := c.GetQuery(_const.BRANDIDFIELD); ok {
		f.brandId = &brandId
	}

	if typeId, ok := c.GetQuery(_const.TYPEIDFIELD); ok {
		f.typeId = &typeId
	}

}
