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

	postProtoV1 "github.com/vantoan19/Petifies/proto/post-service/v1"
	outbox_dispatcher "github.com/vantoan19/Petifies/server/infrastructure/outbox/dispatcher"
	"github.com/vantoan19/Petifies/server/libs/grpcutils"
	logging "github.com/vantoan19/Petifies/server/libs/logging-config"
	cmd "github.com/vantoan19/Petifies/server/services/post-service/cmd"
	commentservice "github.com/vantoan19/Petifies/server/services/post-service/internal/application/services/comment"
	postservice "github.com/vantoan19/Petifies/server/services/post-service/internal/application/services/post"
	eventMongoRepo "github.com/vantoan19/Petifies/server/services/post-service/internal/infra/repositories/post_event/mongo"
	v1 "github.com/vantoan19/Petifies/server/services/post-service/internal/presentation/endpoints/grpc/v1"
	serversV1 "github.com/vantoan19/Petifies/server/services/post-service/internal/presentation/transport/grpc/v1"
)

var logger = logging.New("PostService.Cmd.Grpc")

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

	eventRepo, err := eventMongoRepo.New(cmd.MongoClient)
	if err != nil {
		panic(err)
	}
	// Register user service
	postSvc, err := postservice.NewPostService(
		postservice.WithMongoPostRepository(cmd.MongoClient),
		postservice.WithMongoCommentRepository(cmd.MongoClient),
		postservice.WithKafkaPostEventPublisher(&cmd.PostProducer, eventRepo),
	)
	if err != nil {
		panic(err)
	}
	commentSvc, err := commentservice.NewCommentService(
		commentservice.WithMongoCommentRepository(cmd.MongoClient),
		commentservice.WithMongoPostRepository(cmd.MongoClient),
	)
	if err != nil {
		panic(err)
	}

	postEndpointsV1 := v1.NewPostEndpoints(postSvc, commentSvc)
	if err != nil {
		logger.ErrorData("Finished registerServices: FAILED", logging.Data{"error": err.Error()})
		panic(err)
	}

	postProtoV1.RegisterPostServiceServer(grpcServer, serversV1.NewPostServer(postEndpointsV1))
	logger.Info("Finished registerServices: SUCCESSFUL")
}

func servePostEventDispatcher(endSignal <-chan bool) {
	settings := outbox_dispatcher.DispatcherSettings{
		PublishInterval: 15 * time.Second,
		UnlockInterval:  5 * time.Minute,
		CleanInterval:   24 * time.Hour,
		PublishSettings: outbox_dispatcher.PublishSettings{},
		CleanSettings: outbox_dispatcher.CleanSettings{
			EventLifetime: 18 * time.Hour,
		},
	}
	eventRepo, err := eventMongoRepo.New(cmd.MongoClient)
	errsChan := make(chan error)
	if err != nil {
		panic(err)
	}
	dispatcher := outbox_dispatcher.NewDispatcher(eventRepo, cmd.PostProducer, settings, *logging.New("PostService.Dispatcher"))
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
	go servePostEventDispatcher(endSignal)

	// wait for a terminating signal from the OS
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	sig := <-signalChan

	logger.InfoData("Received signal, shutting down the service", logging.Data{"sig": sig})
	s.GracefulStop()

	// Wait for the server to stop gracefully
	time.Sleep(time.Second)
}

func main() {
	actualMain()
}
