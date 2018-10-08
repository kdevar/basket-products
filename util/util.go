package util

import (
	"github.com/gin-gonic/gin"
	"github.com/kdevar/basket-products/const"
	"github.com/olivere/elastic"
	"strconv"
	"strings"
)

func ConvertStringToInterface(sArr []string) []interface{} {
	s := make([]interface{}, len(sArr))
	for i, v := range sArr {
		s[i] = v
	}
	return s
}

func ConvertIntToInterface(iArr []int) []interface{} {
	s := make([]interface{}, len(iArr))
	for i, v := range iArr {
		s[i] = v
	}
	return s
}

func ConvertGeoPointToString(point *elastic.GeoPoint) string {
	return strconv.FormatFloat(
		point.Lat, 'f', 6, 64) +
		"," +
		strconv.FormatFloat(point.Lon, 'f', 6, 64)
}

func ConvertHeadersToGeoPoint(c *gin.Context) (*elastic.GeoPoint, bool) {
	latitude, laterr := strconv.ParseFloat(c.GetHeader(_const.LATITUDEFIELD), 64)
	longitude, longerr := strconv.ParseFloat(c.GetHeader(_const.LONGITUDEFIELD), 64)

	if laterr == nil && longerr == nil {
		return &elastic.GeoPoint{Lat: latitude, Lon: longitude}, true
	}

	return nil, false
}

func ConvertStringToIntArr(s string) ([]int, bool) {
	values := strings.Split(s, ",")
	integers := []int{}
	valid := true

	for _, id := range values {
		integer, err := strconv.Atoi(id)

		if err != nil {
			integers = nil
			valid = false
			break
		}
		integers = append(integers, integer)
	}

	return integers, valid

}
