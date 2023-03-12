package cmd

import (
	"github.com/vantoan19/Petifies/server/libs/grpcutils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"google.golang.org/grpc"
)

var (
	RelationshipServiceConn *grpc.ClientConn
)

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
