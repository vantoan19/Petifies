package grpcutils

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"time"

	logging "github.com/vantoan19/Petifies/server/libs/logging-config"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
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
	logger.Info("Creating TLS Grpc server")

	credential, err := LoadGrpcServerTLSCreadentials(myKeyPath, myCertPath)
	if err != nil {
		logger.ErrorData("Failed to TLS Grpc server", logging.Data{"error": err.Error()})
		return nil, err
	}

	serverOptions := []grpc.ServerOption{
		grpc.Creds(credential),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(interceptors.UnaryInterceptors...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(interceptors.StreamInterceptors...)),
	}

	logger.Info("Created TLS Grpc server successfully")
	return grpc.NewServer(serverOptions...), nil
}

func NewTLSGrpcClient(serverCAPath string, serverEndpoint string, timeout time.Duration, interceptors ClientInterceptors) (*grpc.ClientConn, error) {
	logger.Info("Creating TLS Grpc client")

	credential, err := LoadTLSCredentialsForGrpcClient(serverCAPath)
	if err != nil {
		logger.ErrorData("Failed to TLS Grpc client", logging.Data{"error": err.Error()})
		return nil, err
	}

	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(credential),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(interceptors.UnaryInterceptors...)),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(interceptors.StreamInterceptors...)),
	}

	conn, err := grpc.Dial(serverEndpoint, dialOptions...)
	if err != nil {
		logger.ErrorData("Failed to TLS Grpc client", logging.Data{"error": err.Error()})
		return nil, err
	}

	return waitForConn(conn, timeout)
}

func NewInsecureGrpcServer(interceptors ServerInterceptors) (*grpc.Server, error) {
	logger.Info("Creating insecure Grpc server")

	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(interceptors.UnaryInterceptors...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(interceptors.StreamInterceptors...)),
	}

	logger.Info("Created insecure Grpc server successfully")
	return grpc.NewServer(serverOptions...), nil
}

func NewInsecureGrpcClient(serverEndpoint string, timeout time.Duration, interceptors ClientInterceptors) (*grpc.ClientConn, error) {
	logger.Info("Creating insecure Grpc client")

	dialOptions := []grpc.DialOption{
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(interceptors.UnaryInterceptors...)),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(interceptors.StreamInterceptors...)),
	}

	conn, err := grpc.Dial(serverEndpoint, dialOptions...)
	if err != nil {
		logger.ErrorData("Failed to insecure Grpc client", logging.Data{"error": err.Error()})
		return nil, err
	}

	return waitForConn(conn, timeout)
}

func LoadGrpcServerTLSCreadentials(serverKeyPath string, serverCertPath string) (credentials.TransportCredentials, error) {
	logger.Info("Loading server creadentials")
	certificate, err := tls.LoadX509KeyPair(serverCertPath, serverKeyPath)
	if err != nil {
		logger.ErrorData("Failed to load certificate", logging.Data{"error": err.Error()})
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.NoClientCert,
	}

	logger.Info("Loaded server creadentials successfully")
	return credentials.NewTLS(config), nil
}

func LoadTLSCredentialsForGrpcClient(serverCAPath string) (credentials.TransportCredentials, error) {
	logger.Info("Loading creadentials for client")
	serverCA, err := ioutil.ReadFile(serverCAPath)
	if err != nil {
		logger.ErrorData("Failed to read CA", logging.Data{"error": err.Error()})
		return nil, err
	}

	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(serverCA)
	if !ok {
		logger.Error("Failed to parse server CA as pem")
		return nil, errors.New("LoadClientTLSCredentials: failed to parse server CA as pem")
	}

	config := &tls.Config{
		RootCAs: certPool,
	}

	logger.Info("Loaded creadentials for client successfully")
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

	logger.Error("Creating Grpc client failed, timed out")
	return nil, errors.New("gRPC connection timed out, unable to connect to the server")
}
