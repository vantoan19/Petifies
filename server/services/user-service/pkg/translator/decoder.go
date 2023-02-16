package translator

import (
	"context"
	"errors"

	"github.com/google/uuid"
	commonProto "github.com/vantoan19/Petifies/proto/common"
	userProtoV1 "github.com/vantoan19/Petifies/proto/user-service/v1"
	"github.com/vantoan19/Petifies/server/services/user-service/pkg/models"
)

func DecodeCreateUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*commonProto.CreateUserRequest)
	if !ok {
		return nil, errors.New("must be proto request")
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
		return nil, errors.New("must be proto response")
	}

	id, err := uuid.Parse(resp.Id)
	if err != nil {
		return nil, err
	}

	return &models.CreateUserResp{
		ID:        id,
		Email:     resp.GetEmail(),
		FirstName: resp.GetFirstName(),
		LastName:  resp.GetLastName(),
		CreatedAt: resp.GetCreatedAt().AsTime(),
		UpdatedAt: resp.GetUpdatedAt().AsTime(),
	}, nil
}

func DecodeLoginRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*commonProto.LoginRequest)
	if !ok {
		return nil, errors.New("must be proto request")
	}

	return &models.LoginReq{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}, nil
}

func DecodeLoginResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*commonProto.LoginResponse)
	if !ok {
		return nil, errors.New("must be proto response")
	}

	return &models.LoginResp{
		AccessToken: resp.GetAccessToken(),
	}, nil
}

func DecodeVerifyTokenRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*userProtoV1.VerifyTokenRequest)
	if !ok {
		return nil, errors.New("must be proto request")
	}

	return &models.VerifyTokenReq{
		Token: req.GetToken(),
	}, nil
}

func DecodeVerifyTokenResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*userProtoV1.VerifyTokenResponse)
	if !ok {
		return nil, errors.New("must be proto response")
	}

	return &models.VerifyTokenResp{
		UserID: resp.GetUserId(),
	}, nil
}
