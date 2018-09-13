package area


type Area struct {
	Location struct {
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
	} `json:"location"`
	PostalCode    string `json:"postalCode"`
	PostalCodeID  int    `json:"postalCodeId"`
	CityID        int    `json:"cityId"`
	CityName      string `json:"cityName"`
	StateName     string `json:"stateName"`
	StateID       int    `json:"stateId"`
	MetroAreaID   int    `json:"metroAreaId"`
	MetroAreaName string `json:"metroAreaName"`
	CountryID     int    `json:"countryId"`
	CountryName   string `json:"countryName"`
	TimeZone      string `json:"timeZone"`
}


type AreaResponse struct {
	Content   []Area      `json:"content"`
	Message   interface{} `json:"message"`
	ErrorCode interface{} `json:"errorCode"`
}