package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	locationProtoV1 "github.com/vantoan19/Petifies/proto/location-service/v1"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/config"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/consumer"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	"github.com/vantoan19/Petifies/server/libs/grpcutils"
	logging "github.com/vantoan19/Petifies/server/libs/logging-config"
	cmd "github.com/vantoan19/Petifies/server/services/location-service/cmd"
	locationservice "github.com/vantoan19/Petifies/server/services/location-service/internal/application/services/location"
	petifieskafkalistener "github.com/vantoan19/Petifies/server/services/location-service/internal/infra/listeners/petifies/kafka"
	locationmongo "github.com/vantoan19/Petifies/server/services/location-service/internal/infra/repositories/location/mongo"
	endpointsV1 "github.com/vantoan19/Petifies/server/services/location-service/internal/presentation/endpoints/grpc/v1"
	serversV1 "github.com/vantoan19/Petifies/server/services/location-service/internal/presentation/transport/grpc/v1"
)

var logger = logging.New("RelationshipService.Cmd.Grpc")

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

	locationSvc, err := locationservice.NewLocationService(
		locationservice.WithMongoLocationRepository(cmd.MongoClient),
	)
	if err != nil {
		panic(err)
	}

	locationEndpoints := endpointsV1.NewLocationEndpoints(locationSvc)
	locationProtoV1.RegisterLocationServiceServer(grpcServer, serversV1.NewLocationServer(locationEndpoints))
	logger.Info("Finished registerServices: SUCCESSFUL")
}

func servePetifiesConsumer() {
	locationRepo := locationmongo.New(cmd.MongoClient)
	listener := petifieskafkalistener.NewKafkaPetifiesEventListener(locationRepo)
	consumerLogger := logging.New("LocationService.PetifiesEventConsumer")

	consumer, err := consumer.NewKafkaConsumer(
		config.NewKafkaConsumerConfig(cmd.Conf.Brokers, cmd.Conf.PetifiesEventTopic, cmd.Conf.PetifiesEventConsumerGroup),
		func(ctx context.Context, message *models.KafkaMessage) error {
			var petifiesEvent models.PetifiesEvent
			err := petifiesEvent.Deserialize(message.Value)
			if err != nil {
				return err
			}

			consumerLogger.InfoData("Received a new Petifies Event", logging.Data{"event": petifiesEvent})
			err = listener.Receive(ctx, petifiesEvent)
			if err != nil {
				return err
			}

			return nil
		},
	)
	if err != nil {
		panic(err)
	}
	consumer.Consume()
}

func actualMain() {
	logger.Info("Relationship Service starting up")
	cmd.Initialize()
	s, err := setupGRPC()
	if err != nil {
		panic(err)
	}

	go serveGRPC(s)
	go servePetifiesConsumer()

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
