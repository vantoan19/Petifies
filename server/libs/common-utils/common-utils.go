package common

import (
	"context"
	"errors"
	"os"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

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
