package main

import (
	"fmt"
	"github.com/kdevar/basket-products/api"
	"github.com/kdevar/basket-products/api/area"
	"github.com/kdevar/basket-products/api/products"
	"github.com/kdevar/basket-products/api/stores"
	"github.com/kdevar/basket-products/api/typeahead"
	"github.com/kdevar/basket-products/config"
	"go.uber.org/dig"
	"os"
	"os/signal"
	"syscall"
)

func CreateContainer() *dig.Container {
	container := dig.New()

	container.Provide(config.NewConfig)
	container.Provide(config.NewElasticClient)
	container.Provide(products.NewProductService)
	container.Provide(products.NewProductController)
	container.Provide(typeahead.NewTypeaheadController)
	container.Provide(typeahead.NewTypeaheadService)
	container.Provide(stores.NewStoreService)
	container.Provide(area.NewAreaService)
	container.Provide(area.NewAreaController)
	container.Provide(api.NewServer)

	return container
}

func RunServer(server api.Server) {
	server.Run()
}

func Cleanup() {
	fmt.Println("done")
}

func main() {

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signals
		Cleanup()
		os.Exit(1)
	}()

	ctr := CreateContainer()

	err := ctr.Invoke(RunServer)

	if err != nil {
		panic(err)
	}
}
