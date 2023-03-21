package cmd

import (
	"github.com/vantoan19/Petifies/server/libs/grpcutils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"google.golang.org/grpc"
)

var (
	UserServiceConn         *grpc.ClientConn
	PostServiceConn         *grpc.ClientConn
	MediaServiceConn        *grpc.ClientConn
	RelationshipServiceConn *grpc.ClientConn
	NewfeedServiceConn      *grpc.ClientConn
)

func initUserServiceClient() error {
	logger.Info("Start initUserServiceClient")

	conn, err := grpcutils.NewInsecureGrpcClient(
		Conf.UserServiceHost,
		10,
		grpcutils.ClientInterceptors{},
	)
	if err != nil {
		logger.ErrorData("Finished initUserServiceClient: FAILED", logging.Data{"error": err.Error()})
		return err
	}

	logger.Info("Finished initUserServiceClient: SUCCESSFUL")
	UserServiceConn = conn
	return nil
}

func initPostServiceClient() error {
	logger.Info("Start initPostServiceClient")

	conn, err := grpcutils.NewInsecureGrpcClient(
		Conf.PostServiceHost,
		10,
		grpcutils.ClientInterceptors{},
	)
	if err != nil {
		logger.ErrorData("Finished initPostServiceClient: FAILED", logging.Data{"error": err.Error()})
		return err
	}

	logger.Info("Finished initPostServiceClient: SUCCESSFUL")
	PostServiceConn = conn
	return nil
}

func initMediaServiceClient() error {
	logger.Info("Start initMediaServiceClient")

	conn, err := grpcutils.NewInsecureGrpcClient(
		Conf.MediaServiceHost,
		10,
		grpcutils.ClientInterceptors{},
	)
	if err != nil {
		logger.ErrorData("Finished initMediaServiceClient: FAILED", logging.Data{"error": err.Error()})
		return err
	}

	logger.Info("Finished initMediaServiceClient: SUCCESSFUL")
	MediaServiceConn = conn
	return nil
}

func initRelationshipServiceClient() error {
	logger.Info("Start initRelationshipServiceClient")

	conn, err := grpcutils.NewInsecureGrpcClient(
		Conf.RelationshipServiceHost,
		10,
		grpcutils.ClientInterceptors{},
	)
	if err != nil {
		logger.ErrorData("Finished initRelationshipServiceClient: FAILED", logging.Data{"error": err.Error()})
		return err
	}

	logger.Info("Finished initRelationshipServiceClient: SUCCESSFUL")
	RelationshipServiceConn = conn
	return nil
}

func initNewfeedServiceClient() error {
	logger.Info("Start initNewfeedServiceClient")

	conn, err := grpcutils.NewInsecureGrpcClient(
		Conf.NewfeedServiceHost,
		10,
		grpcutils.ClientInterceptors{},
	)
	if err != nil {
		logger.ErrorData("Finished initNewfeedServiceClient: FAILED", logging.Data{"error": err.Error()})
		return err
	}

	logger.Info("Finished initNewfeedServiceClient: SUCCESSFUL")
	NewfeedServiceConn = conn
	return nil
}
