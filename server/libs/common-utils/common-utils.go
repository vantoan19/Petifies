package common

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"time"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/google/uuid"
)

// ===========================

type AuthKey struct{}

type AuthData struct {
	UserID string
}

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	userAuth := ctx.Value(AuthKey{})
	switch md := userAuth.(type) {
	case *AuthData:
		id, err := uuid.Parse(md.UserID)
		if err != nil {
			return uuid.UUID{}, err
		}
		return id, nil
	default:
		return uuid.UUID{}, errors.New("not auth data")
	}
}

// ===========================

type MultiError []error

func (m MultiError) Exist() bool {
	return len(m) > 0
}

func (m MultiError) Error() string {
	s, n := "", 0
	for _, e := range m {
		if e != nil {
			s = s + e.Error() + ",\n"
			n++
		}
	}
	if n == 0 {
		return "0 error"
	}
	return s
}

func IsDevEnv() bool {
	return os.Getenv("SERVER_MODE") == "development"
}

// ============================

type Translator func(context.Context, interface{}) (interface{}, error)

func CreateClientForwardDecodeRequestFunc[T interface{}]() grpctransport.DecodeRequestFunc {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(T)
		if !ok {
			return nil, errors.New("unexpected type of request")
		}
		return req, nil
	}
}

func CreateClientForwardDecodeResponseFunc[T interface{}]() grpctransport.DecodeResponseFunc {
	return func(_ context.Context, response interface{}) (interface{}, error) {
		req, ok := response.(T)
		if !ok {
			return nil, errors.New("unexpected type of response")
		}
		return req, nil
	}
}

func CreateClientForwardEncodeRequestFunc[T interface{}]() grpctransport.EncodeRequestFunc {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(T)
		if !ok {
			return nil, errors.New("unexpected type of request")
		}
		return req, nil
	}
}

func CreateClientForwardEncodeResponseFunc[T interface{}]() grpctransport.EncodeResponseFunc {
	return func(_ context.Context, response interface{}) (interface{}, error) {
		req, ok := response.(T)
		if !ok {
			return nil, errors.New("unexpected type of response")
		}
		return req, nil
	}
}

// ============================

func ToSlice[T interface{}](c chan T) []T {
	s := make([]T, 0)
	for i := range c {
		s = append(s, i)
	}
	return s
}

func Filter[T any](data []T, f func(T) bool) []T {
	fltd := make([]T, 0, len(data))

	for _, e := range data {
		if f(e) {
			fltd = append(fltd, e)
		}
	}
	return fltd
}

func Map2[T, U any](data []T, f func(T) U) []U {

	res := make([]U, 0, len(data))

	for _, e := range data {
		res = append(res, f(e))
	}

	return res
}

func FindFirst[T any](data []T, f func(T) bool) int {
	for idx, e := range data {
		if f(e) {
			return idx
		}
	}
	return -1
}

func NullStringToString(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}

func NullUUIDToUUID(u uuid.NullUUID) *uuid.UUID {
	if u.Valid {
		return &u.UUID
	}
	return nil
}

func NullTimeToTime(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func UUIDToNullUUID(u *uuid.UUID) uuid.NullUUID {
	if u == nil {
		return uuid.NullUUID{
			UUID:  uuid.Nil,
			Valid: false,
		}
	}
	return uuid.NullUUID{
		UUID:  *u,
		Valid: true,
	}
}

func TimeToNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		}
	}
	return sql.NullTime{
		Time:  *t,
		Valid: true,
	}
}

func StringToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{
			String: "",
			Valid:  false,
		}
	}
	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}
