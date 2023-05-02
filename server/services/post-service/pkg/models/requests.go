package models

import (
	"github.com/google/uuid"
)

type Image struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

type Video struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

type CreatePostReq struct {
	AuthorID    uuid.UUID
	Visibility  string
	Activity    string
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
	ID         uuid.UUID
	Visibility string
	Activity   string
	Content    string
	Images     []Image
	Videos     []Video
}

type LoveReactReq struct {
	TargetID     uuid.UUID
	AuthorID     uuid.UUID
	IsTargetPost bool
}

type RemoveLoveReactReq struct {
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

type GetLoveCountReq struct {
	TargetID     uuid.UUID
	IsPostParent bool
}

type GetCommentCountReq struct {
	ParentID     uuid.UUID
	IsPostParent bool
}

type GetPostReq struct {
	PostID uuid.UUID
}

type GetCommentReq struct {
	CommentID uuid.UUID
}

type GetLoveReq struct {
	AuthorID uuid.UUID
	TargetID uuid.UUID
}

type ListCommentIDsByParentIDReq struct {
	ParentID       uuid.UUID
	PageSize       int
	AfterCommentID uuid.UUID
}

type ListCommentAncestorsReq struct {
	CommentID uuid.UUID
}
