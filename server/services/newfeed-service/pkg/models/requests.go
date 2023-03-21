package models

import (
	"time"

	"github.com/google/uuid"
)

type ListPostFeedsReq struct {
	UserID     uuid.UUID
	PageSize   int
	BeforeTime time.Time
}

type ListStoryFeedsReq struct {
	UserID     uuid.UUID
	PageSize   int
	BeforeTime time.Time
}
