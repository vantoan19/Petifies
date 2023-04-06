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

	petifiesProtoV1 "github.com/vantoan19/Petifies/proto/petifies-service/v1"
	outbox_dispatcher "github.com/vantoan19/Petifies/server/infrastructure/outbox/dispatcher"
	"github.com/vantoan19/Petifies/server/libs/grpcutils"
	logging "github.com/vantoan19/Petifies/server/libs/logging-config"
	cmd "github.com/vantoan19/Petifies/server/services/petifies-service/cmd"
	petifiesproposalservice "github.com/vantoan19/Petifies/server/services/petifies-service/internal/application/services/petifies-proposal-service"
	reviewservice "github.com/vantoan19/Petifies/server/services/petifies-service/internal/application/services/petifies-review-service"
	petifiesservice "github.com/vantoan19/Petifies/server/services/petifies-service/internal/application/services/petifies-service"
	petifiessessionservice "github.com/vantoan19/Petifies/server/services/petifies-service/internal/application/services/petifies-session-service"
	petifieseventmongo "github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/repositories/petifies_event/mongo"
	endpointsV1 "github.com/vantoan19/Petifies/server/services/petifies-service/internal/presentation/endpoints/grpc/v1"
	serversV1 "github.com/vantoan19/Petifies/server/services/petifies-service/internal/presentation/transport/grpc/v1"
)

var logger = logging.New("PestifiesService.Cmd.Grpc")

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

	petifiesEventRepo, err := petifieseventmongo.New(cmd.MongoClient)
	if err != nil {
		panic(err)
	}
	// Services
	petifiesSvc, err := petifiesservice.NewPetifiesService(
		petifiesservice.WithMongoPetifiesRepository(cmd.MongoClient),
		petifiesservice.WithMongoPetifiesProposalRepository(cmd.MongoClient),
		petifiesservice.WithMongoPetifiesSessionRepository(cmd.MongoClient),
		petifiesservice.WithMongoReviewRepository(cmd.MongoClient),
		petifiesservice.WithKafkaPetifiesEventPublisher(&cmd.PetifiesProducer, petifiesEventRepo),
	)
	if err != nil {
		panic(err)
	}

	petifiesSessionSvc, err := petifiessessionservice.NewPetifiesSessionService(
		petifiessessionservice.WithMongoPetifiesRepository(cmd.MongoClient),
		petifiessessionservice.WithMongoPetifiesProposalRepository(cmd.MongoClient),
		petifiessessionservice.WithMongoPetifiesSessionRepository(cmd.MongoClient),
		petifiessessionservice.WithMongoReviewRepository(cmd.MongoClient),
		petifiessessionservice.WithKafkaPetifiesSessionEventPublisher(&cmd.PetifiesProducer, petifiesEventRepo),
	)
	if err != nil {
		panic(err)
	}

	petifiesProposalSvc, err := petifiesproposalservice.NewPetifiesProposalService(
		petifiesproposalservice.WithMongoPetifiesRepository(cmd.MongoClient),
		petifiesproposalservice.WithMongoPetifiesProposalRepository(cmd.MongoClient),
		petifiesproposalservice.WithMongoPetifiesSessionRepository(cmd.MongoClient),
		petifiesproposalservice.WithMongoReviewRepository(cmd.MongoClient),
		petifiesproposalservice.WithKafkaPetifiesProposalPublisher(&cmd.PetifiesProducer, petifiesEventRepo),
	)
	if err != nil {
		panic(err)
	}

	reviewSvc, err := reviewservice.NewReviewService(
		reviewservice.WithMongoPetifiesRepository(cmd.MongoClient),
		reviewservice.WithMongoPetifiesProposalRepository(cmd.MongoClient),
		reviewservice.WithMongoPetifiesSessionRepository(cmd.MongoClient),
		reviewservice.WithMongoReviewRepository(cmd.MongoClient),
		reviewservice.WithKafkaReviewPublisher(&cmd.PetifiesProducer, petifiesEventRepo),
	)
	if err != nil {
		panic(err)
	}

	petifiesEndpointsV1 := endpointsV1.NewPetifiesEndpoints(petifiesSvc)
	petifiesSessionEndpointsV1 := endpointsV1.NewPetifiesSessionEndpoints(petifiesSessionSvc)
	petifiesProposalEndpointsV1 := endpointsV1.NewPetifiesProposalEndpoints(petifiesProposalSvc)
	reviewEndpointsV1 := endpointsV1.NewReviewEndpoints(reviewSvc)

	petifiesProtoV1.RegisterPetifiesServiceServer(grpcServer, serversV1.NewPetifiesServer(
		petifiesEndpointsV1, petifiesSessionEndpointsV1, petifiesProposalEndpointsV1, reviewEndpointsV1))

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
	petifiesEventRepo, err := petifieseventmongo.New(cmd.MongoClient)
	if err != nil {
		panic(err)
	}
	errsChan := make(chan error)
	dispatcher := outbox_dispatcher.NewDispatcher(petifiesEventRepo, cmd.PetifiesProducer, settings, *logging.New("PetifiesService.Dispatcher"))
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
