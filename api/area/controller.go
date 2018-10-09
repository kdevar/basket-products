package area

import (
	"github.com/gin-gonic/gin"
	"github.com/kdevar/basket-products/errors"
	"net/http"
)

type AreaController struct {
	AreaService AreaService
}

func (ctrl *AreaController) GetAreaInformation(c *gin.Context) *errors.ApiError {
	filter := AreaFilter{}
	filter.transform(c)

	results, err := ctrl.AreaService.GetAreaInformation(filter)

	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, results)
	return nil
}
