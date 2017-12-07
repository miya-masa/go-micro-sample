package products

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

type logmw struct {
	next   ProductService
	server string
	logger log.Logger
}

type LoggingMiddleware func(s ProductService) ProductService

func (l *logmw) ProductByName(ctx context.Context, name string) (res *Product, err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"took", time.Since(begin),
			"res", res,
			"err", err,
		)
	}(time.Now())
	l.logger.Log("server", l.server)
	return l.next.ProductByName(ctx, name)
}

func Logging(logger log.Logger, server string) LoggingMiddleware {

	return func(s ProductService) ProductService {
		return &logmw{
			next:   s,
			server: server,
			logger: logger,
		}
	}
}
