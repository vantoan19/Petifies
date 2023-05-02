package translator

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	commonProto "github.com/vantoan19/Petifies/proto/common"
	userProtoV1 "github.com/vantoan19/Petifies/proto/user-service/v1"
	commonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/services/user-service/pkg/models"
)

var (
	MustBeEndpointReqErr  = status.Error(codes.InvalidArgument, "must be endpoints' request")
	MustBeEndpointRespErr = status.Error(codes.InvalidArgument, "must be endpoints' response")
)

func EncodeCreateUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.CreateUserReq)
	if !ok {
		return nil, MustBeEndpointReqErr
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
		return nil, MustBeEndpointRespErr
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
		return nil, MustBeEndpointReqErr
	}

	return &commonProto.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}, nil
}

func EncodeLoginResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.LoginResp)
	if !ok {
		return nil, MustBeEndpointRespErr
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
		return nil, MustBeEndpointReqErr
	}

	return &userProtoV1.VerifyTokenRequest{
		Token: req.Token,
	}, nil
}

func EncodeVerifyTokenResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.VerifyTokenResp)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &userProtoV1.VerifyTokenResponse{
		UserId: resp.UserID,
	}, nil
}

func EncodeRefreshTokenRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.RefreshTokenReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &commonProto.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	}, nil
}

func EncodeRefreshTokenResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.RefreshTokenResp)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &commonProto.RefreshTokenResponse{
		AccessToken:          resp.AccessToken,
		AccessTokenExpiresAt: timestamppb.New(resp.AccessTokenExpiresAt),
	}, nil
}

func EncodeGetUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.GetUserReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &userProtoV1.GetUserRequest{
		UserId: req.ID.String(),
	}, nil
}

func EncodeGetUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.User)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return encodeUser(resp), nil
}

func EncodeListUsersByIdsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.ListUsersByIdsReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &userProtoV1.ListUsersByIdsRequest{
		UserIds: commonutils.Map2(req.Ids, func(id uuid.UUID) string { return id.String() }),
	}, nil
}

func EncodeListUsersByIdsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.ListUsersByIdsResp)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return commonutils.Map2(resp.Users, func(u *models.User) *commonProto.User { return encodeUser(u) }), nil
}

func encodeUser(user *models.User) *commonProto.User {
	return &commonProto.User{
		Id:          user.ID.String(),
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		IsActivated: user.IsActivated,
		CreatedAt:   timestamppb.New(user.CreatedAt),
		UpdatedAt:   timestamppb.New(user.UpdatedAt),
	}
}
