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
	logger.Info("Init User Client...")
	conn, err := grpcutils.NewInsecureGrpcClient(
		Conf.UserServiceHost,
		10,
		grpcutils.ClientInterceptors{},
	)
	if err != nil {
		logger.ErrorData("Failed to init user client", logging.Data{"error": err.Error()})
		return err
	}

	logger.Info("Finished Init User Client...")
	UserServiceConn = conn
	return nil
}
