package repository

import (
	"context"

	"github.com/relby/diva.back/internal/model"
)

type AdminRepository interface {
	Save(ctx context.Context, admin *model.Admin) error
}
