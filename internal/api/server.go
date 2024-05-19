package api

import (
	"github.com/relby/diva.back/internal/service"
	"github.com/relby/diva.back/pkg/genproto"
)

type GRPCServer struct {
	genproto.UnimplementedCustomersServiceServer

	customerService *service.CustomerService
}

func NewGRPCServer(customerService *service.CustomerService) *GRPCServer {
	return &GRPCServer{
		customerService: customerService,
	}
}
