package grpcutils

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
	"time"

	logging "github.com/vantoan19/Petifies/server/libs/logging-config"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var logger = logging.New("Libs.GRPCUtils")

const timePingInterval = time.Second

type ServerInterceptors struct {
	UnaryInterceptors  []grpc.UnaryServerInterceptor
	StreamInterceptors []grpc.StreamServerInterceptor
}

type ClientInterceptors struct {
	UnaryInterceptors  []grpc.UnaryClientInterceptor
	StreamInterceptors []grpc.StreamClientInterceptor
}

func NewTLSGrpcServer(myKeyPath string, myCertPath string, interceptors ServerInterceptors) (*grpc.Server, error) {
	logger.Info("Start NewTLSGrpcServer")

	credential, err := LoadGrpcServerTLSCreadentials(myKeyPath, myCertPath)
	if err != nil {
		logger.ErrorData("Finised NewTLSGrpcServer: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	serverOptions := []grpc.ServerOption{
		grpc.Creds(credential),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(interceptors.UnaryInterceptors...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(interceptors.StreamInterceptors...)),
	}

	logger.Info("Finished NewTLSGrpcServer: SUCCESSFUL")
	return grpc.NewServer(serverOptions...), nil
}

func NewTLSGrpcClient(serverCAPath string, serverEndpoint string, retries int, interceptors ClientInterceptors) (*grpc.ClientConn, error) {
	logger.Info("Start NewTLSGrpcClient")

	credential, err := LoadTLSCredentialsForGrpcClient(serverCAPath)
	if err != nil {
		logger.ErrorData("Finished NewTLSGrpcClient: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(credential),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(interceptors.UnaryInterceptors...)),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(interceptors.StreamInterceptors...)),
	}

	for i := 1; i <= retries; i++ {
		conn, err := grpc.Dial(serverEndpoint, dialOptions...)
		if err != nil {
			logger.WarningData("Executing NewTLSGrpcClient: failed to dial the server, retrying", logging.Data{"error": err.Error()})
		}

		conn, err = waitForConn(conn, time.Second*30)
		if err != nil {
			logger.WarningData("Executing NewTLSGrpcClient: failed to dial the server, retrying", logging.Data{"error": err.Error()})
		} else {
			logger.Info("Finished NewTLSGrpcClient: SUCCESSFUL")
			return conn, nil
		}
	}

	logger.Error("Finished NewTLSGrpcClient: FAILED")
	return nil, errors.New("unable to connect to the GRPC server")
}

func NewInsecureGrpcServer(interceptors ServerInterceptors) (*grpc.Server, error) {
	logger.Info("Start NewInsecureGrpcServer")

	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(interceptors.UnaryInterceptors...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(interceptors.StreamInterceptors...)),
	}

	logger.Info("Finished NewInsecureGrpcServer: SUCCESSFUL")
	return grpc.NewServer(serverOptions...), nil
}

func NewInsecureGrpcClient(serverEndpoint string, retries int, interceptors ClientInterceptors) (*grpc.ClientConn, error) {
	logger.Info("Start NewInsecureGrpcClient")
	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(interceptors.UnaryInterceptors...)),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(interceptors.StreamInterceptors...)),
	}

	for i := 1; i <= retries; i++ {
		conn, err := grpc.Dial(serverEndpoint, dialOptions...)
		if err != nil {
			logger.WarningData("Executing NewInsecureGrpcClient: failed to dial the server, retrying", logging.Data{"error": err.Error()})
		}

		conn, err = waitForConn(conn, time.Second*5)
		if err != nil {
			logger.WarningData("Executing NewInsecureGrpcClient: failed to dial the server, retrying", logging.Data{"error": err.Error()})
		} else {
			logger.Info("Finished NewInsecureGrpcClient: SUCCESSFUL")
			return conn, nil
		}
	}

	logger.Error("Finished NewInsecureGrpcClient: FAILED")
	return nil, errors.New("unable to connect to the GRPC server")
}

func LoadGrpcServerTLSCreadentials(serverKeyPath string, serverCertPath string) (credentials.TransportCredentials, error) {
	logger.Info("Start LoadGrpcServerTLSCreadentials")
	certificate, err := tls.LoadX509KeyPair(serverCertPath, serverKeyPath)
	if err != nil {
		logger.ErrorData("Finished LoadGrpcServerTLSCreadentials: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.NoClientCert,
	}

	logger.Info("Finished LoadGrpcServerTLSCreadentials: SUCCESSFUL")
	return credentials.NewTLS(config), nil
}

func LoadTLSCredentialsForGrpcClient(serverCAPath string) (credentials.TransportCredentials, error) {
	logger.Info("Start LoadTLSCredentialsForGrpcClient")
	serverCA, err := os.ReadFile(serverCAPath)
	if err != nil {
		logger.ErrorData("Finished LoadTLSCredentialsForGrpcClient: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(serverCA)
	if !ok {
		logger.ErrorData("Finished LoadTLSCredentialsForGrpcClient: FAILED", logging.Data{"error": err.Error()})
		return nil, errors.New("failed to parse server CA as pem")
	}

	config := &tls.Config{
		RootCAs: certPool,
	}

	logger.Info("Finished LoadTLSCredentialsForGrpcClient: SUCCESSFUL")
	return credentials.NewTLS(config), nil
}

func waitForConn(conn *grpc.ClientConn, timeout time.Duration) (*grpc.ClientConn, error) {
	logger.Info("Waiting for Grpc connection to be established")
	waitedTime := time.Second * 0

	for waitedTime < timeout {
		if conn.GetState() != connectivity.Ready {
			time.Sleep(timePingInterval)
			waitedTime += timePingInterval
		} else {
			logger.Info("Connection established successfully")
			return conn, nil
		}
	}

	logger.Warning("Creating Grpc client failed, timed out")
	return nil, errors.New("gRPC connection timed out, unable to connect to the server")
}
