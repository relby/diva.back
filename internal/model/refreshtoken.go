package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/relby/diva.back/internal/domainerrors"
)

type RefreshTokenID uuid.UUID

func NewRefreshTokenID(value uuid.UUID) (RefreshTokenID, error) {
	if value == uuid.Nil {
		return RefreshTokenID(uuid.Nil), domainerrors.NewValidationError("user refresh token id can't be nil")
	}

	return RefreshTokenID(value), nil
}

type RefreshTokenExpiresAt time.Time

func NewRefreshTokenExpiresAt(value time.Time) (RefreshTokenExpiresAt, error) {
	return RefreshTokenExpiresAt(value), nil
}

type RefreshToken struct {
	id        RefreshTokenID
	userId    UserID
	expiresAt RefreshTokenExpiresAt
}

func NewRefreshToken(id RefreshTokenID, userId UserID, expiresAt RefreshTokenExpiresAt) (*RefreshToken, error) {
	return &RefreshToken{
		id:        id,
		userId:    userId,
		expiresAt: expiresAt,
	}, nil
}

func (urt *RefreshToken) ID() RefreshTokenID {
	return urt.id
}

func (urt *RefreshToken) UserID() UserID {
	return urt.userId
}

func (urt *RefreshToken) ExpiresAt() RefreshTokenExpiresAt {
	return urt.expiresAt
}
