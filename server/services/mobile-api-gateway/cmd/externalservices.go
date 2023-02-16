package cmd

import (
	"github.com/vantoan19/Petifies/server/libs/grpcutils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"google.golang.org/grpc"
)

var (
	UserServiceConn *grpc.ClientConn
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
