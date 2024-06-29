package repository

import (
	"context"

	"github.com/relby/diva.back/internal/model"
)

type EmployeeRepository interface {
	GetByID(ctx context.Context, id model.UserID) (*model.Employee, error)
	GetByAccessKey(ctx context.Context, accessKey model.EmployeeAccessKey) (*model.Employee, error)
	List(ctx context.Context) ([]*model.Employee, error)
	Save(ctx context.Context, employee *model.Employee) error
	Delete(ctx context.Context, employee *model.Employee) error
}
