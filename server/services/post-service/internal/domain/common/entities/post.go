package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/valueobjects"
)

var (
	ErrEmptyID       = status.Errorf(codes.InvalidArgument, "id is empty")
	ErrEmptyAuthorID = status.Errorf(codes.InvalidArgument, "author id is empty")
	ErrEmptyActivity = status.Errorf(codes.InvalidArgument, "empty activity")
)

type Post struct {
	ID          uuid.UUID
	AuthorID    uuid.UUID
	Visibility  valueobjects.Visibility
	Activity    string
	TextContent valueobjects.TextContent
	Images      []valueobjects.ImageContent
	Videos      []valueobjects.VideoContent
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Validate validate Post entity
func (p *Post) Validate() (errs common.MultiError) {
	if p.ID == uuid.Nil {
		errs = append(errs, ErrEmptyID)
	}
	if p.AuthorID == uuid.Nil {
		errs = append(errs, ErrEmptyAuthorID)
	}
	for _, image := range p.Images {
		if imgErrs := image.Validate(); imgErrs.Exist() {
			errs = append(errs, imgErrs...)
			break
		}
	}
	for _, video := range p.Videos {
		if vdErrs := video.Validate(); vdErrs.Exist() {
			errs = append(errs, vdErrs...)
			break
		}
	}
	if p.Activity == "" {
		errs = append(errs, ErrEmptyActivity)
	}
	return errs
}

// HasTextContent returns true if the post has text content.
func (p *Post) HasTextContent() bool {
	return p.TextContent.IsEmpty()
}

// HasImageContent returns true if the post has at least one image content.
func (p *Post) HasImageContent() bool {
	return len(p.Images) > 0
}

// HasVideoContent returns true if the post has at least one video content.
func (p *Post) HasVideoContent() bool {
	return len(p.Videos) > 0
}

// UpdateTextContent updates the text content of the post.
func (p *Post) UpdateTextContent(content valueobjects.TextContent) {
	p.TextContent = content
}

// AddImageContent adds an image content to the post.
func (p *Post) AddImageContent(content valueobjects.ImageContent) error {
	if errs := content.Validate(); errs.Exist() {
		return status.Errorf(codes.InvalidArgument, errs[0].Error())
	}
	p.Images = append(p.Images, content)
	return nil
}

// AddVideoContent adds a video content to the post.
func (p *Post) AddVideoContent(content valueobjects.VideoContent) error {
	if errs := content.Validate(); errs.Exist() {
		return status.Errorf(codes.InvalidArgument, errs[0].Error())
	}
	p.Videos = append(p.Videos, content)
	return nil
}

// RemoveImageContentByURL removes an image content from the post by URL.
func (p *Post) RemoveImageContentByURL(url string) error {
	for i, content := range p.Images {
		if content.URL() == url {
			p.Images = append(p.Images[:i], p.Images[i+1:]...)
			return nil
		}
	}
	return status.Errorf(codes.NotFound, fmt.Sprintf("image content with URL %q not found", url))
}

// RemoveVideoContentByURL removes a video content from the post by URL.
func (p *Post) RemoveVideoContentByURL(url string) error {
	for i, content := range p.Videos {
		if content.URL() == url {
			p.Videos = append(p.Videos[:i], p.Videos[i+1:]...)
			return nil
		}
	}
	return status.Errorf(codes.NotFound, fmt.Sprintf("video content with URL %q not found", url))
}
