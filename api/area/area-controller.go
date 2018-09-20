package area

import (
	"github.com/gin-gonic/gin"
	"github.com/kdevar/basket-products/api/errors"
	"github.com/olivere/elastic"
	"strconv"
	"net/http"
)
var Controller *areaController

func init(){
	Controller = &areaController{}
}

type areaController struct {}

func (ctrl *areaController) GetAreaInformation(c *gin.Context) *errors.ApiError{
	lat, _ := strconv.ParseFloat(c.GetHeader("latitude"), 64)
	lon, _ := strconv.ParseFloat(c.GetHeader("longitude"), 64)
	point := elastic.GeoPoint{lat, lon}
	results, err :=Service.GetAreaInformation(point)
	if err != nil {
		return errors.ServerError(err)
	}
	c.JSON(http.StatusOK, results)
	return nil
}