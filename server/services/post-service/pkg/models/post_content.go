package models

import "github.com/google/uuid"

type Image struct {
	URL         string
	Description string
}

type Video struct {
	URL         string
	Description string
}

type PostContent struct {
	AuthorID    uuid.UUID
	TextContent string
	Images      []Image
	Videos      []Video
}
