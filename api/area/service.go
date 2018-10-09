package area

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kdevar/basket-products/api/stores"
	"github.com/kdevar/basket-products/config"
	"github.com/kdevar/basket-products/const"
	"github.com/kdevar/basket-products/errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

type areaServiceImpl struct {
	Config       *config.Config
	storeService stores.StoreService
}

//TODO fix the error handling here
func (svc *areaServiceImpl) GetAreaInformation(filter AreaFilter) (*Area, *errors.ApiError) {
	const CONTENTTYPE string = "application/json"

	jsonData := []map[string]float64{{
		_const.LATITUDEFIELD:  filter.point.Lat,
		_const.LONGITUDEFIELD: filter.point.Lon,
	}}
	jsonValue, _ := json.Marshal(jsonData)

	response, err := http.Post(
		svc.Config.BasketBaseApiPath+svc.Config.AreaContextPath,
		CONTENTTYPE,
		bytes.NewBuffer(jsonValue))

	if err != nil {
		return nil, errors.ServerError(err)
	}

	s, _ := svc.storeService.GetStoresForLocation(filter.point)
	chains := make(map[string]stores.Chain)

	for _, store := range s {
		chains[strconv.Itoa(store.ChainID)] = store.Chain
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

	return &area, nil
}
