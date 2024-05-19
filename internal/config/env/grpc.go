package env

import (
	"net"
	"os"
)

const grpcHostEnvName = "GRPC_HOST"
const grpcPortEnvName = "GRPC_PORT"

type GRPCConfig struct {
	host string
	port string
}

func NewGRPCConfig() (*GRPCConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		host = "127.0.0.1"
	}

	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
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
