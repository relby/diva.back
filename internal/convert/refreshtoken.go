package convert

import (
	"github.com/relby/diva.back/internal/model"
	"github.com/relby/diva.back/pkg/gensqlc"
)

func RefreshTokenFromRowToModel(row *gensqlc.UserRefreshToken) (*model.RefreshToken, error) {
	id, err := model.NewRefreshTokenID(row.ID)
	if err != nil {
		return nil, err
	}

	userID, err := model.NewUserID(row.UserID)
	if err != nil {
		return nil, err
	}

	expiresAt, err := model.NewRefreshTokenExpiresAt(row.ExpiresAt.Time)
	if err != nil {
		return nil, err
	}

	refreshToken, err := model.NewRefreshToken(id, userID, expiresAt)
	if err != nil {
		return nil, err
	}

	return refreshToken, nil
}
