package env

import (
	"os"

	"github.com/relby/diva.back/internal/config"
)

const (
	pgDsnName = "PG_DSN"
)

var _ config.PostgresConfig = (*postgresConfig)(nil)

type postgresConfig struct {
	dsn string
}

func NewPostgresConfig() (*postgresConfig, error) {
	dsn := os.Getenv(pgDsnName)
	if dsn == "" {
		return nil, notFoundError(pgDsnName)
	}
	return &postgresConfig{
		dsn: dsn,
	}, nil
}

func (x postgresConfig) DSN() string {
	return x.dsn
}
