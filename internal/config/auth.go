package config

import "time"

type AuthConfig interface {
	AccessTokenSecret() string
	AccessTokenExpireDuration() time.Duration
	RefreshTokenSecret() string
	RefreshTokenExpireDuration() time.Duration
}
