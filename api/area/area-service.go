package area

import (
	"github.com/olivere/elastic"
	"encoding/json"
	"net/http"
	"bytes"
	"fmt"
	"io/ioutil"
)

var Service *areaService

func init(){
	Service = &areaService{}
}

type areaService struct {}

func (svc *areaService) GetAreaInformation(point elastic.GeoPoint) Area {
	jsonData := []map[string]float64{{"latitude": point.Lat, "longitude": point.Lon}}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post("https://api.basketsavings.com/lookup/internal/area/locate", "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Println("couldn't call lookup api")
	}

	body, readErr := ioutil.ReadAll(response.Body)

	if readErr != nil {
		fmt.Println("couldn't parse body")
	}

	var r AreaResponse
	json.Unmarshal(body, &r)

	return r.Content[0]
}
