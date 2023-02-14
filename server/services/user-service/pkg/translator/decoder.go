package translator

import (
	"context"
	"errors"

	"github.com/google/uuid"
	commonProto "github.com/vantoan19/Petifies/proto/common"
	"github.com/vantoan19/Petifies/server/services/user-service/pkg/models"
)

func DecodeCreateUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*commonProto.CreateUserRequest)
	if !ok {
		return nil, errors.New("must be proto request")
	}

	return &models.CreateUserReq{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
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
		Email:     resp.Email,
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
		CreatedAt: resp.CreatedAt.AsTime(),
		UpdatedAt: resp.UpdatedAt.AsTime(),
	}, nil
}

func DecodeLoginRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*commonProto.LoginRequest)
	if !ok {
		return nil, errors.New("must be proto request")
	}

	return &models.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	}, nil
}

func DecodeLoginResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*commonProto.LoginResponse)
	if !ok {
		return nil, errors.New("must be proto response")
	}

	return &models.LoginResp{
		AccessToken: resp.AccessToken,
	}, nil
}
