package commentaggre

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/entities"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/valueobjects"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
)

var (
	ErrDuplicatedLove    = errors.New("a user cannot add love twice")
	ErrNotChildComment   = errors.New("parent ID does not identical to comment ID")
	ErrPostParent        = errors.New("subcomment cannot have post parent")
	ErrCommentIDNotExist = errors.New("comment ID does not exist in the post")
	ErrCommentIDExist    = errors.New("comment ID already exists in the post")
)

// Comment represents an aggregate for Comment, Loves and its SubComments
type Comment struct {
	comment     *entities.Comment // root
	loves       []*entities.Love
	subcomments []uuid.UUID
}

func New(content models.CommentContent) (*Comment, error) {
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
		return nil, errors.New(errs[0].Error())
	}

	return &Comment{
		comment:     commentEntity,
		loves:       make([]*entities.Love, 0),
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
		return errors.New(errs[0].Error())
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
func (c *Comment) AddLoveByEntity(love entities.Love) error {
	for _, love_ := range c.loves {
		if love_.AuthorID == love.AuthorID {
			return ErrDuplicatedLove
		}
	}

	if errs := love.Validate(); errs.Exist() {
		return errors.New(errs[0].Error())
	}

	c.loves = append(c.loves, &love)
	return nil
}

// AddLoveByAuthorID adds a Love to the Comment
func (c *Comment) AddLoveByAuthorID(authorID uuid.UUID) error {
	for _, love := range c.loves {
		if love.AuthorID == authorID {
			return ErrDuplicatedLove
		}
	}

	love := &entities.Love{
		ID:        uuid.New(),
		CommentID: c.comment.ID,
		AuthorID:  authorID,
		CreatedAt: time.Now(),
	}
	if errs := love.Validate(); errs.Exist() {
		return errors.New(errs[0].Error())
	}

	c.loves = append(c.loves, love)
	return nil
}

// RemoveLoveByAuthorID removes a Love from the Comment
func (c *Comment) RemoveLoveByAuthorID(authorID uuid.UUID) {
	for i, l := range c.loves {
		if l.AuthorID == authorID {
			c.loves = append(c.loves[:i], c.loves[i+1:]...)
			break
		}
	}
}

// GetLoves returns the Loves associated with the Comment
func (c *Comment) GetLoves() []entities.Love {
	loves := make([]entities.Love, 0)
	for _, love := range c.loves {
		loves = append(loves, *love)
	}

	return loves
}

// AddSubcomment adds a UUID of a subcomment to the Comment
// This method is used for DTO
func (c *Comment) AddSubcomment(subcomment *Comment) error {
	if subcomment.GetCommentEntity().ParentID != c.comment.ID {
		return ErrNotChildComment
	}
	if subcomment.GetCommentEntity().IsPostParent {
		return ErrPostParent
	}

	c.subcomments = append(c.subcomments, subcomment.GetCommentEntity().ID)
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
