package auth

import (
	"context"
	"errors"
	"strings"

	grpcAuth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	logging "github.com/vantoan19/Petifies/server/libs/logging-config"
	userclient "github.com/vantoan19/Petifies/server/services/grpc-clients/user-client"
)

var logger = logging.New("MobileGateway.Auth")

const (
	authHeaderKey    = "authorization"
	authHeaderPrefix = "Bearer "
)

type AuthInterceptor struct {
	userClient userclient.UserClient
}

func New(conn *grpc.ClientConn) *AuthInterceptor {
	return &AuthInterceptor{
		userClient: userclient.New(conn),
	}
}

func (m *AuthInterceptor) GetUnaryAuthInterceptor() grpc.UnaryServerInterceptor {
	return grpcAuth.UnaryServerInterceptor(m.authenticate)
}

func (m *AuthInterceptor) GetStreamAuthInterceptor() grpc.StreamServerInterceptor {
	return grpcAuth.StreamServerInterceptor(m.authenticate)
}

func (m *AuthInterceptor) authenticate(ctx context.Context) (context.Context, error) {
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

	token, err := extractAuthMetadata(ctx)
	if err != nil {
		logger.ErrorData("Finished authenticate: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	userID, err := m.userClient.VerifyToken(context.Background(), token)
	newCtx := removeAuthFromCtx(ctx)

	logger.Info("Finished authenticate: SUCCESSFUL")
	return context.WithValue(newCtx, "user_id", userID), nil
}

func extractAuthMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "metadata is missing")
	}

	values := md[authHeaderKey]
	if len(values) != 1 {
		return "", status.Errorf(codes.Unauthenticated, "no or more than 1 auth metadata in the request")
	}

	authMd := values[0]
	if !strings.HasPrefix(authMd, authHeaderPrefix) {
		return "", status.Errorf(codes.Unauthenticated, "missing \"%s\" prefix in \"%s\" header", authHeaderPrefix, authHeaderKey)
	}

	return strings.TrimPrefix(authMd, authHeaderPrefix), nil
}

func getGrpcService(ctx context.Context) (string, error) {
	method, ok := grpc.Method(ctx)
	if !ok {
		return "", errors.New("cannot retrieve the method from the context")
	}

	service := strings.Split(method, "/")[1]
	return service, nil
}

func removeAuthFromCtx(ctx context.Context) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)
	mdCopy := md.Copy()
	mdCopy[authHeaderKey] = nil
	return metadata.NewIncomingContext(ctx, mdCopy)
}
