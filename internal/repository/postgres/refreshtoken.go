package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/relby/diva.back/internal/convert"
	"github.com/relby/diva.back/internal/domainerrors"
	"github.com/relby/diva.back/internal/model"
	"github.com/relby/diva.back/internal/repository"
	"github.com/relby/diva.back/pkg/gensqlc"
)

var _ repository.RefreshTokenRepository = (*RefreshTokenRepository)(nil)

type RefreshTokenRepository struct {
	queries *gensqlc.Queries
}

func NewRefreshTokenRepository(queries *gensqlc.Queries) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		queries: queries,
	}
}

func (repository *RefreshTokenRepository) GetByID(ctx context.Context, id model.RefreshTokenID) (*model.RefreshToken, error) {
	refreshTokenRow, err := repository.queries.SelectUserRefreshTokenById(ctx, uuid.UUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerrors.NewNotFoundError("refresh token not found")
		}
		return nil, err
	}

	refreshToken, err := convert.RefreshTokenFromRowToModel(refreshTokenRow)
	if err != nil {
		return nil, err
	}

	return refreshToken, nil
}

func (repository *RefreshTokenRepository) Save(ctx context.Context, refreshToken *model.RefreshToken) error {
	var expiresAt pgtype.Timestamptz
	if err := expiresAt.Scan(time.Time(refreshToken.ExpiresAt())); err != nil {
		return err
	}

	if err := repository.queries.UpsertUserRefreshToken(ctx, gensqlc.UpsertUserRefreshTokenParams{
		ID:        uuid.UUID(refreshToken.ID()),
		UserID:    uuid.UUID(refreshToken.UserID()),
		ExpiresAt: expiresAt,
	}); err != nil {
		return err
	}

	return nil
}

func (repository *RefreshTokenRepository) Delete(ctx context.Context, refreshToken *model.RefreshToken) error {
	if err := repository.queries.DeleteRefreshTokenById(ctx, uuid.UUID(refreshToken.ID())); err != nil {
		return err
	}

	return nil
}
