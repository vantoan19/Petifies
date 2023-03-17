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

	authProtoV1 "github.com/vantoan19/Petifies/proto/auth-gateway/v1"
	publicProtoV1 "github.com/vantoan19/Petifies/proto/public-gateway/v1"
	"github.com/vantoan19/Petifies/server/libs/grpcutils"
	logging "github.com/vantoan19/Petifies/server/libs/logging-config"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/cmd"
	postService "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/application/services/post"
	services "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/application/services/user"
	endpointsV1 "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/presentation/endpoints/grpc/v1"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/presentation/interceptors/auth"
	grpcServers "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/presentation/transport/grpc/v1"
)

var logger = logging.New("MobileGateway.Cmd.Grpc")

func setupGRPC() (*grpc.Server, error) {
	logger.Info("Start setupGRPC")

	authInterceptor := auth.New(cmd.UserServiceConn)
	interceptors := grpcutils.ServerInterceptors{
		UnaryInterceptors:  []grpc.UnaryServerInterceptor{authInterceptor.GetUnaryAuthInterceptor()},
		StreamInterceptors: []grpc.StreamServerInterceptor{authInterceptor.GetStreamAuthInterceptor()},
	}

	s, err := grpcutils.NewTLSGrpcServer(cmd.Conf.TLSKeyPath, cmd.Conf.TLSCertPath, interceptors)
	if err != nil {
		logger.ErrorData("Finished setupGRPC: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished setupGRPC: SUCCESSFUL")
	return s, nil
}

func serveGRPC(grpcServer *grpc.Server) {
	logger.InfoData("Start serveGRPC", logging.Data{"port": cmd.Conf.GrpcPort})

	grpcEndpoint := fmt.Sprintf(":%d", cmd.Conf.GrpcPort)
	listener, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		logger.ErrorData("Finished serveGRPC: FAILED", logging.Data{
			"error": err.Error(),
			"port":  cmd.Conf.GrpcPort,
		})
		panic(err)
	}

	registerServices(grpcServer)

	reflection.Register(grpcServer)
	err = grpcServer.Serve(listener)
	if err != nil && err != grpc.ErrServerStopped {
		logger.ErrorData("Finished serveGRPC: FAILED", logging.Data{
			"error": err.Error(),
			"port":  cmd.Conf.GrpcPort,
		})
		panic(err)
	}

	logger.Info("Finished serveGRPC: SUCCESSFUL")
}

func registerServices(grpcServer *grpc.Server) {
	logger.Info("Start registerServices")

	// Register user service
	userSvc, err := services.NewUserService(cmd.UserServiceConn)
	if err != nil {
		logger.ErrorData("Finished registerServices: FAILED", logging.Data{"error": err.Error()})
		panic(err)
	}
	postSvc, err := postService.NewPostService(cmd.PostServiceConn)

	userEndpoints := endpointsV1.NewUserEndpoints(userSvc)
	postEndpoints := endpointsV1.NewPostEndpoints(postSvc)
	publicProtoV1.RegisterPublicGatewayServer(grpcServer, grpcServers.NewPublicServer(userEndpoints))
	authProtoV1.RegisterAuthGatewayServer(grpcServer, grpcServers.NewAuthServer(cmd.MediaServiceConn, userEndpoints, postEndpoints))

	logger.Info("Finished registerServices: SUCCESSFUL")
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
