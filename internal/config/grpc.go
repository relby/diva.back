package config

type GRPCConfig interface {
	Host() string
	Port() string
	Address() string
}
