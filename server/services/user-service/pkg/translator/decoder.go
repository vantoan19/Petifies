package translator

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commonProto "github.com/vantoan19/Petifies/proto/common"
	userProtoV1 "github.com/vantoan19/Petifies/proto/user-service/v1"
	"github.com/vantoan19/Petifies/server/services/user-service/pkg/models"
)

var (
	MustBeProtoReqErr  = status.Error(codes.InvalidArgument, "must be proto request")
	MustBeProtoRespErr = status.Error(codes.InvalidArgument, "must be proto response")
)

func DecodeCreateUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*commonProto.CreateUserRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &models.CreateUserReq{
		Email:     req.GetEmail(),
		Password:  req.GetPassword(),
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
	}, nil
}

func DecodeCreateUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*commonProto.User)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	id, err := uuid.Parse(resp.Id)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:          id,
		Email:       resp.GetEmail(),
		FirstName:   resp.GetFirstName(),
		LastName:    resp.GetLastName(),
		IsActivated: resp.GetIsActivated(),
		CreatedAt:   resp.GetCreatedAt().AsTime(),
		UpdatedAt:   resp.GetUpdatedAt().AsTime(),
	}, nil
}

func DecodeLoginRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*commonProto.LoginRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &models.LoginReq{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}, nil
}

func DecodeLoginResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*commonProto.LoginResponse)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	sessionID, err := uuid.Parse(resp.GetSessionId())
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(resp.GetUser().Id)
	if err != nil {
		return nil, err
	}

	return &models.LoginResp{
		SessionID:             sessionID,
		AccessToken:           resp.GetAccessToken(),
		AccessTokenExpiresAt:  resp.GetAccessTokenExpiresAt().AsTime(),
		RefreshToken:          resp.GetRefreshToken(),
		RefreshTokenExpiresAt: resp.GetRefreshTokenExpiresAt().AsTime(),
		User: models.User{
			ID:          userID,
			Email:       resp.GetUser().Email,
			FirstName:   resp.GetUser().FirstName,
			LastName:    resp.GetUser().LastName,
			IsActivated: resp.User.IsActivated,
			CreatedAt:   resp.GetUser().CreatedAt.AsTime(),
			UpdatedAt:   resp.GetUser().UpdatedAt.AsTime(),
		},
	}, nil
}

func DecodeVerifyTokenRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*userProtoV1.VerifyTokenRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &models.VerifyTokenReq{
		Token: req.GetToken(),
	}, nil
}

func DecodeVerifyTokenResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*userProtoV1.VerifyTokenResponse)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &models.VerifyTokenResp{
		UserID: resp.GetUserId(),
	}, nil
}

func DecodeRefreshTokenRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*commonProto.RefreshTokenRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &models.RefreshTokenReq{
		RefreshToken: req.GetRefreshToken(),
	}, nil
}

func DecodeRefreshTokenResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*commonProto.RefreshTokenResponse)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &models.RefreshTokenResp{
		AccessToken:          resp.GetAccessToken(),
		AccessTokenExpiresAt: resp.GetAccessTokenExpiresAt().AsTime(),
	}, nil
}

func DecodeGetUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*userProtoV1.GetUserRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	id, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	return &models.GetUserReq{
		ID: id,
	}, nil
}

func DecodeGetUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*commonProto.User)
	if !ok {
		return nil, MustBeProtoRespErr
	}

	return decodeUser(resp)
}

func DecodeListUsersByIdsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*userProtoV1.ListUsersByIdsRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	var ids []uuid.UUID
	for _, id := range req.UserIds {
		id_, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id_)
	}

	return &models.ListUsersByIdsReq{
		Ids: ids,
	}, nil
}

func DecodeListUsersByIdsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*userProtoV1.ListUsersByIdsResponse)
	if !ok {
		return nil, MustBeProtoRespErr
	}

	var users []*models.User
	for _, u := range resp.Users {
		user, err := decodeUser(u)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return &models.ListUsersByIdsResp{
		Users: users,
	}, nil
}

func decodeUser(user *commonProto.User) (*models.User, error) {
	id, err := uuid.Parse(user.Id)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:          id,
		Email:       user.GetEmail(),
		FirstName:   user.GetFirstName(),
		LastName:    user.GetLastName(),
		IsActivated: user.GetIsActivated(),
		CreatedAt:   user.GetCreatedAt().AsTime(),
		UpdatedAt:   user.GetUpdatedAt().AsTime(),
	}, nil
}
