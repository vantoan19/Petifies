package translator

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	commonProto "github.com/vantoan19/Petifies/proto/common"
	userProtoV1 "github.com/vantoan19/Petifies/proto/user-service/v1"
	"github.com/vantoan19/Petifies/server/services/user-service/pkg/models"
)

func EncodeCreateUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.CreateUserReq)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "must be endpoints' request")
	}

	return &commonProto.CreateUserRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}, nil
}

func EncodeCreateUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.User)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "must be endpoints' response")
	}

	// No password in the response
	return &commonProto.User{
		Id:          resp.ID.String(),
		Email:       resp.Email,
		FirstName:   resp.FirstName,
		LastName:    resp.LastName,
		IsActivated: resp.IsActivated,
		CreatedAt:   timestamppb.New(resp.CreatedAt),
		UpdatedAt:   timestamppb.New(resp.UpdatedAt),
	}, nil
}

func EncodeLoginRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.LoginReq)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "must be endpoints' request")
	}

	return &commonProto.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}, nil
}

func EncodeLoginResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.LoginResp)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "must be endpoints' response")
	}

	return &commonProto.LoginResponse{
		SessionId:             resp.SessionID.String(),
		AccessToken:           resp.AccessToken,
		AccessTokenExpiresAt:  timestamppb.New(resp.AccessTokenExpiresAt),
		RefreshToken:          resp.RefreshToken,
		RefreshTokenExpiresAt: timestamppb.New(resp.RefreshTokenExpiresAt),
		User: &commonProto.User{
			Id:          resp.User.ID.String(),
			Email:       resp.User.Email,
			FirstName:   resp.User.FirstName,
			LastName:    resp.User.LastName,
			IsActivated: resp.User.IsActivated,
			CreatedAt:   timestamppb.New(resp.User.CreatedAt),
			UpdatedAt:   timestamppb.New(resp.User.UpdatedAt),
		},
	}, nil
}

func EncodeVerifyTokenRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.VerifyTokenReq)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "must be endpoints' request")
	}

	return &userProtoV1.VerifyTokenRequest{
		Token: req.Token,
	}, nil
}

func EncodeVerifyTokenResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.VerifyTokenResp)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "must be endpoints' response")
	}

	return &userProtoV1.VerifyTokenResponse{
		UserId: resp.UserID,
	}, nil
}

func EncodeRefreshTokenRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.RefreshTokenReq)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "must be endpoints' request")
	}

	return &commonProto.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	}, nil
}

func EncodeRefreshTokenResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.RefreshTokenResp)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "must be endpoints' response")
	}

	return &commonProto.RefreshTokenResponse{
		AccessToken:          resp.AccessToken,
		AccessTokenExpiresAt: timestamppb.New(resp.AccessTokenExpiresAt),
	}, nil
}

func EncodeGetUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.GetUserReq)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "must be endpoints' request")
	}

	return &userProtoV1.GetUserRequest{
		UserId: req.ID.String(),
	}, nil
}

func EncodeGetUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.User)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "must be endpoints' response")
	}

	return &commonProto.User{
		Id:          resp.ID.String(),
		Email:       resp.Email,
		FirstName:   resp.FirstName,
		LastName:    resp.LastName,
		IsActivated: resp.IsActivated,
		CreatedAt:   timestamppb.New(resp.CreatedAt),
		UpdatedAt:   timestamppb.New(resp.UpdatedAt),
	}, nil
}
