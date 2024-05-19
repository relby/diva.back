package env

import (
	"errors"
	"os"

	"github.com/relby/diva.back/internal/config"
)

const (
	dsnEnvName = "PG_DSN"
)

var _ config.PostgresConfig = (*postgresConfig)(nil)

type postgresConfig struct {
	dsn string
}

func NewPostgresConfig() (*postgresConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("postgres dsn not found")
	}
	return &postgresConfig{
		dsn: dsn,
	}, nil
}

func (x postgresConfig) DSN() string {
	return x.dsn
}
