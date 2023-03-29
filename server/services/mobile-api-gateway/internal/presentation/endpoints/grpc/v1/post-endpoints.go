package v1

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	postService "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/application/services/post"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/pkg/models"
)

type PostEndpoints struct {
	CreatePost          endpoint.Endpoint
	CreateComment       endpoint.Endpoint
	EditPost            endpoint.Endpoint
	EditComment         endpoint.Endpoint
	UserToggleLoveReact endpoint.Endpoint
}

func NewPostEndpoints(s postService.PostService) PostEndpoints {
	return PostEndpoints{
		CreatePost:          makeCreatePostEndpoint(s),
		CreateComment:       makeCreateCommentEndpoint(s),
		EditPost:            makeEditPostEndpoint(s),
		EditComment:         makeEditCommentEndpoint(s),
		UserToggleLoveReact: makeUserToggleLoveReactEndpoint(s),
	}
}

func makeCreatePostEndpoint(s postService.PostService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.UserCreatePostReq)
		resp, err := s.CreatePost(ctx, req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}

func makeCreateCommentEndpoint(s postService.PostService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.UserCreateCommentReq)
		resp, err := s.CreateComment(ctx, req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}

func makeEditPostEndpoint(s postService.PostService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.UserEditPostReq)
		resp, err := s.EditPost(ctx, req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}

func makeEditCommentEndpoint(s postService.PostService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.UserEditCommentReq)
		resp, err := s.EditComment(ctx, req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}

func makeUserToggleLoveReactEndpoint(s postService.PostService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.UserToggleLoveReq)
		resp, err := s.ToggleLoveReactPost(ctx, req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}
