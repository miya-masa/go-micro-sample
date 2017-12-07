package users

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

type logmw struct {
	next   UserService
	server string
	logger log.Logger
}

type LoggingMiddleware func(s UserService) UserService

func (l *logmw) UserByName(ctx context.Context, name string) (res *User, err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"took", time.Since(begin),
			"res", res,
			"err", err,
		)
	}(time.Now())
	l.logger.Log("server", l.server)
	return l.next.UserByName(ctx, name)
}

func Logging(logger log.Logger, server string) LoggingMiddleware {

	return func(s UserService) UserService {
		return &logmw{
			next:   s,
			server: server,
			logger: logger,
		}
	}
}
