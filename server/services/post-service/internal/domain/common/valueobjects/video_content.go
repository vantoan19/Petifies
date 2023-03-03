package valueobjects

import "github.com/vantoan19/Petifies/server/libs/common-utils"

type VideoContent struct {
	url         string
	description string
}

func NewVideoContent(url string, description string) VideoContent {
	return VideoContent{url: url, description: description}
}

func (v VideoContent) Validate() (errs common.MultiError) {
	if v.url == "" {
		errs = append(errs, EmptyURLErr)
	}

	return errs
}

func (v VideoContent) IsEmpty() bool {
	return v.url == ""
}

func (v VideoContent) URL() string {
	return v.url
}

func (v VideoContent) Description() string {
	return v.description
}
