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

	newfeedProtoV1 "github.com/vantoan19/Petifies/proto/newfeed-service/v1"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/config"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/consumer"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	"github.com/vantoan19/Petifies/server/libs/grpcutils"
	logging "github.com/vantoan19/Petifies/server/libs/logging-config"
	relationshipclient "github.com/vantoan19/Petifies/server/services/grpc-clients/relationship-client"
	cmd "github.com/vantoan19/Petifies/server/services/newfeed-service/cmd"
	postfeedservice "github.com/vantoan19/Petifies/server/services/newfeed-service/internal/application/services/postfeed-service"
	"github.com/vantoan19/Petifies/server/services/newfeed-service/internal/infra/listener"
	postfeedCassandraRepo "github.com/vantoan19/Petifies/server/services/newfeed-service/internal/infra/repositories/post/cassandra"
	userCassandraRepo "github.com/vantoan19/Petifies/server/services/newfeed-service/internal/infra/repositories/user/cassandra"
	endpointsV1 "github.com/vantoan19/Petifies/server/services/newfeed-service/internal/presentation/endpoints/grpc/v1"
	serversV1 "github.com/vantoan19/Petifies/server/services/newfeed-service/internal/presentation/transport/grpc/v1"
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

	postfeedSvc, err := postfeedservice.NewPostFeedService(
		postfeedservice.WithCassandraPostfeedRepository(cmd.DB),
	)
	if err != nil {
		panic(err)
	}
	newfeedEndpoints := endpointsV1.NewNewfeedEndpoints(postfeedSvc)
	newfeedProtoV1.RegisterNewfeedServiceServer(grpcServer, serversV1.NewNewfeedServer(newfeedEndpoints))
	logger.Info("Finished registerServices: SUCCESSFUL")
}

func serveUserConsumer() {
	userRepo, err := userCassandraRepo.NewCassandraUserRepository(cmd.DB)
	if err != nil {
		panic(err)
	}
	listener := listener.NewUserEventListener(userRepo)
	consumerLogger := logging.New("NewFeedService.UserEventConsumer")

	consumer, err := consumer.NewKafkaConsumer(
		config.NewKafkaConsumerConfig(cmd.Conf.Brokers, cmd.Conf.UserEventTopic, cmd.Conf.UserEventConsumerGroup),
		func(ctx context.Context, message *models.KafkaMessage) error {
			var userEvent models.UserEvent
			err := userEvent.Deserialize(message.Value)
			if err != nil {
				return err
			}

			consumerLogger.InfoData("Received a new User Event", logging.Data{"event": userEvent})
			_, err = listener.UserCreated(ctx, userEvent)
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

func servePostConsumer() {
	userRepo, err := userCassandraRepo.NewCassandraUserRepository(cmd.DB)
	if err != nil {
		panic(err)
	}
	postRepo, err := postfeedCassandraRepo.NewCassandraPostRepository(cmd.DB)
	if err != nil {
		panic(err)
	}
	relationshipClient := relationshipclient.New(cmd.RelationshipServiceConn)
	listener := listener.NewPostEventListener(userRepo, postRepo, relationshipClient)
	consumerLogger := logging.New("NewFeedService.PostEventConsumer")

	consumer, err := consumer.NewKafkaConsumer(
		config.NewKafkaConsumerConfig(cmd.Conf.Brokers, cmd.Conf.PostEventTopic, cmd.Conf.PostEventConsumerGroup),
		func(ctx context.Context, message *models.KafkaMessage) error {
			var postEvent models.PostEvent
			err := postEvent.Deserialize(message.Value)
			if err != nil {
				return err
			}

			consumerLogger.InfoData("Received a new Post Event", logging.Data{"event": postEvent})
			_, err = listener.PostCreated(ctx, postEvent)
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
	go serveUserConsumer()
	go servePostConsumer()

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
