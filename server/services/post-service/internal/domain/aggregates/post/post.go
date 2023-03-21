package postaggre

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	comment "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/comment"
	loveaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/love"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/entities"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/valueobjects"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
)

var (
	ErrDuplicatedLove    = status.Errorf(codes.AlreadyExists, "a user cannot add love twice")
	ErrLoveNotExists     = status.Errorf(codes.NotFound, "love reaction does not exist")
	ErrNotChildComment   = status.Errorf(codes.InvalidArgument, "parent ID does not identical to comment ID")
	ErrNotPostParent     = status.Errorf(codes.InvalidArgument, "subcomment cannot have post parent")
	ErrCommentIDNotExist = status.Errorf(codes.NotFound, "comment ID does not exist in the post")
	ErrCommentIDExist    = status.Errorf(codes.AlreadyExists, "comment ID already exists in the post")
)

type Post struct {
	post     *entities.Post // root
	loves    []uuid.UUID
	comments []uuid.UUID
}

func NewPost(content *models.CreatePostReq) (*Post, error) {
	imageValues := make([]valueobjects.ImageContent, 0)
	videoValues := make([]valueobjects.VideoContent, 0)

	for _, image := range content.Images {
		imageValues = append(imageValues, valueobjects.NewImageContent(image.URL, image.Description))
	}
	for _, video := range content.Videos {
		videoValues = append(videoValues, valueobjects.NewVideoContent(video.URL, video.Description))
	}

	postEntity := entities.Post{
		ID:          uuid.New(),
		AuthorID:    content.AuthorID,
		Visibility:  valueobjects.Visibility(content.Visibility),
		Activity:    content.Activity,
		TextContent: valueobjects.NewTextContent(content.TextContent),
		Images:      imageValues,
		Videos:      videoValues,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if errs := postEntity.Validate(); errs.Exist() {
		return nil, status.Errorf(codes.InvalidArgument, errs[0].Error())
	}

	return &Post{
		post:     &postEntity,
		comments: make([]uuid.UUID, 0),
		loves:    make([]uuid.UUID, 0),
	}, nil
}

func (p *Post) SetPostEntity(post entities.Post) error {
	if errs := post.Validate(); errs.Exist() {
		return status.Errorf(codes.InvalidArgument, errs[0].Error())
	}

	p.post = &post
	return nil
}

func (p *Post) GetPostEntity() entities.Post {
	return *p.post
}

func (p *Post) UpdateTextContent(content valueobjects.TextContent) {
	p.post.TextContent = content
}

func (p *Post) AddImage(image valueobjects.ImageContent) error {
	return p.post.AddImageContent(image)
}

func (p *Post) AddVideo(video valueobjects.VideoContent) error {
	return p.post.AddVideoContent(video)
}

// AddSubcommentByEntity adds a UUID of a subcomment to the Comment
// This method is used for DTO

// func (p *Post) AddCommentByEntity(comment entities.Comment) error {
// 	if comment.ParentID != p.post.ID {
// 		return ErrNotChildComment
// 	}
// 	if !comment.IsPostParent {
// 		return ErrNotPostParent
// 	}
// 	if errs := comment.Validate(); errs.Exist() {
// 		return status.Errorf(codes.InvalidArgument, errs[0].Error())
// 	}

// 	p.comments = append(p.comments, comment.ID)
// 	return nil
// }

// AddComment adds a new comment to the post
// and Save the comment in the repo
func (p *Post) AddCommentAndSave(comment *comment.Comment, repo comment.CommentRepository) error {
	if comment.GetCommentEntity().ParentID != p.post.ID {
		return ErrNotChildComment
	}
	if !comment.GetCommentEntity().IsPostParent {
		return ErrNotPostParent
	}
	if p.ExistsComment(comment.GetCommentEntity().ID) {
		return ErrCommentIDExist
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	comment, err := repo.SaveComment(ctx, *comment)
	if err != nil {
		return err
	}

	p.comments = append(p.comments, comment.GetCommentEntity().ID)
	return nil
}

// RemoveCommentAndDelete remove comment uuid from post comment
// and Delete the comment in repo
func (p *Post) RemoveCommentAndDelete(commentID uuid.UUID, repo comment.CommentRepository) error {
	if !p.ExistsComment(commentID) {
		return ErrCommentIDNotExist
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	comment, err := repo.DeleteByUUID(ctx, commentID)
	if err != nil {
		return err
	}

	for i, c := range p.comments {
		if c == comment.GetCommentEntity().ID {
			p.comments = append(p.comments[:i], p.comments[i+1:]...)
			break
		}
	}
	return nil
}

// AddLoveByEntity adds a new like to the post

// func (p *Post) AddLoveByEntity(love entities.Love) error {
// 	for _, love_ := range p.loves {
// 		if love_.AuthorID == love.AuthorID {
// 			return ErrDuplicatedLove
// 		}
// 	}
// 	if errs := love.Validate(); errs.Exist() {
// 		return status.Errorf(codes.InvalidArgument, errs[0].Error())
// 	}

// 	p.loves = append(p.loves, &love)
// 	return nil
// }

// AddLoveByAuthorIDAndSave adds a new love to the post
// and save the love to the db
func (p *Post) AddLoveByAuthorIDAndSave(authorID uuid.UUID, repo loveaggre.LoveRepository) error {
	if l, err := repo.GetByTargetIDAndAuthorID(context.Background(), authorID, p.post.ID); l != nil || err != nil {
		if err != nil {
			return err
		}
		return ErrDuplicatedLove
	}
	loveAggre, err := loveaggre.NewLove(&models.Love{
		ID:           uuid.New(),
		TargetID:     p.post.ID,
		IsPostTarget: true,
		AuthorID:     authorID,
		CreatedAt:    time.Now(),
	})
	if err != nil {
		return err
	}

	savedLove, err := repo.SaveLove(context.Background(), *loveAggre)
	if err != nil {
		return err
	}

	p.loves = append(p.loves, savedLove.GetID())
	return nil
}

// RemoveLoveByAuthorIDAndDelete removes a Love from the Post
// and delete the love int the db
func (p *Post) RemoveLoveByAuthorIDAndDelete(authorID uuid.UUID, repo loveaggre.LoveRepository) error {
	love, err := repo.GetByTargetIDAndAuthorID(context.Background(), authorID, p.post.ID)
	if err != nil {
		return err
	}
	if love == nil {
		return ErrLoveNotExists
	}

	for i, l := range p.loves {
		if l == love.GetID() {
			p.loves = append(p.loves[:i], p.loves[i+1:]...)
			break
		}
	}

	return nil
}

// ========== Aggregate Root Getters ===========

func (p *Post) GetPostID() uuid.UUID {
	return p.post.ID
}

func (p *Post) GetPostTextContent() valueobjects.TextContent {
	return p.post.TextContent
}

func (p *Post) GetAuthorID() uuid.UUID {
	return p.post.AuthorID
}

func (p *Post) GetVisibility() valueobjects.Visibility {
	return p.post.Visibility
}

func (p *Post) GetActivity() string {
	return p.post.Activity
}

func (p *Post) GetImages() []valueobjects.ImageContent {
	return p.post.Images
}

func (p *Post) GetVideos() []valueobjects.VideoContent {
	return p.post.Videos
}

func (p *Post) GetCreatedAt() time.Time {
	return p.post.CreatedAt
}

func (p *Post) GetUpdatedAt() time.Time {
	return p.post.UpdatedAt
}

func (p *Post) SetPostTextContent(content valueobjects.TextContent) {
	p.post.TextContent = content
}

func (p *Post) RemoveAllImages() {
	p.post.Images = []valueobjects.ImageContent{}
}

func (p *Post) RemoveAllVideos() {
	p.post.Videos = []valueobjects.VideoContent{}
}

func (p *Post) AddNewImage(image valueobjects.ImageContent) error {
	if errs := image.Validate(); errs.Exist() {
		return status.Errorf(codes.InvalidArgument, errs[0].Error())
	}
	p.post.Images = append(p.post.Images, image)
	return nil
}

func (p *Post) AddNewVideo(video valueobjects.VideoContent) error {
	if errs := video.Validate(); errs.Exist() {
		return status.Errorf(codes.InvalidArgument, errs[0].Error())
	}
	p.post.Videos = append(p.post.Videos, video)
	return nil
}

// =============== Aggregate Entities Getters ================

func (p *Post) GetLovesID() []uuid.UUID {
	return p.loves
}

func (p *Post) GetComments() []uuid.UUID {
	return p.comments
}

func (p *Post) ExistsComment(id uuid.UUID) bool {
	for _, c := range p.comments {
		if c == id {
			return true
		}
	}

	return false
}
