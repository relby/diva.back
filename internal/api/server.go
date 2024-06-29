package api

import (
	"github.com/relby/diva.back/internal/service"
	"github.com/relby/diva.back/pkg/genproto"
)

type GRPCServer struct {
	genproto.UnimplementedCustomersServiceServer
	genproto.UnimplementedEmployeesServiceServer
	genproto.UnimplementedAuthServiceServer

	customerService *service.CustomerService
	employeeService *service.EmployeeService
	authService     *service.AuthService
}

func NewGRPCServer(
	customerService *service.CustomerService,
	employeeService *service.EmployeeService,
	authService *service.AuthService,
) *GRPCServer {
	return &GRPCServer{
		customerService: customerService,
		employeeService: employeeService,
		authService:     authService,
	}
}
