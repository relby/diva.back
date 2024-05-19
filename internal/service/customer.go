package service

import (
	"context"

	"github.com/relby/diva.back/internal/model"
	"github.com/relby/diva.back/internal/repository"
)

type CustomerService struct {
	customerRepository repository.CustomerRepository
}

func NewCustomerService(repository repository.CustomerRepository) *CustomerService {
	return &CustomerService{
		customerRepository: repository,
	}
}

type CustomerServiceGetManyCustomersOptions struct {
	FullName    string
	PhoneNumber string
}

func (service CustomerService) GetManyCustomers(ctx context.Context, options CustomerServiceGetManyCustomersOptions) ([]*model.Customer, error) {
	customers, err := service.customerRepository.FindMany(ctx, repository.CustomerRepositoryFindManyOptions(options))
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (service CustomerService) SetCustomerDiscountByID(ctx context.Context, id model.CustomerID, discount model.CustomerDiscount) (*model.Customer, error) {
	customer, err := service.customerRepository.FindOneByID(ctx, id)
	if err != nil {
		return nil, err
	}

	customer.SetDiscount(discount)

	err = service.customerRepository.Save(ctx, customer)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (service CustomerService) AddCustomer(ctx context.Context, fullName model.CustomerFullName, phoneNumber model.CustomerPhoneNumber, discount model.CustomerDiscount) (*model.Customer, error) {
	customer, err := model.NewCustomer(fullName, phoneNumber, discount)
	if err != nil {
		return nil, err
	}

	if err = service.customerRepository.Save(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil
}
