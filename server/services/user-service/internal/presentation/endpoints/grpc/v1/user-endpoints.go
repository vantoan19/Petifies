package v1

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	commonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	services "github.com/vantoan19/Petifies/server/services/user-service/internal/application/services"
	userAggre "github.com/vantoan19/Petifies/server/services/user-service/internal/domain/aggregates/user"
	"github.com/vantoan19/Petifies/server/services/user-service/pkg/models"
)

type UserEndpoints struct {
	CreateUser     endpoint.Endpoint
	Login          endpoint.Endpoint
	GetUser        endpoint.Endpoint
	ListUsersByIds endpoint.Endpoint
}

func NewUserEndpoints(s services.UserService) UserEndpoints {
	return UserEndpoints{
		CreateUser:     makeCreateUserEndpoint(s),
		Login:          makeLoginEndpoint(s),
		GetUser:        makeGetUserEndpoint(s),
		ListUsersByIds: makeListUsersByIdEndpoint(s),
	}
}

func makeCreateUserEndpoint(s services.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.CreateUserReq)
		result, err := s.CreateUser(ctx, req.Email, req.Password, req.FirstName, req.LastName)
		if err != nil {
			return nil, err
		}
		return &models.User{
			ID:          result.GetID(),
			Email:       result.GetEmail(),
			FirstName:   result.GetName().GetFirstName(),
			LastName:    result.GetName().GetLastName(),
			IsActivated: result.GetIsActivated(),
			CreatedAt:   result.GetCreatedAt(),
			UpdatedAt:   result.GetUpdatedAt(),
		}, nil
	}
}

func makeLoginEndpoint(s services.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.LoginReq)
		sessionID, refreshToken, refreshExpiresAt, accessToken, accessExpiresAt, user, err := s.Login(ctx, req.Email, req.Password)
		if err != nil {
			return nil, err
		}
		return &models.LoginResp{
			SessionID:             sessionID,
			AccessToken:           accessToken,
			AccessTokenExpiresAt:  accessExpiresAt,
			RefreshToken:          refreshToken,
			RefreshTokenExpiresAt: refreshExpiresAt,
			User: models.User{
				ID:          user.GetID(),
				Email:       user.GetEmail(),
				FirstName:   user.GetName().GetFirstName(),
				LastName:    user.GetName().GetLastName(),
				IsActivated: user.GetIsActivated(),
				CreatedAt:   user.GetCreatedAt(),
				UpdatedAt:   user.GetUpdatedAt(),
			},
		}, nil
	}
}

func makeGetUserEndpoint(s services.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.GetUserReq)
		user, err := s.GetUser(ctx, req.ID)
		if err != nil {
			return nil, err
		}
		return &models.User{
			ID:          user.GetID(),
			Email:       user.GetEmail(),
			FirstName:   user.GetName().GetFirstName(),
			LastName:    user.GetName().GetLastName(),
			IsActivated: user.GetIsActivated(),
			CreatedAt:   user.GetCreatedAt(),
			UpdatedAt:   user.GetUpdatedAt(),
		}, nil
	}
}

func makeListUsersByIdEndpoint(s services.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.ListUsersByIdsReq)
		users, err := s.ListUserByIds(ctx, req.Ids)
		if err != nil {
			return nil, err
		}
		return &models.ListUsersByIdsResp{
			Users: commonutils.Map2(users, func(u *userAggre.User) *models.User {
				return &models.User{
					ID:          u.GetID(),
					Email:       u.GetEmail(),
					FirstName:   u.GetName().GetFirstName(),
					LastName:    u.GetName().GetLastName(),
					IsActivated: u.GetIsActivated(),
					CreatedAt:   u.GetCreatedAt(),
					UpdatedAt:   u.GetUpdatedAt(),
				}
			}),
		}, nil
	}
}
