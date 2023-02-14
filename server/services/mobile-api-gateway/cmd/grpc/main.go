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

	publicProtoV1 "github.com/vantoan19/Petifies/proto/public-gateway/v1"
	"github.com/vantoan19/Petifies/server/libs/grpcutils"
	logging "github.com/vantoan19/Petifies/server/libs/logging-config"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/cmd"
	endpointsV1 "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/endpoints/grpc/v1"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/interceptors/auth"
	services "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/services"
	grpcServers "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/transport/grpc/v1"
)

var logger = logging.New("Mobile.APIGateway")

func setupGRPC() (*grpc.Server, error) {
	logger.Info("Setting up GRPC server for Mobile API Gateway")

	interceptors := grpcutils.ServerInterceptors{
		UnaryInterceptors:  []grpc.UnaryServerInterceptor{auth.Auth.GetUnaryAuthInterceptor()},
		StreamInterceptors: []grpc.StreamServerInterceptor{auth.Auth.GetStreamAuthInterceptor()},
	}

	s, err := grpcutils.NewTLSGrpcServer(cmd.Conf.TLSKeyPath, cmd.Conf.TLSCertPath, interceptors)
	if err != nil {
		logger.ErrorData("Failed to create new TLS GRPC server", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Set up GRPC server successfully")
	return s, nil
}

func serveGRPC(grpcServer *grpc.Server) {
	logger.InfoData("Serving GRPC server", logging.Data{"port": cmd.Conf.GrpcPort})

	grpcEndpoint := fmt.Sprintf(":%d", cmd.Conf.GrpcPort)
	listener, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		logger.ErrorData("Failed to serve GRPC server", logging.Data{"error": err.Error(), "port": cmd.Conf.GrpcPort})
	}

	registerServices(grpcServer)

	reflection.Register(grpcServer)
	err = grpcServer.Serve(listener)
	if err != nil && err != grpc.ErrServerStopped {
		logger.ErrorData("Failed to serve GRPC server", logging.Data{"error": err.Error(), "port": cmd.Conf.GrpcPort})
	}

	logger.Info("Shutting down GRPC server for Mobile API Gateway")
}

func registerServices(grpcServer *grpc.Server) {
	// Register user service
	userSvc, err := services.NewUserService(cmd.UserServiceConn)
	if err != nil {
		logger.ErrorData("Failed to create User Service", logging.Data{"error": err.Error()})
		panic(err)
	}
	userEndpoints := endpointsV1.NewUserEndpoints(userSvc)
	publicProtoV1.RegisterPublicGatewayServer(grpcServer, grpcServers.NewPublicServer(userEndpoints))
}

func actualMain() {
	logger.Info("Mobile API Gateway starting up")
	cmd.Initialize()
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
