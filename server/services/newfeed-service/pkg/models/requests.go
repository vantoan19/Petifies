package models

import "github.com/google/uuid"

type ListPostFeedsReq struct {
	UserID uuid.UUID
}

type ListStoryFeedsReq struct {
	UserID uuid.UUID
}
