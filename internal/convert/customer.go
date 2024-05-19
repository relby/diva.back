package convert

import (
	"github.com/relby/diva.back/internal/model"
	"github.com/relby/diva.back/pkg/genproto"
	"github.com/relby/diva.back/pkg/gensqlc"
)

func CustomerFromRowToModel(row *gensqlc.Customer) (*model.Customer, error) {
	id, err := model.NewCustomerID(row.ID)
	if err != nil {
		return nil, err
	}
	fullName, err := model.NewCustomerFullName(row.FullName)
	if err != nil {
		return nil, err
	}
	phoneNumber, err := model.NewCustomerPhoneNumber(row.PhoneNumber)
	if err != nil {
		return nil, err
	}
	discount, err := model.NewCustomerDiscount(row.Discount)
	if err != nil {
		return nil, err
	}

	customer, err := model.NewCustomerWithId(
		id,
		fullName,
		phoneNumber,
		discount,
	)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func CustomerFromModelToProto(customer *model.Customer) *genproto.Customer {
	return &genproto.Customer{
		Id:          uint64(customer.ID()),
		FullName:    string(customer.FullName()),
		PhoneNumber: string(customer.PhoneNumber()),
		Discount:    uint32(customer.Discount()),
	}
}
