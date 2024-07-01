package model

import (
	"github.com/google/uuid"
	"github.com/relby/diva.back/internal/domainerrors"
)

type UserType string

const (
	UserTypeAdmin    UserType = "admin"
	UserTypeEmployee UserType = "employee"
)

func NewUserType(userType string) (UserType, error) {
	switch userType {
	case string(UserTypeAdmin):
		return UserTypeAdmin, nil
	case string(UserTypeEmployee):
		return UserTypeEmployee, nil
	default:
		return "", nil
	}
}

type UserID uuid.UUID

func NewUserID(id uuid.UUID) (UserID, error) {
	if id == uuid.Nil {
		return UserID(uuid.Nil), domainerrors.NewValidationError("user id can't be nil")
	}

	return UserID(id), nil
}

func NewUserIDFromString(id string) (UserID, error) {
	idUuid, err := uuid.Parse(id)
	if err != nil {
		return UserID(uuid.Nil), err
	}

	return NewUserID(idUuid)
}

func (id UserID) String() string {
	return uuid.UUID(id).String()
}

type UserFullName string

func NewUserFullName(name string) (UserFullName, error) {
	return UserFullName(name), nil
}

type User interface {
	ID() UserID
	FullName() UserFullName

	SetFullName(fullName UserFullName)
}

var _ User = (*user)(nil)

type user struct {
	id       UserID
	fullName UserFullName
}

func (user *user) ID() UserID {
	return user.id
}

func (user *user) FullName() UserFullName {
	return user.fullName
}

func (user *user) SetFullName(fullName UserFullName) {
	user.fullName = fullName
}
