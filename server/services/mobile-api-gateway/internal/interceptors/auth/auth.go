package auth

import (
	"context"
	"errors"
	"strings"

	grpcAuth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"

	logging "github.com/vantoan19/Petifies/server/libs/logging-config"
)

var logger = logging.New("MobileGateway.Auth")

type AuthInterceptor struct{}

var Auth AuthInterceptor

func (m *AuthInterceptor) GetUnaryAuthInterceptor() grpc.UnaryServerInterceptor {
	return grpcAuth.UnaryServerInterceptor(authenticate)
}

func (m *AuthInterceptor) GetStreamAuthInterceptor() grpc.StreamServerInterceptor {
	return grpcAuth.StreamServerInterceptor(authenticate)
}

func authenticate(ctx context.Context) (context.Context, error) {
	logger.Info("Start authenticate")

	callingService, err := getGrpcService(ctx)
	if err != nil {
		logger.ErrorData("Finished authenticate: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	// Bypass authenticate for public apis (apis that do not require auth)
	if callingService == "PublicGateway" {
		logger.Info("Finished authenticate: SUCCESSFUL")
		return ctx, nil
	}

	logger.Info("Finished authenticate: SUCCESSFUL")
	return ctx, nil
}

func getGrpcService(ctx context.Context) (string, error) {
	method, ok := grpc.Method(ctx)
	if !ok {
		return "", errors.New("cannot retrieve the method from the context")
	}

	service := strings.Split(method, "/")[1]
	return service, nil
}
