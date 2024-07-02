package env

import (
	"fmt"
	"os"
	"time"
)

const accessTokenSecretName = "ACCESS_TOKEN_SECRET"
const accessTokenExpireDurationName = "ACCESS_TOKEN_EXPIRE_DURATION"
const refreshTokenSecretName = "REFRESH_TOKEN_SECRET"
const refreshTokenExpireDurationName = "REFRESH_TOKEN_EXPIRE_DURATION"

type AuthConfig struct {
	accessTokenSecret          string
	accessTokenExpireDuration  time.Duration
	refreshTokenSecret         string
	refreshTokenExpireDuration time.Duration
}

func NewAuthConfig() (*AuthConfig, error) {
	accessTokenSecret := os.Getenv(accessTokenSecretName)
	if accessTokenSecret == "" {
		return nil, notFoundError(accessTokenSecretName)
	}

	accessTokenExpireDurationString := os.Getenv(accessTokenExpireDurationName)
	if accessTokenExpireDurationString == "" {
		return nil, notFoundError(accessTokenExpireDurationName)
	}

	accessTokenExpireDuration, err := time.ParseDuration(accessTokenExpireDurationString)
	if err != nil {
		return nil, fmt.Errorf("`%s` has invalid format: %w", accessTokenExpireDurationName, err)
	}

	refreshTokenSecret := os.Getenv(refreshTokenSecretName)
	if refreshTokenSecret == "" {
		return nil, notFoundError(refreshTokenSecretName)
	}

	refreshTokenExpireDurationString := os.Getenv(refreshTokenExpireDurationName)
	if refreshTokenExpireDurationString == "" {
		return nil, notFoundError(refreshTokenExpireDurationName)
	}

	refreshTokenExpireDuration, err := time.ParseDuration(refreshTokenExpireDurationString)
	if err != nil {
		return nil, fmt.Errorf("`%s` has invalid format: %w", refreshTokenExpireDurationName, err)
	}

	return &AuthConfig{
		accessTokenSecret:          accessTokenSecret,
		accessTokenExpireDuration:  accessTokenExpireDuration,
		refreshTokenSecret:         refreshTokenSecret,
		refreshTokenExpireDuration: refreshTokenExpireDuration,
	}, nil
}

func (config *AuthConfig) AccessTokenSecret() string {
	return config.accessTokenSecret
}

func (config *AuthConfig) AccessTokenExpireDuration() time.Duration {
	return config.accessTokenExpireDuration
}

func (config *AuthConfig) RefreshTokenSecret() string {
	return config.refreshTokenSecret
}

func (config *AuthConfig) RefreshTokenExpireDuration() time.Duration {
	return config.refreshTokenExpireDuration
}
