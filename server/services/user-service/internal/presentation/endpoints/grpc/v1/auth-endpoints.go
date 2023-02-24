package v1

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	services "github.com/vantoan19/Petifies/server/services/user-service/internal/application/services"
	"github.com/vantoan19/Petifies/server/services/user-service/pkg/models"
)

type AuthEndpoints struct {
	VerifyToken  endpoint.Endpoint
	RefreshToken endpoint.Endpoint
}

func NewAuthEndpoints(s services.UserService) AuthEndpoints {
	return AuthEndpoints{
		VerifyToken:  makeVerifyTokenEndpoint(s),
		RefreshToken: makeRefreshTokenEndpoint(s),
	}
}

func makeVerifyTokenEndpoint(s services.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.VerifyTokenReq)
		result, err := s.VerifyToken(ctx, req.Token)
		if err != nil {
			return nil, err
		}
		return &models.VerifyTokenResp{
			UserID: result,
		}, nil
	}
}

func makeRefreshTokenEndpoint(s services.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.RefreshTokenReq)
		accessToken, expiresAt, err := s.RefreshToken(ctx, req.RefreshToken)
		if err != nil {
			return nil, err
		}
		return &models.RefreshTokenResp{
			AccessToken:          accessToken,
			AccessTokenExpiresAt: expiresAt,
		}, nil
	}
}
