package products

import (
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
	"strconv"
	"strings"
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
	latitude, laterr := strconv.ParseFloat(c.GetHeader("latitude"), 64)
	longitude, longerr := strconv.ParseFloat(c.GetHeader("longitude"), 64)
	keyword, kok := c.GetQuery("keyword")
	categoryId, cok := c.GetQuery("categoryId")
	brandId, bok := c.GetQuery("brandId")
	typeId, tok := c.GetQuery("typeId")

	potentialProductIds := strings.Split(keyword, ",")
	parsedIds := []int{}
	allvalidids := true
	if kok && len(potentialProductIds) > 0 {
		for _, id := range potentialProductIds {
			parsedId, err := strconv.Atoi(id)

			if err == nil {
				parsedIds = append(parsedIds, parsedId)
			} else {
				allvalidids = false
				break
			}
		}
	}

	if laterr == nil && longerr == nil {
		f.location = &elastic.GeoPoint{Lat: latitude, Lon: longitude}
	}

	if cok {
		f.categoryId = &categoryId
	}

	if bok {
		f.brandId = &brandId
	}

	if tok {
		f.typeId = &typeId
	}

	if len(parsedIds) > 0 && allvalidids {
		f.productIds = parsedIds
	} else {
		f.keyword = keyword
	}
}
