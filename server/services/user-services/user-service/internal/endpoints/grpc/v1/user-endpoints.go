package v1

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	userService "github.com/vantoan19/Petifies/server/services/user-services/user-service/internal/service"
)

type UserEndpoints struct {
	SayHello endpoint.Endpoint
}

func MakeAuthenticateEndpoint(s userService.UserService) UserEndpoints {
	return UserEndpoints{
		SayHello: makeSayHelloEndpoint(s),
	}
}

// Definition of endpoints below

type SayHelloReq struct{}

type SayHelloResp struct {
	Greeting string
}

func makeSayHelloEndpoint(s userService.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		result, _ := s.SayHello(ctx)
		return SayHelloResp{Greeting: result}, nil
	}
}
