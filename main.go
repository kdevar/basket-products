package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kdevar/basket-products/api/errors"
	"github.com/kdevar/basket-products/api/products"
	"github.com/gin-gonic/contrib/static"
	typeahead2 "github.com/kdevar/basket-products/api/typeahead"
	area2 "github.com/kdevar/basket-products/api/area"
)

type AppHandlerFunc func(*gin.Context) *errors.ApiError

func withErrorHandling(fn AppHandlerFunc) gin.HandlerFunc {
	return errors.WithError(fn).Handle
}

func withContextPath(path string) string{
	return "/basket-products" + path
}

func main(){

	router := gin.Default();

	router.Use(static.Serve("/", static.LocalFile("./views", true)))



	api := router.Group("/api")

	area := api.Group("/area")

	typeahead := api.Group("/typeahead")

	area.GET("/", withErrorHandling(area2.Controller.GetAreaInformation))

	typeahead.GET("/products", withErrorHandling(typeahead2.Controller.GetSuggestedProducts))

	api.GET(
		withContextPath("/"),
		withErrorHandling(products.Controller.GetProducts))

	api.GET(
		withContextPath("/:productid"),
		withErrorHandling(products.Controller.GetProduct))

	api.GET(
		withContextPath("/:productid/estimated-prices"),
		withErrorHandling(products.Controller.GetProductEstimates))

	api.GET(withContextPath("/:productid/prices"),
		withErrorHandling(products.Controller.GetProductPrices))


	router.Run(":8080")
}
