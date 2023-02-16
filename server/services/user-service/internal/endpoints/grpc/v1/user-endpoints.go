package v1

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	services "github.com/vantoan19/Petifies/server/services/user-service/internal/services"
	"github.com/vantoan19/Petifies/server/services/user-service/pkg/models"
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
		req := request.(*models.CreateUserReq)
		result, err := s.CreateUser(ctx, req.Email, req.Password, req.FirstName, req.LastName)
		if err != nil {
			return nil, err
		}
		return &models.CreateUserResp{
			ID:        result.GetID(),
			Email:     result.GetEmail(),
			FirstName: result.GetName().GetFirstName(),
			LastName:  result.GetName().GetLastName(),
			CreatedAt: result.GetCreatedAt(),
			UpdatedAt: result.GetUpdatedAt(),
		}, nil
	}
}

func makeLoginEndpoint(s services.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.LoginReq)
		token, err := s.Login(ctx, req.Email, req.Password)
		if err != nil {
			return nil, err
		}
		return &models.LoginResp{
			AccessToken: token,
		}, nil
	}
}
