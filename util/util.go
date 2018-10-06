package util

import (
	"github.com/olivere/elastic"
	"strconv"
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

func ConvertPointToString(point *elastic.GeoPoint) string {
	return strconv.FormatFloat(
		point.Lat, 'f', 6, 64) +
		"," +
		strconv.FormatFloat(point.Lon, 'f', 6, 64)
}
