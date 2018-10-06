package api

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/kdevar/basket-products/api/area"
	"github.com/kdevar/basket-products/api/products"
	"github.com/kdevar/basket-products/api/typeahead"
	"github.com/kdevar/basket-products/config"
	"github.com/kdevar/basket-products/errors"
	"go.uber.org/dig"
)

type AppHandlerFunc func(*gin.Context) *errors.ApiError

func withErrorHandling(fn AppHandlerFunc) gin.HandlerFunc {
	return errors.WithError(fn).Handle
}

type ServerParams struct {
	dig.In
	Config    *config.Config
	Area      *area.AreaController
	Products  *products.ProductsController
	Typeahead *typeahead.TypeaheadController
}

type Server struct {
	ServerParams
}

func (s *Server) Run() {
	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./views", true)))

	api := router.Group("/api")

	area := api.Group("/area")

	basketProducts := api.Group("/basket-products")

	typeahead := api.Group("/typeahead")

	area.GET("/",
		withErrorHandling(s.Area.GetAreaInformation))

	typeahead.GET("/products",
		withErrorHandling(s.Typeahead.GetSuggestedProducts))

	basketProducts.GET("/",
		withErrorHandling(s.Products.SearchProducts))

	basketProducts.GET("/:productid/estimated-prices",
		withErrorHandling(s.Products.GetEstimatedProductPrices))

	basketProducts.GET("/:productid/prices",
		withErrorHandling(s.Products.GetLiveProductPrices))

	router.Run(":8080")
}

func NewServer(params ServerParams) *Server {
	return &Server{params}
}
