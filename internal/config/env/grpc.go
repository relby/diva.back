package env

import (
	"net"
	"os"
)

const grpcHostName = "GRPC_HOST"
const grpcPortName = "GRPC_PORT"

type GRPCConfig struct {
	host string
	port string
}

func NewGRPCConfig() (*GRPCConfig, error) {
	host := os.Getenv(grpcHostName)
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv(grpcPortName)
	if port == "" {
		port = "50051"
	}

	return &GRPCConfig{
		host: host,
		port: port,
	}, nil
}

func (config *GRPCConfig) Host() string {
	return config.host
}

func (config *GRPCConfig) Port() string {
	return config.port
}

func (config *GRPCConfig) Address() string {
	return net.JoinHostPort(config.host, config.port)
}
