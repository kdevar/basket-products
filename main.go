package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kdevar/basket-products/api/errors"
	"github.com/kdevar/basket-products/api/products"
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

	router.GET(
		withContextPath("/"),
		withErrorHandling(products.Controller.GetProducts))

	router.GET(
		withContextPath("/:productid"),
		withErrorHandling(products.Controller.GetProduct))


	router.Run(":8080")
}
