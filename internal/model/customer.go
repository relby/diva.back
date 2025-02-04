package model

import (
	"github.com/relby/diva.back/internal/domainerrors"
	"golang.org/x/exp/constraints"
)

type CustomerID uint64

func NewCustomerID[T constraints.Integer](value T) (CustomerID, error) {
	if value < 1 {
		return 0, domainerrors.NewValidationError("customer id must be >= 1")
	}

	return CustomerID(value), nil
}

type CustomerFullName string

func NewCustomerFullName(value string) (CustomerFullName, error) {
	if value == "" {
		return "", domainerrors.NewValidationError("customer full name is empty")
	}
	return CustomerFullName(value), nil
}

type CustomerPhoneNumber string

func NewCustomerPhoneNumber(value string) (CustomerPhoneNumber, error) {
	// TODO: maybe introduce validation for customer phone number
	return CustomerPhoneNumber(value), nil
}

type CustomerDiscount uint8

func NewCustomerDiscount[T constraints.Integer](value T) (CustomerDiscount, error) {
	if value < 0 || value > 100 {
		return 0, domainerrors.NewValidationError("customer discount must in range 0-100")
	}

	return CustomerDiscount(value), nil
}

type Customer struct {
	id          CustomerID
	fullName    CustomerFullName
	phoneNumber CustomerPhoneNumber
	discount    CustomerDiscount
}

func NewCustomer(fullName CustomerFullName, phoneNumber CustomerPhoneNumber, discount CustomerDiscount) (*Customer, error) {
	return &Customer{
		id:          CustomerID(0),
		fullName:    fullName,
		phoneNumber: phoneNumber,
		discount:    discount,
	}, nil
}

func NewCustomerWithId(id CustomerID, fullName CustomerFullName, phoneNumber CustomerPhoneNumber, discount CustomerDiscount) (*Customer, error) {
	customer, err := NewCustomer(fullName, phoneNumber, discount)
	if err != nil {
		return nil, err
	}

	customer.SetID(id)

	return customer, nil
}

func (customer *Customer) ID() CustomerID {
	return customer.id
}
func (customer *Customer) FullName() CustomerFullName {
	return customer.fullName
}
func (customer *Customer) PhoneNumber() CustomerPhoneNumber {
	return customer.phoneNumber
}
func (customer *Customer) Discount() CustomerDiscount {
	return customer.discount
}

func (customer *Customer) SetID(id CustomerID) {
	customer.id = id
}

func (customer *Customer) SetFullName(fullName CustomerFullName) {
	customer.fullName = fullName
}

func (customer *Customer) SetPhoneNumber(phoneNumber CustomerPhoneNumber) {
	customer.phoneNumber = phoneNumber
}

func (customer *Customer) SetDiscount(discount CustomerDiscount) {
	customer.discount = discount
}
