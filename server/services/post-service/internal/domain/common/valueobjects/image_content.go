package valueobjects

import (
	"github.com/vantoan19/Petifies/server/libs/common-utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ImageContent struct {
	url         string
	description string
}

var (
	EmptyURLErr = status.Errorf(codes.InvalidArgument, "image URL cannot be empty")
)

func NewImageContent(url string, description string) ImageContent {
	return ImageContent{url: url, description: description}
}

func (i ImageContent) Validate() (errs common.MultiError) {
	if i.url == "" {
		errs = append(errs, EmptyURLErr)
	}

	return errs
}

func (i ImageContent) IsEmpty() bool {
	return i.url == ""
}

func (i ImageContent) URL() string {
	return i.url
}

func (i ImageContent) Description() string {
	return i.description
}
