package api

import (
	"github.com/relby/diva.back/internal/service"
	"github.com/relby/diva.back/pkg/genproto"
)

type GRPCServer struct {
	genproto.UnimplementedCustomersServiceServer
	genproto.UnimplementedEmployeesServiceServer

	customerService *service.CustomerService
	employeeService *service.EmployeeService
}

func NewGRPCServer(customerService *service.CustomerService, employeeService *service.EmployeeService) *GRPCServer {
	return &GRPCServer{
		customerService: customerService,
		employeeService: employeeService,
	}
}
