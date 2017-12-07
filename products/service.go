package products

import (
	"context"
	"errors"
)

var (
	products = map[string]Product{
		"product1": Product{ID: 1, Name: "product1"},
		"product2": Product{ID: 2, Name: "product2"},
		"product3": Product{ID: 3, Name: "product3"},
	}

	ErrNotFound = errors.New("product not found.")
)

type ProductService interface {
	ProductByName(ctx context.Context, name string) (*Product, error)
}

type impl struct {
}

func NewService() ProductService {
	return &impl{}
}

func (u *impl) ProductByName(ctx context.Context, name string) (*Product, error) {

	product, ok := products[name]
	if ok {
		return &product, nil
	}

	return nil, ErrNotFound
}
