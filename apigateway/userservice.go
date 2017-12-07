package apigateway

import (
	"context"
	"strings"
	"time"

	"golang.org/x/time/rate"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	"github.com/miya-masa/go-micro/users"
	"github.com/sony/gobreaker"
)

type ProxyUserService struct {
	userByName endpoint.Endpoint
}

func (m *ProxyUserService) UserByName(ctx context.Context, name string) (*users.User, error) {
	response, err := m.userByName(ctx, name)
	if err != nil {
		return nil, err
	}
	resp := response.(*users.User)
	return resp, nil
}

func NewUserService(proxyURLs []string) users.UserService {

	var (
		qps         = 100
		maxAttempts = 3
		maxTime     = 1 * time.Second
	)

	var (
		subscriber sd.FixedEndpointer
	)
	for _, url := range proxyURLs {
		if !strings.HasPrefix(url, "http") {
			url = "http://" + url
		}
		var e endpoint.Endpoint
		e = NewEndpoints(url).ProductByName
		e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
		e = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), qps))(e)
		subscriber = append(subscriber, e)
	}

	balancer := lb.NewRoundRobin(subscriber)
	retry := lb.Retry(maxAttempts, maxTime, balancer)

	return &ProxyUserService{retry}
}
