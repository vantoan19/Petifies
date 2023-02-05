package v1

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"

	userProtoV1 "github.com/vantoan19/Petifies/proto/user-service/v1"
	userEndpointsV1 "github.com/vantoan19/Petifies/server/services/user-services/user-service/internal/endpoints/grpc/v1"
)

type gRPCUserServer struct {
	sayHello grpctransport.Handler
}

func NewGRPCUserServer(endpoints userEndpointsV1.UserEndpoints) userProtoV1.UserServer {
	return &gRPCUserServer{
		sayHello: grpctransport.NewServer(
			endpoints.SayHello,
			decodeSayHelloRequest,
			encodeSayHelloResponse,
		),
	}
}

func (s *gRPCUserServer) SayHello(ctx context.Context, req *userProtoV1.HelloWorldRequest) (*userProtoV1.HelloWorldResponse, error) {
	_, resp, err := s.sayHello.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*userProtoV1.HelloWorldResponse), nil
}

func decodeSayHelloRequest(_ context.Context, request interface{}) (interface{}, error) {
	return userEndpointsV1.SayHelloReq{}, nil
}

func encodeSayHelloResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(userEndpointsV1.SayHelloResp)
	return &userProtoV1.HelloWorldResponse{Greeting: resp.Greeting}, nil
}
