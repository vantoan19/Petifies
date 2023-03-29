package commentaggre

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	loveaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/love"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/entities"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/valueobjects"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
)

var (
	ErrDuplicatedLove    = status.Errorf(codes.AlreadyExists, "a user cannot add love twice")
	ErrLoveNotExists     = status.Errorf(codes.NotFound, "love reaction does not exist")
	ErrNotChildComment   = status.Errorf(codes.InvalidArgument, "parent ID does not identical to comment ID")
	ErrPostParent        = status.Errorf(codes.InvalidArgument, "subcomment cannot have post parent")
	ErrCommentIDNotExist = status.Errorf(codes.NotFound, "comment ID does not exist in the post")
	ErrCommentIDExist    = status.Errorf(codes.AlreadyExists, "comment ID already exists in the post")
)

// Comment represents an aggregate for Comment, Loves and its SubComments
type Comment struct {
	comment     *entities.Comment // root
	loves       []uuid.UUID
	subcomments []uuid.UUID
}

func New(content *models.CreateCommentReq) (*Comment, error) {
	commentEntity := &entities.Comment{
		ID:           uuid.New(),
		PostID:       content.PostID,
		AuthorID:     content.AuthorID,
		ParentID:     content.ParentID,
		IsPostParent: content.IsPostParent,
		Content:      valueobjects.NewTextContent(content.Content),
		ImageContent: valueobjects.NewImageContent(content.ImageContent.URL, content.ImageContent.Description),
		VideoContent: valueobjects.NewVideoContent(content.VideoContent.URL, content.VideoContent.Description),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if errs := commentEntity.Validate(); errs.Exist() {
		return nil, status.Errorf(codes.InvalidArgument, errs[0].Error())
	}

	return &Comment{
		comment:     commentEntity,
		loves:       make([]uuid.UUID, 0),
		subcomments: make([]uuid.UUID, 0),
	}, nil
}

// GetComment returns the Comment entity associated with the aggregate
func (c *Comment) GetCommentEntity() entities.Comment {
	return *c.comment
}

// SetComment sets the Comment entity associated with the aggregate
func (c *Comment) SetCommentEntity(comment entities.Comment) error {
	if errs := comment.Validate(); errs.Exist() {
		return status.Errorf(codes.InvalidArgument, errs[0].Error())
	}

	c.comment = &comment
	return nil
}

// UpdateCommentTextContent update the text content of the comment
func (c *Comment) UpdateCommentTextContent(textContent valueobjects.TextContent) {
	c.comment.Content = textContent
}

// AddLoveByEntity adds a Love to Comment
// This method is used for DTO
// func (c *Comment) AddLoveByEntity(love entities.Love) error {
// 	for _, love_ := range c.loves {
// 		if love_.AuthorID == love.AuthorID {
// 			return ErrDuplicatedLove
// 		}
// 	}
// 	if errs := love.Validate(); errs.Exist() {
// 		return status.Errorf(codes.InvalidArgument, errs[0].Error())
// 	}

// 	c.loves = append(c.loves, &love)
// 	return nil
// }

// AddLoveByAuthorIDAndSave adds a new love to the comment
// and save the love to the db
func (c *Comment) AddLoveByAuthorIDAndSave(authorID uuid.UUID, repo loveaggre.LoveRepository) error {
	if exists, err := repo.ExistsLoveByTargetIDAndAuthorID(context.Background(), authorID, c.comment.ID); err != nil {
		return err
	} else if exists {
		return ErrDuplicatedLove
	}
	loveAggre, err := loveaggre.NewLove(&models.Love{
		ID:           uuid.New(),
		TargetID:     c.comment.ID,
		IsPostTarget: false,
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

	c.loves = append(c.loves, savedLove.GetID())
	return nil
}

// RemoveLoveByAuthorIDAndDelete removes a Love from the comment
// and delete the love int the db
func (c *Comment) RemoveLoveByAuthorIDAndDelete(authorID uuid.UUID, repo loveaggre.LoveRepository) error {
	if exists, err := repo.ExistsLoveByTargetIDAndAuthorID(context.Background(), authorID, c.comment.ID); err != nil {
		return err
	} else if !exists {
		return ErrLoveNotExists
	}
	love, err := repo.GetByTargetIDAndAuthorID(context.Background(), authorID, c.comment.ID)
	if err != nil {
		return err
	}
	err = repo.DeleteByUUID(context.Background(), love.GetID())
	if err != nil {
		return err
	}

	for i, l := range c.loves {
		if l == love.GetID() {
			c.loves = append(c.loves[:i], c.loves[i+1:]...)
			break
		}
	}
	return nil
}

// GetLoves returns the Loves associated with the Comment
func (c *Comment) GetLovesID() []uuid.UUID {
	return c.loves
}

// AddSubcommentByEntity adds a UUID of a subcomment to the Comment
// This method is used for DTO
func (c *Comment) AddSubcommentByEntity(subcomment entities.Comment) error {
	if subcomment.ParentID != c.comment.ID {
		return ErrNotChildComment
	}
	if subcomment.IsPostParent {
		return ErrPostParent
	}
	if errs := subcomment.Validate(); errs.Exist() {
		return status.Errorf(codes.InvalidArgument, errs[0].Error())
	}

	c.subcomments = append(c.subcomments, subcomment.ID)
	return nil
}

// AddSubcommentAndSave adds a UUID of a subcomment to the Comment
// and Save the comment in the repo
func (c *Comment) AddSubcommentAndSave(subcomment *Comment, repo CommentRepository) error {
	if subcomment.GetCommentEntity().ParentID != c.comment.ID {
		return ErrNotChildComment
	}
	if subcomment.GetCommentEntity().IsPostParent {
		return ErrPostParent
	}
	if c.ExistsSubcomment(subcomment.GetCommentEntity().ID) {
		return ErrCommentIDNotExist
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	subcomment, err := repo.SaveComment(ctx, *subcomment)
	if err != nil {
		return err
	}

	c.subcomments = append(c.subcomments, subcomment.GetCommentEntity().ID)
	return nil
}

// RemoveSubcommentAndDelete removes a UUID of a subcomment from the Comment
// and Delete the comment in the repo
func (c *Comment) RemoveSubcommentAndDelete(subcommentID uuid.UUID, repo CommentRepository) error {
	if !c.ExistsSubcomment(subcommentID) {
		return ErrCommentIDNotExist
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	subcomment, err := repo.DeleteByUUID(ctx, subcommentID)
	if err != nil {
		return err
	}

	for i, sc := range c.subcomments {
		if sc == subcomment.GetCommentEntity().ID {
			c.subcomments = append(c.subcomments[:i], c.subcomments[i+1:]...)
			break
		}
	}
	return nil
}

// ExistsSubcomment checks whether a comment ID exists in the subcomments or not
func (c *Comment) ExistsSubcomment(id uuid.UUID) bool {
	for _, sc := range c.subcomments {
		if sc == id {
			return true
		}
	}

	return false
}

func (c *Comment) GetSubcommentsID() []uuid.UUID {
	return c.subcomments
}

// ============= Root Entity Getters =================

func (c *Comment) GetID() uuid.UUID {
	return c.comment.ID
}

func (c *Comment) GetPostID() uuid.UUID {
	return c.comment.PostID
}

func (c *Comment) GetAuthorID() uuid.UUID {
	return c.comment.AuthorID
}

func (c *Comment) GetParentID() uuid.UUID {
	return c.comment.ParentID
}

func (c *Comment) GetIsPostParent() bool {
	return c.comment.IsPostParent
}

func (c *Comment) GetContent() valueobjects.TextContent {
	return c.comment.Content
}

func (c *Comment) GetImageContent() valueobjects.ImageContent {
	return c.comment.ImageContent
}

func (c *Comment) GetVideoContent() valueobjects.VideoContent {
	return c.comment.VideoContent
}

func (c *Comment) GetCreatedAt() time.Time {
	return c.comment.CreatedAt
}

func (c *Comment) GetUpdatedAt() time.Time {
	return c.comment.UpdatedAt
}
