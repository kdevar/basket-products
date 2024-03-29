package typeahead

import (
	"github.com/kdevar/basket-products/config"
	"github.com/kdevar/basket-products/errors"
)

type TypeadheadService interface {
	GetSuggestedProducts(filter Filter) ([]Suggestions, *errors.ApiError)
}

func NewTypeaheadService(c *config.Config) TypeadheadService {
	return &TypeaheadServiceImpl{
		Config: c,
	}
}

func NewTypeaheadController(s TypeadheadService) *TypeaheadController {
	return &TypeaheadController{
		Service: s,
	}
}
