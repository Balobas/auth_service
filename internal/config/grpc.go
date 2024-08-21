package config

import (
	"net"
	"os"

	"github.com/pkg/errors"
)

type ConfigGRPC struct {
	host string
	port string
}

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"
)

func NewConfigGRPC() (*ConfigGRPC, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.Errorf("%s is empty", grpcHostEnvName)
	}
	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, errors.Errorf("%s is empty", grpcPortEnvName)
	}

	return &ConfigGRPC{
		host: host,
		port: port,
	}, nil
}

func (c *ConfigGRPC) Address() string {
	return net.JoinHostPort(c.host, c.port)
}
