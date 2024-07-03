package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
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
	authConfig             config.AuthConfig
	employeeRepository     repository.EmployeeRepository
	adminRepository        repository.AdminRepository
	refreshTokenRepository repository.RefreshTokenRepository
}

func NewAuthService(
	authConfig config.AuthConfig,
	employeeRepository repository.EmployeeRepository,
	adminRepository repository.AdminRepository,
	refreshTokenRepository repository.RefreshTokenRepository,
) *AuthService {
	return &AuthService{
		authConfig:             authConfig,
		employeeRepository:     employeeRepository,
		adminRepository:        adminRepository,
		refreshTokenRepository: refreshTokenRepository,
	}
}

func (service *AuthService) EmployeeLogin(ctx context.Context, accessKey model.EmployeeAccessKey) (string, string, error) {
	employee, err := service.employeeRepository.GetByAccessKey(ctx, accessKey)
	if err != nil {
		return "", "", err
	}

	accessToken, err := jwt.NewAccessToken(&jwt.AccessTokenClaims{
		UserID:   employee.ID(),
		UserType: model.UserTypeEmployee,
	}, service.authConfig.AccessTokenSecret(), service.authConfig.AccessTokenExpireDuration())
	if err != nil {
		return "", "", err
	}

	refreshTokenID, err := model.NewRefreshTokenID(uuid.New())
	if err != nil {
		return "", "", err
	}

	refreshTokenExpiresAt, err := model.NewRefreshTokenExpiresAt(time.Now().Add(service.authConfig.RefreshTokenExpireDuration()))
	if err != nil {
		return "", "", err
	}

	refreshTokenModel, err := model.NewRefreshToken(refreshTokenID, employee.ID(), refreshTokenExpiresAt)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewRefreshToken(&jwt.RefreshTokenClaims{
		ID:       refreshTokenModel.ID(),
		UserID:   refreshTokenModel.UserID(),
		UserType: model.UserTypeEmployee,
	}, service.authConfig.RefreshTokenSecret(), service.authConfig.RefreshTokenExpireDuration())
	if err != nil {
		return "", "", err
	}

	if err := service.refreshTokenRepository.Save(ctx, refreshTokenModel); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (service *AuthService) AdminLogin(ctx context.Context, login model.AdminLogin, password string) (string, string, error) {
	admin, err := service.adminRepository.GetByLogin(ctx, login)
	if err != nil {
		return "", "", err
	}

	if !admin.PasswordMathes(password) {
		return "", "", errors.New("TODO")
	}

	accessToken, err := jwt.NewAccessToken(&jwt.AccessTokenClaims{
		UserID:   admin.ID(),
		UserType: model.UserTypeAdmin,
	}, service.authConfig.AccessTokenSecret(), service.authConfig.AccessTokenExpireDuration())

	if err != nil {
		// TODO: create domain error
		return "", "", err
	}

	refreshTokenID, err := model.NewRefreshTokenID(uuid.New())
	if err != nil {
		return "", "", err
	}

	refreshTokenExpiresAt, err := model.NewRefreshTokenExpiresAt(time.Now().Add(service.authConfig.RefreshTokenExpireDuration()))
	if err != nil {
		return "", "", err
	}

	refreshTokenModel, err := model.NewRefreshToken(refreshTokenID, admin.ID(), refreshTokenExpiresAt)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewRefreshToken(&jwt.RefreshTokenClaims{
		ID:       refreshTokenModel.ID(),
		UserID:   refreshTokenModel.UserID(),
		UserType: model.UserTypeAdmin,
	}, service.authConfig.RefreshTokenSecret(), service.authConfig.RefreshTokenExpireDuration())
	if err != nil {
		return "", "", err
	}

	if err := service.refreshTokenRepository.Save(ctx, refreshTokenModel); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (service *AuthService) Refresh(ctx context.Context, oldRefreshToken string) (string, string, error) {
	oldRefreshTokenClaims, err := jwt.ParseRefreshToken(oldRefreshToken, service.authConfig.RefreshTokenSecret())
	if err != nil {
		return "", "", err
	}

	oldRefreshTokenModel, err := service.refreshTokenRepository.GetByID(ctx, oldRefreshTokenClaims.ID)
	if err != nil {
		return "", "", err
	}

	accessToken, err := jwt.NewAccessToken(&jwt.AccessTokenClaims{
		UserID:   oldRefreshTokenModel.UserID(),
		UserType: oldRefreshTokenClaims.UserType,
	}, service.authConfig.AccessTokenSecret(), service.authConfig.AccessTokenExpireDuration())
	if err != nil {
		return "", "", err
	}

	newRefreshTokenID, err := model.NewRefreshTokenID(uuid.New())
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := jwt.NewRefreshToken(&jwt.RefreshTokenClaims{
		ID:       newRefreshTokenID,
		UserID:   oldRefreshTokenModel.UserID(),
		UserType: oldRefreshTokenClaims.UserType,
	}, service.authConfig.RefreshTokenSecret(), service.authConfig.RefreshTokenExpireDuration())
	if err != nil {
		return "", "", err
	}

	newRefreshTokenExpiresAt, err := model.NewRefreshTokenExpiresAt(time.Now().Add(service.authConfig.RefreshTokenExpireDuration()))
	if err != nil {
		return "", "", err
	}

	newRefreshTokenModel, err := model.NewRefreshToken(
		newRefreshTokenID,
		oldRefreshTokenModel.UserID(),
		newRefreshTokenExpiresAt,
	)
	if err != nil {
		return "", "", err
	}

	// TODO: all of the repository operations must be within a transaction. introduce transaction manager
	if err := service.refreshTokenRepository.Save(ctx, newRefreshTokenModel); err != nil {
		return "", "", err
	}

	if err := service.refreshTokenRepository.Delete(ctx, oldRefreshTokenModel); err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func (service *AuthService) Logout(ctx context.Context, refreshToken string) error {
	refreshTokenClaims, err := jwt.ParseRefreshToken(refreshToken, service.authConfig.RefreshTokenSecret())
	if err != nil {
		return err
	}

	refreshTokenModel, err := service.refreshTokenRepository.GetByID(ctx, refreshTokenClaims.ID)
	if err != nil {
		return err
	}

	err = service.refreshTokenRepository.Delete(ctx, refreshTokenModel)
	if err != nil {
		return err
	}

	return nil
}
