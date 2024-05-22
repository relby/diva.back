package api

import (
	"context"
	"errors"

	"github.com/relby/diva.back/internal/convert"
	"github.com/relby/diva.back/internal/model"
	"github.com/relby/diva.back/internal/service"
	"github.com/relby/diva.back/pkg/genproto"
)

func (server *GRPCServer) ListCustomers(ctx context.Context, req *genproto.GetCustomersRequest) (*genproto.GetCustomersResponse, error) {
	customers, err := server.customerService.ListCustomers(ctx, &service.CustomerServiceGetManyCustomersOptions{
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
	})
	if err != nil {
		return nil, err
	}

	customersProto := make([]*genproto.Customer, len(customers))
	for i, customer := range customers {
		customersProto[i] = convert.CustomerFromModelToProto(customer)
	}

	return &genproto.GetCustomersResponse{
		Customers: customersProto,
	}, nil
}

func (server *GRPCServer) UpdateCustomer(ctx context.Context, req *genproto.UpdateCustomerRequest) (*genproto.UpdateCustomerResponse, error) {
	id, err := model.NewCustomerID(req.Id)
	if err != nil {
		return nil, err
	}

	customer, err := server.customerService.GetCustomer(ctx, id)
	if err != nil {
		return nil, err
	}

	var values service.UpdateCustomerValues
	if req.FullName != nil {
		fullName, err := model.NewCustomerFullName(*req.FullName)
		if err != nil {
			return nil, err
		}

		values.FullName = &fullName
	}
	if req.PhoneNumber != nil {
		phoneNumber, err := model.NewCustomerPhoneNumber(*req.PhoneNumber)
		if err != nil {
			return nil, err
		}

		values.PhoneNumber = &phoneNumber
	}
	if req.Discount != nil {
		discount, err := model.NewCustomerDiscount(*req.Discount)
		if err != nil {
			return nil, err
		}

		values.Discount = &discount
	}

	err = server.customerService.UpdateCustomer(ctx, customer, &values)
	if err != nil {
		return nil, err
	}

	return &genproto.UpdateCustomerResponse{
		Customer: convert.CustomerFromModelToProto(customer),
	}, nil
}

func (server *GRPCServer) AddCustomer(ctx context.Context, req *genproto.AddCustomerRequest) (*genproto.AddCustomerResponse, error) {
	fullName, err := model.NewCustomerFullName(req.FullName)
	if err != nil {
		return nil, err
	}

	phoneNumber, err := model.NewCustomerPhoneNumber(req.PhoneNumber)
	if err != nil {
		return nil, err
	}

	discount, err := model.NewCustomerDiscount(req.Discount)
	if err != nil {
		return nil, err
	}

	customer, err := server.customerService.AddCustomer(ctx, fullName, phoneNumber, discount)
	if err != nil {
		return nil, err
	}

	return &genproto.AddCustomerResponse{
		Customer: convert.CustomerFromModelToProto(customer),
	}, nil
}

func (server *GRPCServer) DeleteCustomer(ctx context.Context, req *genproto.DeleteCustomerRequest) (*genproto.DeleteCustomerResponse, error) {
	id, err := model.NewCustomerID(req.Id)
	if err != nil {
		return nil, err
	}

	customer, err := server.customerService.GetCustomer(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = server.customerService.DeleteCustomer(ctx, customer); err != nil {
		return nil, err
	}

	return &genproto.DeleteCustomerResponse{
		Customer: convert.CustomerFromModelToProto(customer),
	}, nil
}

func (server *GRPCServer) ExportCustomersToExcel(ctx context.Context, req *genproto.ExportCustomersToExcelRequest) (*genproto.ExportCustomersToExcelResponse, error) {
	// TODO
	return nil, errors.New("UNIMPLEMENTED")
}
