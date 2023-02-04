package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	authProtoV1 "github.com/vantoan19/Petifies/proto/auth-service/v1"
	"github.com/vantoan19/Petifies/server/libs/grpcutils"
	logging "github.com/vantoan19/Petifies/server/libs/logging-config"
	"github.com/vantoan19/Petifies/server/services/user-services/auth-service/internal/config"
	authenEndpointsV1 "github.com/vantoan19/Petifies/server/services/user-services/auth-service/internal/endpoints/grpc/v1/authenticate"
	authenService "github.com/vantoan19/Petifies/server/services/user-services/auth-service/internal/service/authenticate"
	authServerV1 "github.com/vantoan19/Petifies/server/services/user-services/auth-service/internal/transport/grpc/v1"
)

var logger = logging.NewLogger("AuthService")

func setupGRPC() (*grpc.Server, error) {
	logger.Info("Setting up GRPC server for Auth Service")

	interceptors := grpcutils.ServerInterceptors{
		UnaryInterceptors:  []grpc.UnaryServerInterceptor{},
		StreamInterceptors: []grpc.StreamServerInterceptor{},
	}

	s, err := grpcutils.NewInsecureGrpcServer(interceptors)
	if err != nil {
		logger.ErrorData("Failed to create new GRPC server", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Set up GRPC server successfully")
	return s, nil
}

func serveGRPC(grpcServer *grpc.Server) {
	logger.InfoData("Serving GRPC server", logging.Data{"port": config.Conf.GrpcPort})

	grpcEndpoint := fmt.Sprintf(":%d", config.Conf.GrpcPort)
	listener, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		logger.ErrorData("Failed to serve GRPC server", logging.Data{"error": err.Error(), "port": config.Conf.GrpcPort})
	}

	authenSvc := authenService.NewAuthenticateService()
	authenEndpoints := authenEndpointsV1.MakeAuthenticateEndpoint(authenSvc)
	authProtoV1.RegisterAuthServer(grpcServer, authServerV1.NewGRPCAuthServer(authenEndpoints))

	reflection.Register(grpcServer)
	err = grpcServer.Serve(listener)
	if err != nil && err != grpc.ErrServerStopped {
		logger.ErrorData("Failed to serve GRPC server", logging.Data{"error": err.Error(), "port": config.Conf.GrpcPort})
	}

	logger.Info("Shutting down GRPC server")
}

func actualMain() {
	logger.Info("Mobile API Gateway starting up")
	initialize()
	s, err := setupGRPC()
	if err != nil {
		panic(err)
	}

	go serveGRPC(s)

	// wait for a terminating signal from the OS
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	sig := <-signalChan

	logger.InfoData("Received signal, shutting down the API Gateway", logging.Data{"sig": sig})
	s.GracefulStop()

	// Wait for the server to stop gracefully
	time.Sleep(time.Second)
}

func main() {
	actualMain()
}
