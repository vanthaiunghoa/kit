package service

import "context"

// Hello1025Service describes the service.
type Hello1025Service interface {
	// Add your methods here
	Foo(ctx context.Context, s string) (rs string, err error)
}

type basicHello1025Service struct{}

func (b *basicHello1025Service) Foo(ctx context.Context, s string) (rs string, err error) {
	// TODO implement the business logic of Foo
	return rs, err
}

// NewBasicHello1025Service returns a naive, stateless implementation of Hello1025Service.
func NewBasicHello1025Service() Hello1025Service {
	return &basicHello1025Service{}
}

// New returns a Hello1025Service with all of the expected middleware wired in.
func New(middleware []Middleware) Hello1025Service {
	var svc Hello1025Service = NewBasicHello1025Service()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
