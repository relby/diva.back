package config

type PostgresConfig interface {
	DSN() string
}
