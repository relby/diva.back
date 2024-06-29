package repository

import (
	"context"

	"github.com/relby/diva.back/internal/model"
)

type AdminRepository interface {
	GetByID(ctx context.Context, id model.UserID) (*model.Admin, error)
	GetByLogin(ctx context.Context, login model.AdminLogin) (*model.Admin, error)
	Save(ctx context.Context, admin *model.Admin) error
}
