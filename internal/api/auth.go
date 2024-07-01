package api

import (
	"context"

	"github.com/relby/diva.back/internal/model"
	"github.com/relby/diva.back/pkg/genproto"
)

func (server *GRPCServer) EmployeeLogin(ctx context.Context, req *genproto.EmployeeLoginRequest) (*genproto.EmployeeLoginResponse, error) {
	accessKey, err := model.NewEmployeeAccessKey(req.AccessKey)
	if err != nil {
		return nil, err
	}

	accessToken, err := server.authService.EmployeeLogin(ctx, accessKey)
	if err != nil {
		return nil, err
	}

	return &genproto.EmployeeLoginResponse{
		AccessToken: accessToken,
	}, nil
}

func (server *GRPCServer) AdminLogin(ctx context.Context, req *genproto.AdminLoginRequest) (*genproto.AdminLoginResponse, error) {
	login, err := model.NewAdminLogin(req.Login)
	if err != nil {
		return nil, err
	}

	accessToken, err := server.authService.AdminLogin(ctx, login, req.Password)
	if err != nil {
		return nil, err
	}

	return &genproto.AdminLoginResponse{
		AccessToken: accessToken,
	}, nil
}