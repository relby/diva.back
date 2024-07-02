package repository

import (
	"context"

	"github.com/relby/diva.back/internal/model"
)

type RefreshTokenRepository interface {
	GetByID(ctx context.Context, id model.RefreshTokenID) (*model.RefreshToken, error)
	Save(ctx context.Context, refreshToken *model.RefreshToken) error
	Delete(ctx context.Context, refreshToken *model.RefreshToken) error
}
