package valueobjects

import (
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vantoan19/Petifies/server/libs/common-utils"
)

var (
	ErrInvalidTimeRange = status.Errorf(codes.InvalidArgument, "fromTime have to be less than toTime")
)

type TimeRange struct {
	fromTime time.Time
	toTime   time.Time
}

func NewTimeRange(fromTime, toTime time.Time) TimeRange {
	return TimeRange{
		fromTime: fromTime,
		toTime:   toTime,
	}
}

func (t TimeRange) Validate() (errs common.MultiError) {
	if t.fromTime.After(t.toTime) {
		errs = append(errs, ErrInvalidTimeRange)
	}

	return errs
}

func (t TimeRange) Intersects(another TimeRange) bool {
	if another.toTime.Before(t.fromTime) || t.toTime.Before(another.fromTime) {
		return false
	}
	return true
}

func (t TimeRange) GetFromTime() time.Time {
	return t.fromTime
}

func (t TimeRange) GetToTime() time.Time {
	return t.toTime
}
