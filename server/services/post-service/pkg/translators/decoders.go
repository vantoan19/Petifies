package translators

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commonProto "github.com/vantoan19/Petifies/proto/common"
	postProtoV1 "github.com/vantoan19/Petifies/proto/post-service/v1"
	utils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
)

var (
	MustBeProtoReqErr     = status.Error(codes.InvalidArgument, "must be proto request")
	MustBeProtoRespErr    = status.Error(codes.InvalidArgument, "must be proto response")
	MustBeEndpointReqErr  = status.Error(codes.InvalidArgument, "must be endpoints' request")
	MustBeEndpointRespErr = status.Error(codes.InvalidArgument, "must be endpoints' response")
)

func DecodeCreatePostRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*postProtoV1.CreatePostRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	authorId, err := uuid.Parse(req.GetAuthorId())
	if err != nil {
		return nil, err
	}

	return &models.CreatePostReq{
		AuthorID:    authorId,
		Visibility:  req.Visibility,
		Activity:    req.Actitivty,
		TextContent: req.GetContent(),
		Images:      utils.Map2(req.GetImages(), func(i *commonProto.Image) models.Image { return models.Image{URL: i.Uri, Description: i.Description} }),
		Videos:      utils.Map2(req.GetVideos(), func(v *commonProto.Video) models.Video { return models.Video{URL: v.Uri, Description: v.Description} }),
	}, nil
}

func DecodePostResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*commonProto.Post)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return decodePostProtoModel(resp)
}

func DecodeCreateCommentRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*postProtoV1.CreateCommentRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	postId, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, err
	}
	authorId, err := uuid.Parse(req.AuthorId)
	if err != nil {
		return nil, err
	}
	parentId, err := uuid.Parse(req.ParentId)
	if err != nil {
		return nil, err
	}

	return &models.CreateCommentReq{
		PostID:       postId,
		AuthorID:     authorId,
		ParentID:     parentId,
		IsPostParent: req.IsPostParent,
		Content:      req.Content,
		ImageContent: models.Image{
			URL:         req.Image.Uri,
			Description: req.Image.Description,
		},
		VideoContent: models.Video{
			URL:         req.Video.Uri,
			Description: req.Video.Description,
		},
	}, nil
}

func DecodeCommentResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*commonProto.Comment)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return decodeCommentProtoModel(resp)
}

func DecodeLoveReactRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*postProtoV1.LoveReactRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	targetID, err := uuid.Parse(req.GetTargetId())
	if err != nil {
		return nil, err
	}
	authorID, err := uuid.Parse(req.GetAuthorId())
	if err != nil {
		return nil, err
	}

	return &models.LoveReactReq{
		TargetID:     targetID,
		AuthorID:     authorID,
		IsTargetPost: req.GetIsTargetPost(),
	}, nil
}

func DecodeLoveResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*commonProto.Love)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return decodeLoveProtoModel(resp)
}

func DecodeEditPostRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*postProtoV1.EditPostRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	postID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	return &models.EditPostReq{
		ID:         postID,
		Visibility: req.Visibility,
		Activity:   req.Activity,
		Content:    req.GetContent(),
		Images:     utils.Map2(req.GetImages(), func(i *commonProto.Image) models.Image { return models.Image{URL: i.Uri, Description: i.Description} }),
		Videos:     utils.Map2(req.GetVideos(), func(v *commonProto.Video) models.Video { return models.Video{URL: v.Uri, Description: v.Description} }),
	}, nil
}

func DecodeEditCommentRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*postProtoV1.EditCommentRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	commentID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	return &models.EditCommentReq{
		ID:      commentID,
		Content: req.GetContent(),
		Image: models.Image{
			URL:         req.Image.Uri,
			Description: req.Image.Description,
		},
		Video: models.Video{
			URL:         req.Video.Uri,
			Description: req.Video.Description,
		},
	}, nil
}

func DecodeListCommentsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*postProtoV1.ListCommentsRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	commentIDs := make([]uuid.UUID, 0)
	for _, id_ := range req.CommentIds {
		id, err := uuid.Parse(id_)
		if err != nil {
			return nil, err
		}
		commentIDs = append(commentIDs, id)
	}

	return &models.ListCommentsReq{
		CommentIDs: commentIDs,
	}, nil
}

func DecodeListCommentsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*postProtoV1.ListCommentsResponse)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	comments := make([]*models.Comment, 0)
	for _, c_ := range resp.GetComments() {
		c, err := decodeCommentProtoModel(c_)
		if err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}

	return &models.ListCommentsResp{
		Comments: comments,
	}, nil
}

func DecodeListPostsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*postProtoV1.ListPostsRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	postIDs := make([]uuid.UUID, 0)
	for _, id_ := range req.PostIds {
		id, err := uuid.Parse(id_)
		if err != nil {
			return nil, err
		}
		postIDs = append(postIDs, id)
	}

	return &models.ListPostsReq{
		PostIDs: postIDs,
	}, nil
}

func DecodeListPostsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*postProtoV1.ListPostsResponse)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	posts := make([]*models.Post, 0)
	for _, p_ := range resp.GetPosts() {
		p, err := decodePostProtoModel(p_)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	return &models.ListPostsResp{
		Posts: posts,
	}, nil
}

func DecodeGetLoveCountRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*postProtoV1.GetLoveCountRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	targetID, err := uuid.Parse(req.TargetId)
	if err != nil {
		return nil, err
	}

	return &models.GetLoveCountReq{
		TargetID:     targetID,
		IsPostParent: req.IsPostTarget,
	}, nil
}

func DecodeGetLoveCountResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*postProtoV1.GetLoveCountReponse)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &models.GetLoveCountResp{
		Count: int(resp.Count),
	}, nil
}

func DecodeGetCommentCountRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*postProtoV1.GetCommentCountRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	parentID, err := uuid.Parse(req.ParentId)
	if err != nil {
		return nil, err
	}

	return &models.GetCommentCountReq{
		ParentID:     parentID,
		IsPostParent: req.IsPostParent,
	}, nil
}

func DecodeGetCommentCountResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*postProtoV1.GetCommentCountReponse)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &models.GetCommentCountResp{
		Count: int(resp.Count),
	}, nil
}

func DecodeRemoveLoveReactRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*postProtoV1.RemoveLoveReactRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	targetID, err := uuid.Parse(req.GetTargetId())
	if err != nil {
		return nil, err
	}
	authorID, err := uuid.Parse(req.GetAuthorId())
	if err != nil {
		return nil, err
	}

	return &models.RemoveLoveReactReq{
		TargetID:     targetID,
		AuthorID:     authorID,
		IsTargetPost: req.GetIsTargetPost(),
	}, nil
}

func DecodeRemoveLoveReactResponse(_ context.Context, response interface{}) (interface{}, error) {
	_, ok := response.(*postProtoV1.RemoveLoveReactResponse)
	if !ok {
		return nil, MustBeProtoRespErr
	}

	return &models.RemoveLoveReactResp{}, nil
}

func DecodeGetPostRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*postProtoV1.GetPostRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, err
	}

	return &models.GetPostReq{
		PostID: postID,
	}, nil
}

func DecodeGetCommentRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*postProtoV1.GetCommentRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	commentID, err := uuid.Parse(req.CommentId)
	if err != nil {
		return nil, err
	}

	return &models.GetCommentReq{
		CommentID: commentID,
	}, nil
}

func DecodeGetLoveRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*postProtoV1.GetLoveRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	authorID, err := uuid.Parse(req.AuthorId)
	if err != nil {
		return nil, err
	}

	targetID, err := uuid.Parse(req.TargetId)
	if err != nil {
		return nil, err
	}

	return &models.GetLoveReq{
		AuthorID: authorID,
		TargetID: targetID,
	}, nil
}

func DecodeListCommentIDsByParentIDRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*postProtoV1.ListCommentIDsByParentIDRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	parentID, err := uuid.Parse(req.ParentId)
	if err != nil {
		return nil, err
	}
	commentID, err := uuid.Parse(req.AfterCommentId)
	if err != nil {
		return nil, err
	}

	return &models.ListCommentIDsByParentIDReq{
		ParentID:       parentID,
		PageSize:       int(req.PageSize),
		AfterCommentID: commentID,
	}, nil
}

func DecodeListCommentIDsByParentIDResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*postProtoV1.ListCommentIDsByParentIDResponse)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	commentIDs := make([]uuid.UUID, 0)
	for _, id_ := range resp.GetCommentIds() {
		id, err := uuid.Parse(id_)
		if err != nil {
			return nil, err
		}

		commentIDs = append(commentIDs, id)
	}

	return &models.ListCommentIDsByParentIDResp{
		CommentIDs: commentIDs,
	}, nil
}

func DecodeListCommentAncestorsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*postProtoV1.ListCommentAncestorsRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	commentID, err := uuid.Parse(req.CommentId)
	if err != nil {
		return nil, err
	}

	return &models.ListCommentAncestorsReq{
		CommentID: commentID,
	}, nil
}

func DecodeListCommentAncestorsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*postProtoV1.ListCommentAncestorsResponse)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	comments := make([]*models.Comment, 0)
	for _, c := range resp.AncestorComments {
		comment, err := decodeCommentProtoModel(c)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	return &models.ListCommentAncestorsResp{
		AncestorComments: comments,
	}, nil
}

func decodePostProtoModel(post *commonProto.Post) (*models.Post, error) {
	id, err := uuid.Parse(post.GetId())
	if err != nil {
		return nil, err
	}
	authorID, err := uuid.Parse(post.GetAuthorId())
	if err != nil {
		return nil, err
	}

	return &models.Post{
		ID:           id,
		AuthorID:     authorID,
		Content:      post.GetContent(),
		Images:       utils.Map2(post.GetImages(), func(i *commonProto.Image) models.Image { return models.Image{URL: i.Uri, Description: i.Description} }),
		Videos:       utils.Map2(post.GetVideos(), func(v *commonProto.Video) models.Video { return models.Video{URL: v.Uri, Description: v.Description} }),
		LoveCount:    int(post.GetLoveCount()),
		CommentCount: int(post.GetCommentCount()),
		Visibility:   post.GetVisibility(),
		Activity:     post.GetActivity(),
		CreatedAt:    post.GetCreatedAt().AsTime(),
		UpdatedAt:    post.GetUpdatedAt().AsTime(),
	}, nil
}

func decodeCommentProtoModel(comment *commonProto.Comment) (*models.Comment, error) {
	id, err := uuid.Parse(comment.Id)
	if err != nil {
		return nil, err
	}
	postId, err := uuid.Parse(comment.PostId)
	if err != nil {
		return nil, err
	}
	authorId, err := uuid.Parse(comment.AuthorId)
	if err != nil {
		return nil, err
	}
	parentId, err := uuid.Parse(comment.ParentId)
	if err != nil {
		return nil, err
	}

	return &models.Comment{
		ID:           id,
		PostID:       postId,
		AuthorID:     authorId,
		ParentID:     parentId,
		IsPostParent: comment.GetIsPostParent(),
		Content:      comment.GetContent(),
		Image: models.Image{
			URL:         comment.GetImage().Uri,
			Description: comment.GetImage().Description,
		},
		Video: models.Video{
			URL:         comment.GetVideo().Uri,
			Description: comment.GetVideo().Description,
		},
		LoveCount:       int(comment.GetLoveCount()),
		SubcommentCount: int(comment.GetSubcommentCount()),
		CreatedAt:       comment.CreatedAt.AsTime(),
		UpdatedAt:       comment.UpdatedAt.AsTime(),
	}, nil
}

func decodeLoveProtoModel(love *commonProto.Love) (*models.Love, error) {
	loveId, err := uuid.Parse(love.Id)
	if err != nil {
		return nil, err
	}
	targetID, err := uuid.Parse(love.TargetId)
	if err != nil {
		return nil, err
	}
	authorId, err := uuid.Parse(love.AuthorId)
	if err != nil {
		return nil, err
	}
	return &models.Love{
		ID:           loveId,
		TargetID:     targetID,
		IsPostTarget: love.IsPostTarget,
		AuthorID:     authorId,
		CreatedAt:    love.CreatedAt.AsTime(),
	}, nil
}
