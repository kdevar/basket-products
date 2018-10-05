package products

import (
	"github.com/gin-gonic/gin"
	"github.com/kdevar/basket-products/errors"
	"net/http"
)

type ProductsController struct {
	ProductService ProductService
}

//GetLiveProductPrices returns live prices (not estimated)
func (ctrl *ProductsController) GetLiveProductPrices(c *gin.Context) *errors.ApiError {
	filter := LivePriceFilter{}
	filter.transform(c)

	products, err := ctrl.ProductService.GetLiveProductPrices(filter)

	if err != nil {
		return errors.ServerError(err)
	}

	c.JSON(http.StatusOK, products)

	return nil
}

//SearchProducts with parameters
func (ctrl *ProductsController) SearchProducts(c *gin.Context) *errors.ApiError {
	filter := SearchFilter{}
	filter.transform(c)

	products, err := ctrl.ProductService.SearchProducts(filter)

	if err != nil {
		return errors.ServerError(err)
	}

	c.JSON(http.StatusOK, products)
	return nil
}

//GetEstimatedProductPrices returns only esitmated prices
func (ctrl *ProductsController) GetEstimatedProductPrices(c *gin.Context) *errors.ApiError {
	filter := EstimatedPriceFilter{}
	filter.Transform(c)

	result, err := ctrl.ProductService.GetEstimatedProductPrices(filter)

	if err != nil {
		return errors.ServerError(err)
	}

	c.JSON(http.StatusOK, result)

	return nil

}
