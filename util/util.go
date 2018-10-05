package util

import (
	"github.com/olivere/elastic"
	"strconv"
)

func ConvertSToInterface(sarr []string) []interface{} {
	s := make([]interface{}, len(sarr))
	for i, v := range sarr {
		s[i] = v
	}
	return s
}

func ConvertIToInterface(sarr []int) []interface{} {
	s := make([]interface{}, len(sarr))
	for i, v := range sarr {
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
