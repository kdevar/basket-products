package products

import (
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
	"strconv"
)

type EstimatedPriceFilter struct {
	ProductId       string
	ChainId         []string
	MetroAreaId     string
	CityId          string
	ZipCodeId       string
	Location        *elastic.GeoPoint
	IncludeEstimate bool
	IncludeDetails  bool
}

func (filter *EstimatedPriceFilter) GetLatLongString() string {
	return strconv.FormatFloat(filter.Location.Lat, 'f', 6, 64) + "," + strconv.FormatFloat(filter.Location.Lon, 'f', 6, 64)
}

func (filter *EstimatedPriceFilter) Transform(c *gin.Context) {
	cityId, ok := c.GetQuery("cityId")

	if ok {
		filter.CityId = cityId
	}

	chainIds, ok := c.GetQueryArray("chainId")

	if ok {
		filter.ChainId = chainIds
	}

	metroAreaId, ok := c.GetQuery("metroAreaId")

	if ok {
		filter.MetroAreaId = metroAreaId
	}

	zipCodeId, ok := c.GetQuery("zipCodeId")

	if ok {
		filter.ZipCodeId = zipCodeId
	}

	latitude, laterr := strconv.ParseFloat(c.GetHeader("latitude"), 64)
	longitude, lonerr := strconv.ParseFloat(c.GetHeader("longitude"), 64)

	if laterr == nil && lonerr == nil {
		filter.Location = &elastic.GeoPoint{Lat: latitude, Lon: longitude}
	}

	productId := c.Params.ByName("productid")
	filter.ProductId = productId
}