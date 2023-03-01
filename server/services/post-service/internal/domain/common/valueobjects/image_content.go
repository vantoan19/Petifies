package valueobjects

import (
	"errors"

	"github.com/vantoan19/Petifies/server/libs/common-utils"
)

type ImageContent struct {
	url         string
	description string
}

var (
	EmptyURLErr = errors.New("image URL cannot be empty")
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

func (i ImageContent) URL() string {
	return i.url
}

func (i ImageContent) Description() string {
	return i.description
}
