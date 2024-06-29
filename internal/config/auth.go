package config

import "time"

type AuthConfig interface {
	JWTSecret() string
	JWTExpireDuration() time.Duration
}
