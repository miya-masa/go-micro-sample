package products

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	ProductByName endpoint.Endpoint
}

func NewEndpoints(srv ProductService) *Endpoints {
	return &Endpoints{
		ProductByName: makeProductByNameEndpoint(srv),
	}
}

func makeProductByNameEndpoint(svc ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		product, err := svc.ProductByName(ctx, req)
		if err != nil {
			return nil, err
		}
		return product, nil
	}
}
