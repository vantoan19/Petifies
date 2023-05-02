package models

import "github.com/google/uuid"

type ListPostFeedsResp struct {
	PostIDs       []uuid.UUID
	NextPageToken string
}

type ListStoryFeedsResp struct {
	StoryIDs      []uuid.UUID
	NextPageToken string
}
