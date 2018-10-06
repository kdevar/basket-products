package typeahead

import (
	"github.com/gin-gonic/gin"
	"github.com/kdevar/basket-products/errors"
	"net/http"
)

var Controller *TypeaheadController

type TypeaheadController struct {
	Service TypeadheadService
}

func (ctrl *TypeaheadController) GetSuggestedProducts(c *gin.Context) *errors.ApiError {
	filter := Filter{}
	filter.transform(c)

	if len(filter.keyword) < 3 {
		c.JSON(http.StatusOK, nil)
		return nil
	}
	suggestions, err := ctrl.Service.GetSuggestedProducts(filter)

	if err != nil {
		return errors.ServerError(err)
	}
	c.JSON(http.StatusOK, suggestions)

	return nil
}
