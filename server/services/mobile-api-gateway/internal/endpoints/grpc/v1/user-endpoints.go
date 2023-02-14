package v1

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	commonProto "github.com/vantoan19/Petifies/proto/common"
	services "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/services"
)

type UserEndpoints struct {
	CreateUser endpoint.Endpoint
	Login      endpoint.Endpoint
}

func NewUserEndpoints(s services.UserService) UserEndpoints {
	return UserEndpoints{
		CreateUser: makeCreateUserEndpoint(s),
		Login:      makeLoginEndpoint(s),
	}
}

func makeCreateUserEndpoint(s services.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*commonProto.CreateUserRequest)
		resp, err := s.CreateUser(ctx, req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}

func makeLoginEndpoint(s services.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*commonProto.LoginRequest)
		resp, err := s.Login(ctx, req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}
