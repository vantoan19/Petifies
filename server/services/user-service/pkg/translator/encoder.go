package translator

import (
	"context"
	"errors"

	"google.golang.org/protobuf/types/known/timestamppb"

	commonProto "github.com/vantoan19/Petifies/proto/common"
	"github.com/vantoan19/Petifies/server/services/user-service/pkg/models"
)

func EncodeCreateUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.CreateUserReq)
	if !ok {
		return nil, errors.New("must be endpoints' request")
	}

	return &commonProto.CreateUserRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}, nil
}

func EncodeCreateUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.CreateUserResp)
	if !ok {
		return nil, errors.New("must be endpoints' response")
	}

	// No password in the response
	return &commonProto.User{
		Id:        resp.ID.String(),
		Email:     resp.Email,
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
		CreatedAt: timestamppb.New(resp.CreatedAt),
		UpdatedAt: timestamppb.New(resp.UpdatedAt),
	}, nil
}

func EncodeLoginRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.LoginReq)
	if !ok {
		return nil, errors.New("must be endpoints' request")
	}

	return &commonProto.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}, nil
}

func EncodeLoginResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.LoginResp)
	if !ok {
		return nil, errors.New("must be endpoints' response")
	}

	return &commonProto.LoginResponse{
		AccessToken: resp.AccessToken,
	}, nil
}
