package service

import (
	"context"
	"errors"

	"github.com/relby/diva.back/internal/config"
	"github.com/relby/diva.back/internal/model"
	"github.com/relby/diva.back/internal/repository"
	"github.com/relby/diva.back/pkg/jwt"
)

func (service *AuthService) CheckAccess(ctx context.Context, userID model.UserID, userType model.UserType, isAdminOnly bool, permissions model.EmployeePermissions) bool {
	switch userType {
	case model.UserTypeAdmin:
		// TODO: if i introduce user repository change this to userRepository.Exists
		_, err := service.adminRepository.GetByID(ctx, userID)
		if err != nil {
			return false
		}

		return true
	case model.UserTypeEmployee:
		employee, err := service.employeeRepository.GetByID(ctx, userID)
		if err != nil {
			return false
		}

		if isAdminOnly {
			return false
		}

		return employee.HasPermissions(permissions)
	default:
		return false
	}
}

type AuthService struct {
	authConfig         config.AuthConfig
	employeeRepository repository.EmployeeRepository
	adminRepository    repository.AdminRepository
}

func NewAuthService(
	authConfig config.AuthConfig,
	employeeRepository repository.EmployeeRepository,
	adminRepository repository.AdminRepository,
) *AuthService {
	return &AuthService{
		authConfig:         authConfig,
		employeeRepository: employeeRepository,
		adminRepository:    adminRepository,
	}
}

func (service *AuthService) EmployeeLogin(ctx context.Context, accessKey model.EmployeeAccessKey) (string, error) {
	employee, err := service.employeeRepository.GetByAccessKey(ctx, accessKey)
	if err != nil {
		return "", err
	}

	accessToken, err := jwt.NewAccessToken(&jwt.AccessTokenClaims{
		UserID:   employee.ID(),
		UserType: model.UserTypeEmployee,
	}, service.authConfig.JWTSecret(), service.authConfig.JWTExpireDuration())

	if err != nil {
		return "", err
	}

	return accessToken, nil
}
func (service *AuthService) AdminLogin(ctx context.Context, login model.AdminLogin, password string) (string, error) {
	admin, err := service.adminRepository.GetByLogin(ctx, login)
	if err != nil {
		return "", err
	}

	if !admin.PasswordMathes(password) {
		return "", errors.New("TODO")
	}

	accessToken, err := jwt.NewAccessToken(&jwt.AccessTokenClaims{
		UserID:   admin.ID(),
		UserType: model.UserTypeAdmin,
	}, service.authConfig.JWTSecret(), service.authConfig.JWTExpireDuration())

	if err != nil {
		return "", err
	}

	return accessToken, nil
}
