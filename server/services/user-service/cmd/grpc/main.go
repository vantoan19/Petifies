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

	userProtoV1 "github.com/vantoan19/Petifies/proto/user-service/v1"
	outbox_dispatcher "github.com/vantoan19/Petifies/server/infrastructure/outbox/dispatcher"
	"github.com/vantoan19/Petifies/server/libs/grpcutils"
	logging "github.com/vantoan19/Petifies/server/libs/logging-config"
	cmd "github.com/vantoan19/Petifies/server/services/user-service/cmd"
	services "github.com/vantoan19/Petifies/server/services/user-service/internal/application/services"
	"github.com/vantoan19/Petifies/server/services/user-service/internal/infra/repositories/user_event/postgres"
	endpointsV1 "github.com/vantoan19/Petifies/server/services/user-service/internal/presentation/endpoints/grpc/v1"
	serversV1 "github.com/vantoan19/Petifies/server/services/user-service/internal/presentation/transport/grpc/v1"
)

var logger = logging.New("UserService.Cmd.Grpc")

func setupGRPC() (*grpc.Server, error) {
	logger.Info("Start setupGRPC")

	interceptors := grpcutils.ServerInterceptors{
		UnaryInterceptors:  []grpc.UnaryServerInterceptor{},
		StreamInterceptors: []grpc.StreamServerInterceptor{},
	}

	s, err := grpcutils.NewInsecureGrpcServer(interceptors)
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
		logger.ErrorData("Finished serveGRPC: FAILED", logging.Data{"error": err.Error(), "port": cmd.Conf.GrpcPort})
		panic(err)
	}

	registerServices(grpcServer)

	reflection.Register(grpcServer)
	err = grpcServer.Serve(listener)
	if err != nil && err != grpc.ErrServerStopped {
		logger.ErrorData("Finished serveGRPC: FAILED", logging.Data{"error": err.Error(), "port": cmd.Conf.GrpcPort})
		panic(err)
	}

	logger.Info("Finished serveGRPC: SUCCESSFUL")
}

func registerServices(grpcServer *grpc.Server) {
	logger.Info("Start registerServices")

	eventRepo, err := postgres.New(cmd.DB)
	if err != nil {
		panic(err)
	}
	// Register user service
	userSvc, err := services.NewUserService(
		services.WithPostgreUserRepository(cmd.DB),
		services.WithKafkaUserEventPublisher(&cmd.UserProducer, eventRepo),
	)
	if err != nil {
		logger.ErrorData("Finished registerServices: FAILED", logging.Data{"error": err.Error()})
		panic(err)
	}

	userEndpoints := endpointsV1.NewUserEndpoints(userSvc)
	authEndpoints := endpointsV1.NewAuthEndpoints(userSvc)
	userProtoV1.RegisterUserServiceServer(grpcServer, serversV1.NewUserServer(userEndpoints))
	userProtoV1.RegisterAuthServiceServer(grpcServer, serversV1.NewAuthServer(authEndpoints))
	logger.Info("Finished registerServices: SUCCESSFUL")
}

func serveEventDispatcher(endSignal <-chan bool) {
	settings := outbox_dispatcher.DispatcherSettings{
		PublishInterval: 15 * time.Second,
		UnlockInterval:  5 * time.Minute,
		CleanInterval:   24 * time.Hour,
		PublishSettings: outbox_dispatcher.PublishSettings{},
		CleanSettings: outbox_dispatcher.CleanSettings{
			EventLifetime: 18 * time.Hour,
		},
	}
	eventRepo, err := postgres.New(cmd.DB)
	errsChan := make(chan error)
	if err != nil {
		panic(err)
	}
	dispatcher := outbox_dispatcher.NewDispatcher(eventRepo, cmd.UserProducer, settings, *logging.New("UserService.Dispatcher"))
	dispatcher.Run(errsChan, endSignal)

	go func() {
		for {
			err = <-errsChan
			fmt.Printf("Received err from dispatcher: %s", err.Error())
		}
	}()

	<-endSignal
}

func actualMain() {
	logger.Info("User Service starting up")
	cmd.Initialize()
	s, err := setupGRPC()
	if err != nil {
		panic(err)
	}
	endSignal := make(chan bool)

	go serveGRPC(s)
	go serveEventDispatcher(endSignal)

	// wait for a terminating signal from the OS
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	sig := <-signalChan
	endSignal <- true

	logger.InfoData("Received signal, shutting down the service", logging.Data{"sig": sig})
	s.GracefulStop()

	// Wait for the server to stop gracefully
	time.Sleep(time.Second)
}

func main() {
	actualMain()
}
