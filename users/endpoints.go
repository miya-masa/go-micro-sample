package users

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	UserByName endpoint.Endpoint
}

func NewEndpoints(srv UserService) *Endpoints {
	return &Endpoints{
		UserByName: makeUserByNameEndpoint(srv),
	}
}

func makeUserByNameEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		user, err := svc.UserByName(ctx, req)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
}
