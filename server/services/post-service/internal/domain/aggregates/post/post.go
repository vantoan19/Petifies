package postaggre

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	comment "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/comment"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/entities"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/valueobjects"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
)

var (
	ErrDuplicatedLove    = errors.New("a user cannot add love twice")
	ErrNotChildComment   = errors.New("parent ID does not identical to comment ID")
	ErrNotPostParent     = errors.New("subcomment cannot have post parent")
	ErrCommentIDNotExist = errors.New("comment ID does not exist in the post")
	ErrCommentIDExist    = errors.New("comment ID already exists in the post")
)

type Post struct {
	post     *entities.Post // root
	loves    []*entities.Love
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
		TextContent: valueobjects.NewTextContent(content.TextContent),
		Images:      imageValues,
		Videos:      videoValues,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if errs := postEntity.Validate(); errs.Exist() {
		return nil, errors.New(errs[0].Error())
	}

	return &Post{
		post:     &postEntity,
		comments: make([]uuid.UUID, 0),
		loves:    make([]*entities.Love, 0),
	}, nil
}

func (p *Post) SetPostEntity(post entities.Post) error {
	if errs := post.Validate(); errs.Exist() {
		return errors.New(errs[0].Error())
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
func (p *Post) AddCommentByEntity(comment entities.Comment) error {
	if comment.ParentID != p.post.ID {
		return ErrNotChildComment
	}
	if !comment.IsPostParent {
		return ErrNotPostParent
	}
	if errs := comment.Validate(); errs.Exist() {
		return errors.New(errs[0].Error())
	}

	p.comments = append(p.comments, comment.ID)
	return nil
}

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

// AddLike adds a new like to the post
func (p *Post) AddLoveByEntity(love entities.Love) error {
	for _, love_ := range p.loves {
		if love_.AuthorID == love.AuthorID {
			return ErrDuplicatedLove
		}
	}
	if errs := love.Validate(); errs.Exist() {
		return errors.New(errs[0].Error())
	}

	p.loves = append(p.loves, &love)
	return nil
}

// AddLike adds a new like to the post
func (p *Post) AddLoveByAuthorID(authorID uuid.UUID) error {
	for _, love := range p.loves {
		if love.AuthorID == authorID {
			return ErrDuplicatedLove
		}
	}
	love := &entities.Love{
		ID:        uuid.New(),
		PostID:    p.post.ID,
		AuthorID:  authorID,
		CreatedAt: time.Now(),
	}
	if errs := love.Validate(); errs.Exist() {
		return errors.New(errs[0].Error())
	}

	p.loves = append(p.loves, love)
	return nil
}

// RemoveLove removes a Love from the Post
func (p *Post) RemoveLove(authorID uuid.UUID) {
	for i, l := range p.loves {
		if l.AuthorID == authorID {
			p.loves = append(p.loves[:i], p.loves[i+1:]...)
			break
		}
	}
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
		return errors.New(errs[0].Error())
	}
	p.post.Images = append(p.post.Images, image)
	return nil
}

func (p *Post) AddNewVideo(video valueobjects.VideoContent) error {
	if errs := video.Validate(); errs.Exist() {
		return errors.New(errs[0].Error())
	}
	p.post.Videos = append(p.post.Videos, video)
	return nil
}

// =============== Aggregate Entities Getters ================

func (p *Post) GetLoves() []entities.Love {
	loves := make([]entities.Love, 0)
	for _, love := range p.loves {
		loves = append(loves, *love)
	}

	return loves
}

func (p *Post) GetLovesByAuthorID(authorId uuid.UUID) entities.Love {
	for _, love := range p.loves {
		if love.AuthorID == authorId {
			return *love
		}
	}

	return entities.Love{}
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
