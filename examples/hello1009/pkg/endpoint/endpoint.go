package endpoint

import (
	"context"
	"hello1009/pb/gen-go/pb"
	service "hello1009/pkg/service"

	endpoint "github.com/go-kit/kit/endpoint"
)

// FooRequest collects the request parameters for the Foo method.
type FooRequest struct {
	S string `json:"s"`
}

// FooResponse collects the response parameters for the Foo method.
type FooResponse struct {
	Rs  string `json:"rs"`
	Err error  `json:"err"`
}

// MakeFooEndpoint returns an endpoint that invokes Foo on the service.
func MakeFooEndpoint(s service.Hello1009Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*FooRequest)
		rs, err := s.Foo(ctx, req.S)
		return &FooResponse{
			Err: err,
			Rs:  rs,
		}, nil
	}
}

// Failed implements Failer.
func (r FooResponse) Failed() error {
	return r.Err
}

// UpdateUserInfoRequest collects the request parameters for the UpdateUserInfo method.
type UpdateUserInfoRequest struct {
	Req *pb.UpdateUserInfoRequest `json:"req"`
}

// UpdateUserInfoResponse collects the response parameters for the UpdateUserInfo method.
type UpdateUserInfoResponse struct {
	Rsp *pb.UpdateUserInfoResponse `json:"rsp"`
	E1  error                      `json:"e1"`
}

// MakeUpdateUserInfoEndpoint returns an endpoint that invokes UpdateUserInfo on the service.
func MakeUpdateUserInfoEndpoint(s service.Hello1009Service) endpoint.Endpoint {
	return func(c0 context.Context, request interface{}) (interface{}, error) {
		req := request.(*UpdateUserInfoRequest)
		rsp, e1 := s.UpdateUserInfo(c0, req.Req)
		return &UpdateUserInfoResponse{
			E1:  e1,
			Rsp: rsp,
		}, nil
	}
}

// Failed implements Failer.
func (r UpdateUserInfoResponse) Failed() error {
	return r.E1
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

// Foo implements Service. Primarily useful in a client.
func (e Endpoints) Foo(ctx context.Context, s string) (rs string, err error) {
	request := &FooRequest{S: s}
	response, err := e.FooEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(*FooResponse).Rs, response.(*FooResponse).Err
}

// UpdateUserInfo implements Service. Primarily useful in a client.
func (e Endpoints) UpdateUserInfo(c0 context.Context, req *pb.UpdateUserInfoRequest) (rsp *pb.UpdateUserInfoResponse, e1 error) {
	request := &UpdateUserInfoRequest{Req: req}
	response, err := e.UpdateUserInfoEndpoint(c0, request)
	if err != nil {
		return
	}
	return response.(*UpdateUserInfoResponse).Rsp, response.(*UpdateUserInfoResponse).E1
}
