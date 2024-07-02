package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/relby/diva.back/internal/model"
)

type AccessTokenClaims struct {
	UserID   model.UserID
	UserType model.UserType
}

type RefreshTokenClaims struct {
	ID       model.RefreshTokenID
	UserID   model.UserID
	UserType model.UserType
}

type accessTokenClaims struct {
	jwt.RegisteredClaims
	UserType model.UserType `json:"typ"`
}

type refreshTokenClaims struct {
	jwt.RegisteredClaims
	UserType model.UserType `json:"typ"`
}

func NewAccessToken(claims *AccessTokenClaims, secret string, expireDuration time.Duration) (string, error) {
	now := time.Now()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   claims.UserID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expireDuration)),
		},
		UserType: claims.UserType,
	})

	accessTokenString, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}

func ParseAccessToken(accessTokenString string, secret string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(accessTokenString, &accessTokenClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*accessTokenClaims)
	if !ok {
		return nil, errors.New("failed to parse claims")
	}

	subjectUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, err
	}

	userID, err := model.NewUserID(subjectUUID)
	if err != nil {
		return nil, err
	}

	return &AccessTokenClaims{
		UserID:   userID,
		UserType: claims.UserType,
	}, nil
}

func NewRefreshToken(claims *RefreshTokenClaims, secret string, expireDuration time.Duration) (string, error) {
	now := time.Now()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.UUID(claims.ID).String(),
			Subject:   claims.UserID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expireDuration)),
		},
		UserType: claims.UserType,
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return refreshTokenString, nil
}

func ParseRefreshToken(refreshTokenString string, secret string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(refreshTokenString, &refreshTokenClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*refreshTokenClaims)
	if !ok {
		return nil, errors.New("failed to parse claims")
	}

	idUUID, err := uuid.Parse(claims.ID)
	if err != nil {
		return nil, err
	}

	refreshTokenID, err := model.NewRefreshTokenID(idUUID)
	if err != nil {
		return nil, err
	}

	subjectUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, err
	}

	userID, err := model.NewUserID(subjectUUID)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenClaims{
		ID:       refreshTokenID,
		UserID:   userID,
		UserType: claims.UserType,
	}, nil
}
