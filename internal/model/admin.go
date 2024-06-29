package model

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AdminLogin string

func NewAdminLogin(login string) (AdminLogin, error) {
	return AdminLogin(login), nil
}

type AdminHashedPassword string

func NewAdminHashedPassword(hashedPassword string) (AdminHashedPassword, error) {
	if hashedPassword == "" {
		return "", errors.New("admin password is empty")
	}

	return AdminHashedPassword(hashedPassword), nil
}

func NewAdminHashedPasswordFromPassword(password string) (AdminHashedPassword, error) {
	if password == "" {
		return "", errors.New("admin password is empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return NewAdminHashedPassword(string(hashedPassword))
}

type Admin struct {
	user
	login          AdminLogin
	hashedPassword AdminHashedPassword
}

func NewAdmin(id UserID, fullName UserFullName, login AdminLogin, hashedPassword AdminHashedPassword) (*Admin, error) {
	return &Admin{
		user: user{
			id:       id,
			fullName: fullName,
		},
		login:          login,
		hashedPassword: hashedPassword,
	}, nil
}

func NewAdminWithRandomID(fullName UserFullName, login AdminLogin, hashedPassword AdminHashedPassword) (*Admin, error) {
	id, err := NewUserID(uuid.New())
	if err != nil {
		return nil, err
	}

	admin, err := NewAdmin(id, fullName, login, hashedPassword)
	if err != nil {
		return nil, err
	}

	return admin, nil
}

func (admin *Admin) User() user {
	return admin.user
}

func (admin *Admin) Login() AdminLogin {
	return admin.login
}

func (admin *Admin) HashedPassword() AdminHashedPassword {
	return admin.hashedPassword
}

func (admin *Admin) PasswordMathes(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(string(admin.hashedPassword)), []byte(password))

	return err == nil
}
