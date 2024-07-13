package config

type PostgresConfig interface {
	User() string
	Password() string
	Host() string
	Port() string
	DB() string
	DSN() string
}
