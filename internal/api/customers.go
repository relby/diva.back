package api

import (
	"context"
	"errors"

	"github.com/relby/diva.back/internal/convert"
	"github.com/relby/diva.back/internal/model"
	"github.com/relby/diva.back/internal/service"
	"github.com/relby/diva.back/pkg/genproto"
)

func (server *GRPCServer) GetCustomers(ctx context.Context, req *genproto.GetCustomersRequest) (*genproto.GetCustomersResponse, error) {
	customers, err := server.customerService.GetManyCustomers(ctx, service.CustomerServiceGetManyCustomersOptions{
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

func (server *GRPCServer) SetCustomerDiscountById(ctx context.Context, req *genproto.SetCustomerDiscountByIdRequest) (*genproto.SetCustomerDiscountByIdResponse, error) {
	customerID, err := model.NewCustomerID(req.Id)
	if err != nil {
		return nil, err
	}
	customerDiscount, err := model.NewCustomerDiscount(req.Discount)
	if err != nil {
		return nil, err
	}

	customer, err := server.customerService.SetCustomerDiscountByID(ctx, customerID, customerDiscount)
	if err != nil {
		return nil, err
	}

	return &genproto.SetCustomerDiscountByIdResponse{
		Customer: convert.CustomerFromModelToProto(customer),
	}, nil
}

func (server *GRPCServer) ExportCustomersToExcel(ctx context.Context, req *genproto.ExportCustomersToExcelRequest) (*genproto.ExportCustomersToExcelResponse, error) {
	// TODO
	return nil, errors.New("UNIMPLEMENTED")
}
