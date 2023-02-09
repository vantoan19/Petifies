package v1

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"

	userService "github.com/vantoan19/Petifies/server/services/user-services/user-service/internal/service"
)

type UserEndpoints struct {
	CreateUser endpoint.Endpoint
}

func New(s userService.UserService) UserEndpoints {
	return UserEndpoints{
		CreateUser: makeCreateUserEndpoint(s),
	}
}

// Definition of endpoints below

type CreateUserReq struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

type CreateUserResp struct {
	ID        uuid.UUID
	Email     string
	Password  string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func makeCreateUserEndpoint(s userService.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateUserReq)
		result, err := s.CreateUser(ctx, req.Email, req.Password, req.FirstName, req.LastName)
		if err != nil {
			return CreateUserResp{}, nil
		}
		return CreateUserResp{
			ID:        result.GetID(),
			Email:     result.GetEmail(),
			Password:  result.GetPassword(),
			FirstName: result.GetName().GetFirstName(),
			LastName:  result.GetName().GetLastName(),
			CreatedAt: result.GetCreatedAt(),
			UpdatedAt: result.GetUpdatedAt(),
		}, nil
	}
}
