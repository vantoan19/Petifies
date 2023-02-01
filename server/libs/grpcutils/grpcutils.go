package grpcutils

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
)

const timePingInterval = time.Second

func NewTLSServer(myKeyPath string, myCertPath string) (*grpc.Server, error) {
	credential, err := LoadServerTLSCreadentials(myKeyPath, myCertPath)
	if err != nil {
		return nil, err
	}

	serverOptions := []grpc.ServerOption{
		grpc.Creds(credential),
	}

	return grpc.NewServer(serverOptions...), nil
}

func NewTLSClient(serverCAPath string, serverEndpoint string, timeout time.Duration) (*grpc.ClientConn, error) {
	credential, err := LoadTLSCredentialsForClient(serverCAPath)
	if err != nil {
		return nil, err
	}

	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(credential),
	}

	conn, err := grpc.Dial(serverEndpoint, dialOptions...)
	if err != nil {
		return nil, err
	}

	return waitForConn(conn, timeout)
}

func LoadServerTLSCreadentials(serverKeyPath string, serverCertPath string) (credentials.TransportCredentials, error) {
	certificate, err := tls.LoadX509KeyPair(serverCertPath, serverKeyPath)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}

func LoadTLSCredentialsForClient(serverCAPath string) (credentials.TransportCredentials, error) {
	serverCA, err := ioutil.ReadFile(serverCAPath)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(serverCA)
	if !ok {
		return nil, errors.New("LoadClientTLSCredentials: failed to parse server CA as pem")
	}

	config := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(config), nil
}

func waitForConn(conn *grpc.ClientConn, timeout time.Duration) (*grpc.ClientConn, error) {
	waitedTime := time.Second * 0

	for waitedTime < timeout {
		if conn.GetState() != connectivity.Ready {
			time.Sleep(timePingInterval)
			waitedTime += timePingInterval
		} else {
			return conn, nil
		}
	}

	return nil, errors.New("gRPC connection timed out, unable to connect to the server")
}
