package env

import (
	"fmt"
	"os"

	"github.com/relby/diva.back/internal/config"
)

const (
	postgresUserName     = "POSTGRES_USER"
	postgresPasswordName = "POSTGRES_PASSWORD"
	postgresHostName     = "POSTGRES_HOST"
	postgresPortName     = "POSTGRES_PORT"
	postgresDbName       = "POSTGRES_DB"
)

var _ config.PostgresConfig = (*postgresConfig)(nil)

type postgresConfig struct {
	user     string
	password string
	host     string
	port     string
	db       string
}

func NewPostgresConfig() (*postgresConfig, error) {
	user := os.Getenv(postgresUserName)
	if user == "" {
		return nil, notFoundError(postgresUserName)
	}
	password := os.Getenv(postgresPasswordName)
	if password == "" {
		return nil, notFoundError(postgresPasswordName)
	}
	host := os.Getenv(postgresHostName)
	if host == "" {
		return nil, notFoundError(postgresHostName)
	}
	port := os.Getenv(postgresPortName)
	if port == "" {
		return nil, notFoundError(postgresPortName)
	}
	db := os.Getenv(postgresDbName)
	if db == "" {
		return nil, notFoundError(postgresDbName)
	}
	return &postgresConfig{
		user:     user,
		password: password,
		host:     host,
		port:     port,
		db:       db,
	}, nil
}

func (pc postgresConfig) User() string {
	return pc.user
}

func (pc postgresConfig) Password() string {
	return pc.password
}

func (pc postgresConfig) Host() string {
	return pc.host
}
func (pc postgresConfig) Port() string {
	return pc.port
}

func (pc postgresConfig) DB() string {
	return pc.db
}
func (pc postgresConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", pc.user, pc.password, pc.host, pc.port, pc.db)
}
