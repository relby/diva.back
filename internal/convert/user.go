package convert

import (
	"github.com/relby/diva.back/internal/model"
	"github.com/relby/diva.back/pkg/gensqlc"
)

func userFromRowToValueObjects(userRow gensqlc.User) (model.UserID, model.UserFullName, error) {
	id, err := model.NewUserID(userRow.ID)
	if err != nil {
		return model.UserID{}, "", err
	}

	fullName, err := model.NewUserFullName(userRow.FullName)
	if err != nil {
		return model.UserID{}, "", err
	}

	return id, fullName, nil
}
