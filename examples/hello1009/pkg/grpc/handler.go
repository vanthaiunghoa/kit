package grpc

import (
	"context"
	"errors"
	grpc "github.com/go-kit/kit/transport/grpc"
	context1 "golang.org/x/net/context"
	pb "hello1009/pb/gen-go/pb"
	endpoint "hello1009/pkg/endpoint"
)

// makeFooHandler creates the handler logic
func makeFooHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.FooEndpoint, decodeFooRequest, encodeFooResponse, options...)
}

// decodeFooResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain Foo request.
// TODO implement the decoder
func decodeFooRequest(_ context.Context, req interface{}) (interface{}, error) {
	return nil, errors.New("'decodeFooRequest' is not impelemented")
}

// encodeFooResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
// TODO implement the encoder
func encodeFooResponse(_ context.Context, rsp interface{}) (interface{}, error) {
	return nil, errors.New("'encodeFooResponse' is not impelemented")
}
func (g *grpcServer) Foo(ctx context1.Context, req *pb.FooRequest) (*pb.FooResponse, error) {
	_, rep, err := g.foo.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.FooResponse), nil
}

// makeUpdateUserInfoHandler creates the handler logic
func makeUpdateUserInfoHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.UpdateUserInfoEndpoint, decodeUpdateUserInfoRequest, encodeUpdateUserInfoResponse, options...)
}

// decodeUpdateUserInfoResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain UpdateUserInfo request.
// TODO implement the decoder
func decodeUpdateUserInfoRequest(_ context.Context, req interface{}) (interface{}, error) {
	return nil, errors.New("'decodeUpdateUserInfoRequest' is not impelemented")
}

// encodeUpdateUserInfoResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
// TODO implement the encoder
func encodeUpdateUserInfoResponse(_ context.Context, rsp interface{}) (interface{}, error) {
	return nil, errors.New("'encodeUpdateUserInfoResponse' is not impelemented")
}
func (g *grpcServer) UpdateUserInfo(ctx context1.Context, req *pb.UpdateUserInfoRequest) (*pb.UpdateUserInfoResponse, error) {
	_, rep, err := g.updateUserInfo.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UpdateUserInfoResponse), nil
}
