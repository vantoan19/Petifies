package models

import "github.com/google/uuid"

type ListPostFeedsResp struct {
	PostIDs []uuid.UUID
}

type ListStoryFeedsResp struct {
	StoryIDs []uuid.UUID
}
