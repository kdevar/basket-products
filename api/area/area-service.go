package area

import (
	"github.com/olivere/elastic"
	"encoding/json"
	"net/http"
	"bytes"
	"fmt"
	"io/ioutil"
	"github.com/kdevar/basket-products/api/errors"
	"github.com/kdevar/basket-products/api/stores"
	"strconv"
)

var Service *areaService

func init(){
	Service = &areaService{}
}

type areaService struct {}

func (svc *areaService) GetAreaInformation(point elastic.GeoPoint) (Area, *errors.ApiError) {
	jsonData := []map[string]float64{{"latitude": point.Lat, "longitude": point.Lon}}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post("https://api.basketsavings.com/lookup/internal/area/locate", "application/json", bytes.NewBuffer(jsonValue))


		s, _ := stores.Service.GetStoresForLocation(point)
	chains := make(map[string]stores.Chain)
	for _,store := range s {
		chains[strconv.Itoa(store.ChainID)] = store.Chain
	}



	if err != nil {
		fmt.Println("couldn't call lookup api")
	}

	body, readErr := ioutil.ReadAll(response.Body)

	if readErr != nil {
		fmt.Println("couldn't parse body")
	}

	var r AreaResponse
	json.Unmarshal(body, &r)

	area := r.Content[0]

	area.Stores = s
	area.Chains = chains

	return area, nil
}

func (svc *areaService) GetTotalAreaInformation(point elastic.GeoPoint) (){

}
