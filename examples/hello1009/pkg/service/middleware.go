package service

import (
	"context"

	log "github.com/go-kit/kit/log"
)

// Middleware describes a service middleware.
type Middleware func(Hello1009Service) Hello1009Service

type loggingMiddleware struct {
	logger log.Logger
	next   Hello1009Service
}

// LoggingMiddleware takes a logger as a dependency
// and returns a Hello1009Service Middleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Hello1009Service) Hello1009Service {
		return &loggingMiddleware{logger, next}
	}

}

func (l loggingMiddleware) Foo(ctx context.Context, s string) (rs string, err error) {
	defer func() {
		l.logger.Log("method", "Foo", "s", s, "rs", rs, "err", err)
	}()
	return l.next.Foo(ctx, s)
}
func (l loggingMiddleware) UpdateUserInfo(c0 context.Context, p1 *pb.UpdateUserInfoRequest) (p0 *pb.UpdateUserInfoReply, e1 error) {
	defer func() {
		l.logger.Log("method", "UpdateUserInfo", "p1", p1, "p0", p0, "e1", e1)
	}()
	return l.next.UpdateUserInfo(c0, p1)
}
