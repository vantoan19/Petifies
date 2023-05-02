package valueobjects

import (
	"github.com/vantoan19/Petifies/server/libs/common-utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	EmptyURLErr = status.Errorf(codes.InvalidArgument, "image uri cannot be empty")
)

type Image struct {
	uri         string
	description string
}

func NewImage(uri string, description string) Image {
	return Image{uri: uri, description: description}
}

func (i Image) Validate() (errs common.MultiError) {
	if i.uri == "" {
		errs = append(errs, EmptyURLErr)
	}

	return errs
}

func (i Image) GetURI() string {
	return i.uri
}

func (i Image) GetDescription() string {
	return i.description
}
