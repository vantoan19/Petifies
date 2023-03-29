package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID           uuid.UUID
	AuthorID     uuid.UUID
	Content      string
	Images       []Image
	Videos       []Video
	LoveCount    int
	CommentCount int
	Visibility   string
	Activity     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Comment struct {
	ID              uuid.UUID
	PostID          uuid.UUID
	AuthorID        uuid.UUID
	ParentID        uuid.UUID
	IsPostParent    bool
	Content         string
	Image           Image
	Video           Video
	LoveCount       int
	SubcommentCount int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Love struct {
	ID           uuid.UUID `json:"id"`
	TargetID     uuid.UUID `json:"target_id"`
	IsPostTarget bool      `json:"is_post_target"`
	AuthorID     uuid.UUID `json:"author_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type ListCommentsResp struct {
	Comments []*Comment
}

type ListPostsResp struct {
	Posts []*Post
}

type GetLoveCountResp struct {
	Count int
}

type GetCommentCountResp struct {
	Count int
}

type RemoveLoveReactResp struct{}

type ListCommentIDsByParentIDResp struct {
	CommentIDs []uuid.UUID
}

type ListCommentAncestorsResp struct {
	AncestorComments []*Comment
}
