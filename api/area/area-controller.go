package area

import (
	"github.com/gin-gonic/gin"
	"github.com/kdevar/basket-products/errors"
	"github.com/olivere/elastic"
	"net/http"
	"strconv"
)

type AreaController struct {
	AreaService AreaService
}

func (ctrl *AreaController) GetAreaInformation(c *gin.Context) *errors.ApiError {
	lat, _ := strconv.ParseFloat(c.GetHeader("latitude"), 64)
	lon, _ := strconv.ParseFloat(c.GetHeader("longitude"), 64)
	point := elastic.GeoPoint{lat, lon}
	results, err := ctrl.AreaService.GetAreaInformation(point)
	if err != nil {
		return errors.ServerError(err)
	}
	c.JSON(http.StatusOK, results)
	return nil
}
