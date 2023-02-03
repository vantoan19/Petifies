package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"github.com/vantoan19/Petifies/server/libs/grpcutils"
	logging "github.com/vantoan19/Petifies/server/libs/logging-config"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/config"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/interceptors/auth"
)

var logger = logging.NewLogger("Mobile.APIGateway")

func setupGRPC() (*grpc.Server, error) {
	logger.Info("Setting up GRPC server for Mobile API Gateway")

	interceptors := grpcutils.ServerInterceptors{
		UnaryInterceptors:  []grpc.UnaryServerInterceptor{auth.Auth.GetUnaryAuthInterceptor()},
		StreamInterceptors: []grpc.StreamServerInterceptor{auth.Auth.GetStreamAuthInterceptor()},
	}

	s, err := grpcutils.NewTLSGrpcServer(config.Conf.TLSKeyPath, config.Conf.TLSCertPath, interceptors)
	if err != nil {
		logger.ErrorData("Failed to create new TLS GRPC server", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Set up GRPC server successfully")
	return s, nil
}

func serveGRPC(grpcServer *grpc.Server) {
	logger.InfoData("Serving GRPC server", logging.Data{"port": config.Conf.GrpcPort})

	grpcEndpoint := fmt.Sprintf(":%d", config.Conf.GrpcPort)
	listener, err := net.Listen("tcp", grpcEndpoint)

	err = grpcServer.Serve(listener)
	if err != nil && err != grpc.ErrServerStopped {
		logger.ErrorData("Failed to serve GRPC server", logging.Data{"error": err.Error(), "port": config.Conf.GrpcPort})
	}

	logger.Info("Shutting down GRPC server for Mobile API Gateway")
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
