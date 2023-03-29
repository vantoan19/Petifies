package models

import (
	"github.com/google/uuid"
)

type ListPostFeedsReq struct {
	UserID      uuid.UUID
	PageSize    int
	AfterPostID uuid.UUID
}

type ListStoryFeedsReq struct {
	UserID       uuid.UUID
	PageSize     int
	AfterStoryID uuid.UUID
}
