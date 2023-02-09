package v1

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"

	authProtoV1 "github.com/vantoan19/Petifies/proto/auth-service/v1"
	authEndpointsV1 "github.com/vantoan19/Petifies/server/services/user-services/auth-service/internal/endpoints/grpc/v1"
)

type gRPCAuthServer struct {
	sayHello grpctransport.Handler
}

func New(endpoints authEndpointsV1.AuthenticateEndpoints) authProtoV1.AuthServer {
	return &gRPCAuthServer{
		sayHello: grpctransport.NewServer(
			endpoints.SayHello,
			decodeSayHelloRequest,
			encodeMathResponse,
		),
	}
}

func (s *gRPCAuthServer) SayHello(ctx context.Context, req *authProtoV1.HelloWorldRequest) (*authProtoV1.HelloWorldResponse, error) {
	_, resp, err := s.sayHello.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*authProtoV1.HelloWorldResponse), nil
}

func decodeSayHelloRequest(_ context.Context, request interface{}) (interface{}, error) {
	return authEndpointsV1.SayHelloReq{}, nil
}

func encodeMathResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(authEndpointsV1.SayHelloResp)
	return &authProtoV1.HelloWorldResponse{Greeting: resp.Greeting}, nil
}
