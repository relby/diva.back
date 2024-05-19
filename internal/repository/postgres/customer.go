package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/relby/diva.back/internal/convert"
	"github.com/relby/diva.back/internal/model"
	"github.com/relby/diva.back/internal/repository"
	"github.com/relby/diva.back/pkg/gensqlc"
)

var _ repository.CustomerRepository = (*CustomerRepository)(nil)

type CustomerRepository struct {
	queries *gensqlc.Queries
}

func NewCustomerRepository(queries *gensqlc.Queries) *CustomerRepository {
	return &CustomerRepository{
		queries: queries,
	}
}

func (repository *CustomerRepository) FindMany(ctx context.Context, options repository.CustomerRepositoryFindManyOptions) ([]*model.Customer, error) {
	var fullName, phoneNumber pgtype.Text
	if options.FullName != "" {
		fullName.String = options.FullName
		fullName.Valid = true
	}
	if options.PhoneNumber != "" {
		phoneNumber.String = options.PhoneNumber
		phoneNumber.Valid = true
	}
	customersRow, err := repository.queries.SelectCustomers(ctx, gensqlc.SelectCustomersParams{
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	})
	if err != nil {
		return nil, err
	}

	domainCustomers := make([]*model.Customer, len(customersRow))
	for i, row := range customersRow {
		domainCustomers[i], err = convert.CustomerFromRowToModel(row)
		if err != nil {
			return nil, err
		}
	}

	return domainCustomers, nil
}

func (repository *CustomerRepository) FindOneByID(ctx context.Context, id model.CustomerID) (*model.Customer, error) {
	customerRow, err := repository.queries.SelectCustomerById(ctx, int64(id))
	if err != nil {
		return nil, err
	}

	customer, err := convert.CustomerFromRowToModel(customerRow)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (repository *CustomerRepository) Save(ctx context.Context, customer *model.Customer) error {
	// TODO: Find the way to do that more concise
	var id int64
	var err error
	if customer.ID() == 0 {
		id, err = repository.queries.InsertCustomer(ctx, gensqlc.InsertCustomerParams{
			FullName:    string(customer.FullName()),
			PhoneNumber: string(customer.PhoneNumber()),
			Discount:    int16(customer.Discount()),
		})
	} else {
		id, err = repository.queries.UpsertCustomer(ctx, gensqlc.UpsertCustomerParams{
			ID:          int64(customer.ID()),
			FullName:    string(customer.FullName()),
			PhoneNumber: string(customer.PhoneNumber()),
			Discount:    int16(customer.Discount()),
		})
	}
	if err != nil {
		return err
	}

	customerID, err := model.NewCustomerID(id)
	if err != nil {
		return err
	}

	customer.SetID(customerID)

	return nil
}
