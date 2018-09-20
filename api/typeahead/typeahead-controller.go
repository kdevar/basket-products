package typeahead

import (
	"github.com/gin-gonic/gin"
	"github.com/kdevar/basket-products/api/errors"
	"net/http"
	"encoding/json"
	"io/ioutil"
)
var Controller *typeaheadController

type Suggestions struct {
	Content struct {
		Suggests []struct {
			Category interface{} `json:"category"`
			ID       int         `json:"id"`
			Name     string      `json:"name"`
			Type     string      `json:"type"`
		} `json:"suggests"`
	} `json:"content"`
	ErrorCode interface{} `json:"errorCode"`
	Message   interface{} `json:"message"`
}

func init(){
	Controller = &typeaheadController{}
}

type typeaheadController struct{}

func (ctrl *typeaheadController) GetSuggestedProducts(c *gin.Context) *errors.ApiError {
	keyword, _ := c.GetQuery("query")

	if len(keyword) < 3{
		c.JSON(http.StatusOK, nil)
		return nil
	}

	latitude := c.GetHeader("latitude")
	longitude :=c.GetHeader("longitude")
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.basketsavings.com/search2/search/suggested2", nil)

	if err != nil {
		return errors.ServerError(err)
	}

	req.Header.Add("Authorization", "c9764fd946b54b9195adc4b7a2ca58cd")
	q := req.URL.Query()
	q.Add("query", keyword)
	q.Add("latitude", latitude )
	q.Add("longitude", longitude )
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return errors.ServerError(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var s Suggestions
	json.Unmarshal(body, &s)

	c.JSON(http.StatusOK, s.Content.Suggests)

	return nil
}
