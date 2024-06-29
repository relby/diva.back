package interceptor

import (
	"context"
	"strings"

	"github.com/relby/diva.back/internal/config"
	"github.com/relby/diva.back/internal/model"
	"github.com/relby/diva.back/internal/service"
	"github.com/relby/diva.back/pkg/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	authConfig  config.AuthConfig
	authService *service.AuthService
}

func NewAuthInterceptor(authConfig config.AuthConfig, authService *service.AuthService) *AuthInterceptor {
	return &AuthInterceptor{authConfig: authConfig, authService: authService}
}

func (interceptor *AuthInterceptor) Unary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	methodPermissions, ok := getMethodPermissions(info.FullMethod)
	if !ok {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, status.Error(codes.InvalidArgument, "metadata not found")
	}

	authTokens := md.Get("Authorization")

	if len(authTokens) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization token is missing")
	}

	authHeader := authTokens[0]
	fields := strings.Fields(authHeader)

	if len(fields) < 2 {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth header format. use `Bearer <token>`")
	}

	authType := fields[0]

	if authType != "Bearer" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization type: %v", authType)
	}

	accessToken := fields[1]

	claims, err := jwt.ParseAccessToken(accessToken, interceptor.authConfig.JWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	if !interceptor.authService.CheckAccess(ctx, claims.UserID, claims.UserType, methodPermissions.isAdminOnly, methodPermissions.employeePermissions) {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	return handler(ctx, req)
}

type MethodPermission struct {
	employeePermissions model.EmployeePermissions
	isAdminOnly         bool
}

func getMethodPermissions(fullMethod string) (*MethodPermission, bool) {
	methodPermissions := map[string]MethodPermission{}

	methodPermissions["/customers.CustomersService/ListCustomers"] = MethodPermission{
		employeePermissions: []model.EmployeePermission{
			model.EmployeePermissionRead,
		},
	}
	methodPermissions["/customers.CustomersService/UpdateCustomer"] = MethodPermission{
		employeePermissions: []model.EmployeePermission{
			model.EmployeePermissionUpdate,
		},
	}
	methodPermissions["/customers.CustomersService/AddCustomer"] = MethodPermission{
		employeePermissions: []model.EmployeePermission{
			model.EmployeePermissionCreate,
		},
	}
	methodPermissions["/customers.CustomersService/DeleteCustomer"] = MethodPermission{
		employeePermissions: []model.EmployeePermission{
			model.EmployeePermissionDelete,
		},
	}
	methodPermissions["/customers.CustomersService/ExportCustomersToExcel"] = MethodPermission{
		isAdminOnly: true,
	}

	methodPermissions["/employees.EmployeesService/GetEmployees"] = MethodPermission{
		isAdminOnly: true,
	}
	methodPermissions["/employees.EmployeesService/AddEmployee"] = MethodPermission{
		isAdminOnly: true,
	}
	methodPermissions["/employees.EmployeesService/UpdateEmployee"] = MethodPermission{
		isAdminOnly: true,
	}
	methodPermissions["/employees.EmployeesService/DeleteEmployee"] = MethodPermission{
		isAdminOnly: true,
	}

	methodPermission, ok := methodPermissions[fullMethod]
	if !ok {
		return nil, false
	}

	return &methodPermission, true
}
