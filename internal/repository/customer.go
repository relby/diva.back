package repository

import (
	"context"

	"github.com/relby/diva.back/internal/model"
)

type CustomerRepositoryFindManyOptions struct {
	FullName    *string
	PhoneNumber *string
}

type CustomerRepository interface {
	List(ctx context.Context, options *CustomerRepositoryFindManyOptions) ([]*model.Customer, error)
	GetByID(ctx context.Context, id model.CustomerID) (*model.Customer, error)
	Save(ctx context.Context, customer *model.Customer) error
	Delete(ctx context.Context, customer *model.Customer) error
}
