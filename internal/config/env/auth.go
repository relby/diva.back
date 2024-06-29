package env

import (
	"fmt"
	"os"
	"time"
)

const jwtSecretName = "JWT_SECRET"
const jwtExpireDurationName = "JWT_EXPIRE_DURATION"

type AuthConfig struct {
	jwtSecret         string
	jwtExpireDuration time.Duration
}

func NewAuthConfig() (*AuthConfig, error) {
	jwtSecret := os.Getenv(jwtSecretName)
	if jwtSecret == "" {
		return nil, notFoundError(jwtSecretName)
	}

	jwtExpireDurationString := os.Getenv(jwtExpireDurationName)
	if jwtExpireDurationString == "" {
		return nil, notFoundError(jwtExpireDurationName)
	}

	jwtExpireDuration, err := time.ParseDuration(jwtExpireDurationString)
	if err != nil {
		return nil, fmt.Errorf("`%s` has invalid format: %w", jwtExpireDurationName, err)
	}

	return &AuthConfig{
		jwtSecret:         jwtSecret,
		jwtExpireDuration: jwtExpireDuration,
	}, nil
}

func (config *AuthConfig) JWTSecret() string {
	return config.jwtSecret
}

func (config *AuthConfig) JWTExpireDuration() time.Duration {
	return config.jwtExpireDuration
}
