package convert

import (
	"github.com/relby/diva.back/internal/model"
	"github.com/relby/diva.back/pkg/gensqlc"
)

func AdminFromRowToModel(userRow gensqlc.User, adminRow gensqlc.Admin) (*model.Admin, error) {
	id, fullName, err := userFromRowToValueObjects(userRow)
	if err != nil {
		return nil, err
	}

	login, err := model.NewAdminLogin(adminRow.Login)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := model.NewAdminHashedPassword(adminRow.HashedPassword)
	if err != nil {
		return nil, err
	}

	admin, err := model.NewAdmin(id, fullName, login, hashedPassword)
	if err != nil {
		return nil, err
	}

	return admin, nil
}
