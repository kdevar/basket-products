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

func (ctrl *productsController) GetProductPrices (c *gin.Context) *errors.ApiError{
	latitude, laterr := strconv.ParseFloat(c.GetHeader("latitude"), 64)
	longitude, longerr := strconv.ParseFloat(c.GetHeader("longitude"), 64)
	productId := c.Params.ByName("productid")
	storeIds, _ := c.GetQueryArray("storeId")

	filter := ProductPriceFilter{productId: productId, storeIds:storeIds}
	if laterr == nil && longerr == nil {
		filter.location = &elastic.GeoPoint{Lat:latitude, Lon:longitude}
	}
	products, err := Service.GetProductPrices(filter)

	if err != nil {
		return errors.ServerError(err)
	}

	c.JSON(http.StatusOK, products)

	return nil
}

func (ctrl *productsController) GetProduct(c *gin.Context) *errors.ApiError{
	latitude, _ := strconv.ParseFloat(c.GetHeader("latitude"), 64)
	longitude, _ := strconv.ParseFloat(c.GetHeader("longitude"), 64)
	productId := c.Params.ByName("productid")

	chainIds, _ := c.GetQueryArray("chainId")
	metroAreaId, _ := c.GetQuery("metroAreaId")
	cityId, _ := c.GetQuery("cityId")
	zipcodeID, _ := c.GetQuery("zipCodeId")

	request := ProductEstimateRequest{
		Location:    &elastic.GeoPoint{Lat:latitude, Lon:longitude},
		ProductId:   productId,
		MetroAreaId: metroAreaId,
		CityId:      cityId,
		ZipCodeId:   zipcodeID,
	}

	if len(chainIds) > 0{
		for _, chainId := range chainIds {
			request.ChainId = append(request.ChainId, chainId)
		}
	}


	c.JSON(http.StatusOK, Service.GetProductEstimate(request))

	return nil
}

func (ctrl *productsController) GetProducts(c *gin.Context) *errors.ApiError {

	latitude, laterr := strconv.ParseFloat(c.GetHeader("latitude"), 64)
	longitude, longerr := strconv.ParseFloat(c.GetHeader("longitude"), 64)
	keyword, kok := c.GetQuery("keyword")
	categoryId, cok := c.GetQuery("categoryId")
	brandId, bok := c.GetQuery("brandId")
	typeId, tok := c.GetQuery("typeId")
	productIds,ok := c.GetQueryArray("productIds")

	filter := ProductFilter{}

	if kok {
		filter.keyword = keyword
	}

	if laterr == nil && longerr == nil {
		filter.location = &elastic.GeoPoint{Lat:latitude, Lon:longitude}
	}

	if cok {
		filter.categoryId = &categoryId
	}

	if bok {
		filter.brandId = &brandId
	}

	if tok {
		filter.typeId = &typeId
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

func (ctrl *productsController) GetProductEstimates(c *gin.Context) *errors.ApiError{
	request := ProductEstimateRequest{}
	cityId, ok := c.GetQuery("cityId")

	if ok {
		request.CityId = cityId
	}

	chainIds, ok := c.GetQueryArray("chainId")

	if ok {
		request.ChainId = chainIds
	}

	metroAreaId, ok := c.GetQuery("metroAreaId")

	if ok {
		request.MetroAreaId = metroAreaId
	}

	zipCodeId, ok := c.GetQuery("zipCodeId")

	if ok {
		request.ZipCodeId = zipCodeId
	}

	latitude, laterr := strconv.ParseFloat(c.GetHeader("latitude"), 64)
	longitude, lonerr := strconv.ParseFloat(c.GetHeader("longitude"), 64)


	if laterr == nil && lonerr == nil {
		request.Location = &elastic.GeoPoint{Lat:latitude, Lon:longitude}
	}

	productId := c.Params.ByName("productid")
	request.ProductId = productId

	result := Service.GetProductEstimate(request)


	c.JSON(http.StatusOK, result)

	return nil


}