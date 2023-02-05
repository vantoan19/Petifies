package v1

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	authenService "github.com/vantoan19/Petifies/server/services/user-services/auth-service/internal/service"
)

type AuthenticateEndpoints struct {
	SayHello endpoint.Endpoint
}

func MakeAuthenticateEndpoint(s authenService.AuthenticateService) AuthenticateEndpoints {
	return AuthenticateEndpoints{
		SayHello: makeSayHelloEndpoint(s),
	}
}

// Definition of endpoints below

type SayHelloReq struct{}

type SayHelloResp struct {
	Greeting string
}

func makeSayHelloEndpoint(s authenService.AuthenticateService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		result, _ := s.SayHello(ctx)
		return SayHelloResp{Greeting: result}, nil
	}
}
