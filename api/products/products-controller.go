package products

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/kdevar/basket-products/api/errors"
	"net/http"
	"strconv"
	"github.com/olivere/elastic"
)

var Controller *productsController
var ctx = context.Background()

func init() {
	Controller = &productsController{}
}

type productsController struct{}

func (ctrl *productsController) GetProduct(c *gin.Context) *errors.ApiError{
	latitude, _ := strconv.ParseFloat(c.GetHeader("latitude"), 64)
	longitude, _ := strconv.ParseFloat(c.GetHeader("longitude"), 64)
	productId := c.Params.ByName("productid")
	chainId, _ := c.GetQuery("chainId")

	request := ProductEstimateRequest{
		Location: elastic.GeoPoint{Lat:latitude, Lon:longitude},
		ProductId:productId,
		ChainId:chainId,

	}

	c.JSON(http.StatusOK, Service.GetProductEstimate(request))
	return nil
}

func (ctrl *productsController) GetProducts(c *gin.Context) *errors.ApiError {

	latitude, laterr := strconv.ParseFloat(c.GetHeader("latitude"), 64)
	longitude, longerr := strconv.ParseFloat(c.GetHeader("longitude"), 64)
	keyword, _ := c.GetQuery("keyword")
	categoryId, cok := c.GetQuery("categoryId")

	productIds,ok := c.GetQueryArray("productIds")

	filter := ProductFilter{
		keyword: keyword,
	}

	if laterr == nil && longerr == nil {
		filter.location = &elastic.GeoPoint{Lat:latitude, Lon:longitude}
	}

	if cok {
		filter.categoryId = &categoryId
	}

	if ok  {
		for _,id := range productIds {
			i, _ := strconv.Atoi(id)
			filter.productIds = append(filter.productIds, i )
		}

	}

	products, err  := Service.GetProducts(filter)

	if err != nil {
		return errors.ServerError(err)
	}

	c.JSON(http.StatusOK, products)
	return nil
}