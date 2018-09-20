package area

import "github.com/kdevar/basket-products/api/stores"

type Area struct {
	Location struct {
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
	} `json:"location"`
	PostalCode    string                  `json:"postalCode"`
	PostalCodeID  int                     `json:"postalCodeId"`
	CityID        int                     `json:"cityId"`
	CityName      string                  `json:"cityName"`
	StateName     string                  `json:"stateName"`
	StateID       int                     `json:"stateId"`
	MetroAreaID   int                     `json:"metroAreaId"`
	MetroAreaName string                  `json:"metroAreaName"`
	CountryID     int                     `json:"countryId"`
	CountryName   string                  `json:"countryName"`
	TimeZone      string                  `json:"timeZone"`
	Stores        []stores.Store          `json:"stores"`
	Chains        map[string]stores.Chain `json:"chains"`
}

type AreaResponse struct {
	Content   []Area      `json:"content"`
	Message   interface{} `json:"message"`
	ErrorCode interface{} `json:"errorCode"`
}

type TotalArea struct {
	Area   Area
	Stores []stores.Store
	Chains []stores.Chain
}
