package typeahead

import (
	"github.com/kdevar/basket-products/config"
	"net/http"
	"github.com/kdevar/basket-products/errors"
	"io/ioutil"
	"encoding/json"
)

type TypeaheadServiceImpl struct {
	Config *config.Config
}

func (svc *TypeaheadServiceImpl) GetSuggestedProducts(filter Filter) ([]Suggestions,*errors.ApiError){
	client := &http.Client{}
	req, err := http.NewRequest("GET", svc.Config.BasketBaseApiPath + svc.Config.TypeAheadContextPath, nil)

	if err != nil {
		return nil, errors.ServerError(err)
	}

	req.Header.Add("Authorization", svc.Config.TypeAheadToken)
	q := req.URL.Query()
	q.Add("query", filter.keyword)
	q.Add("latitude", filter.latitude)
	q.Add("longitude", filter.longitude)
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil,errors.ServerError(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var s SuggestionResponse
	json.Unmarshal(body, &s)

	return s.Content.Suggests, nil
}