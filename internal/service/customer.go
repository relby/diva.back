package service

import (
	"context"

	"github.com/relby/diva.back/internal/model"
	"github.com/relby/diva.back/internal/repository"
)

type CustomerService struct {
	customerRepository repository.CustomerRepository
}

func NewCustomerService(customerRepository repository.CustomerRepository) *CustomerService {
	return &CustomerService{
		customerRepository: customerRepository,
	}
}

type CustomerServiceGetManyCustomersOptions struct {
	FullName    *string
	PhoneNumber *string
}

func (service *CustomerService) GetCustomer(ctx context.Context, id model.CustomerID) (*model.Customer, error) {
	customer, err := service.customerRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (service *CustomerService) ListCustomers(ctx context.Context, options *CustomerServiceGetManyCustomersOptions) ([]*model.Customer, error) {
	customers, err := service.customerRepository.List(ctx, (*repository.CustomerRepositoryFindManyOptions)(options))
	if err != nil {
		return nil, err
	}

	return customers, nil
}

type UpdateCustomerValues struct {
	FullName    *model.CustomerFullName
	PhoneNumber *model.CustomerPhoneNumber
	Discount    *model.CustomerDiscount
}

func (service *CustomerService) UpdateCustomer(ctx context.Context, customer *model.Customer, values *UpdateCustomerValues) error {
	if values.FullName != nil {
		customer.SetFullName(*values.FullName)
	}
	if values.PhoneNumber != nil {
		customer.SetPhoneNumber(*values.PhoneNumber)
	}
	if values.Discount != nil {
		customer.SetDiscount(*values.Discount)
	}

	if err := service.customerRepository.Save(ctx, customer); err != nil {
		return err
	}

	return nil
}

func (service *CustomerService) AddCustomer(ctx context.Context, fullName model.CustomerFullName, phoneNumber model.CustomerPhoneNumber, discount model.CustomerDiscount) (*model.Customer, error) {
	customer, err := model.NewCustomer(fullName, phoneNumber, discount)
	if err != nil {
		return nil, err
	}

	if err = service.customerRepository.Save(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil
}

func (service *CustomerService) DeleteCustomer(ctx context.Context, customer *model.Customer) error {
	if err := service.customerRepository.Delete(ctx, customer); err != nil {
		return err
	}

	return nil
}
