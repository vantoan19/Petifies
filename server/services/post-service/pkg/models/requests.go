package models

import (
	"github.com/google/uuid"
)

type Image struct {
	URL         string
	Description string
}

type Video struct {
	URL         string
	Description string
}

type CreatePostReq struct {
	AuthorID    uuid.UUID
	TextContent string
	Images      []Image
	Videos      []Video
}

type CreateCommentReq struct {
	PostID       uuid.UUID
	AuthorID     uuid.UUID
	ParentID     uuid.UUID
	IsPostParent bool
	Content      string
	ImageContent Image
	VideoContent Video
}

type EditCommentReq struct {
	ID      uuid.UUID
	Content string
	Image   Image
	Video   Video
}

type EditPostReq struct {
	ID      uuid.UUID
	Content string
	Images  []Image
	Videos  []Video
}

type LoveReactReq struct {
	TargetID     uuid.UUID
	AuthorID     uuid.UUID
	IsTargetPost bool
}

type ListCommentsReq struct {
	CommentIDs []uuid.UUID
}

type ListPostsReq struct {
	PostIDs []uuid.UUID
}
