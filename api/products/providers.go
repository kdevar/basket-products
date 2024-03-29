package products

import (
	"github.com/kdevar/basket-products/errors"
	"github.com/olivere/elastic"
)

type ProductService interface {
	GetLiveProductPrices(filter LivePriceFilter) ([]Product, *errors.ApiError)
	GetEstimatedProductPrices(filter EstimatedPriceFilter) (EstimatedPriceResponse, *errors.ApiError)
	SearchProducts(filter SearchFilter) ([]Product, *errors.ApiError)
}

func NewProductService(e *elastic.Client) ProductService {
	return &productServiceImpl{
		elasticClient: e,
	}
}

func NewProductController(p ProductService) *ProductsController {
	return &ProductsController{
		ProductService: p,
	}
}
