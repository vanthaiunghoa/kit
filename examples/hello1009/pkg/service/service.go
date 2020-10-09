package service

import (
	"context"
	"hello1009/pb/gen-go/pb"
)

// Hello1009Service describes the service.
type Hello1009Service interface {
	// Add your methods here
	Foo(ctx context.Context, s string) (rs string, err error)
	UpdateUserInfo(context.Context, *pb.UpdateUserInfoRequest) (*pb.UpdateUserInfoResponse, error)
}

type basicHello1009Service struct{}

// NewBasicHello1009Service returns a naive, stateless implementation of Hello1009Service.
func NewBasicHello1009Service() Hello1009Service {
	return &basicHello1009Service{}
}

// New returns a Hello1009Service with all of the expected middleware wired in.
func New(middleware []Middleware) Hello1009Service {
	var svc Hello1009Service = NewBasicHello1009Service()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}

func (b *basicHello1009Service) Foo(ctx context.Context, s string) (rs string, err error) {
	// TODO implement the business logic of Foo
	return rs, err
}
func (b *basicHello1009Service) UpdateUserInfo(c0 context.Context, req *pb.UpdateUserInfoRequest) (rsp *pb.UpdateUserInfoResponse, e1 error) {
	// TODO implement the business logic of UpdateUserInfo
	return rsp, e1
}
