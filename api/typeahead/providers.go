package typeahead

import (
	"github.com/kdevar/basket-products/errors"
	"github.com/kdevar/basket-products/config"
)

type TypeadheadService interface {
	GetSuggestedProducts(filter Filter) ([]Suggestions, *errors.ApiError)
}

func NewTypeaheadService(c *config.Config) *TypeaheadServiceImpl{
	return &TypeaheadServiceImpl{
		Config: c,
	}
}

func NewTypeaheadController(s TypeadheadService) *TypeaheadController{
	return &TypeaheadController{
		Service: s,
	}
}