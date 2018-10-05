package main

import (
	"github.com/kdevar/basket-products/api"
	"github.com/kdevar/basket-products/api/area"
	"github.com/kdevar/basket-products/api/products"
	"github.com/kdevar/basket-products/api/stores"
	"github.com/kdevar/basket-products/api/typeahead"
	"github.com/kdevar/basket-products/config"
)

func main() {
	//manual wiring of the object graph
	cfg := config.NewConfig()
	elasticClient := config.NewElasticClient(cfg)
	productsService := products.NewProductService(elasticClient)
	productController := products.NewProductController(productsService)
	typeaheadController := &typeahead.TypeaheadController{}
	storeService := stores.NewStoreServie(elasticClient)
	areaService := area.NewAreaService(cfg, storeService)
	areaController := area.NewAreaController(areaService)

	server := api.NewServer(api.ServerParams{
		Config:    cfg,
		Products:  productController,
		Typeahead: typeaheadController,
		Area:      areaController,
	})

	server.Run()
}
