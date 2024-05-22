package model

import (
	"errors"

	"github.com/google/uuid"
)

type UserID uuid.UUID

func NewUserID(id uuid.UUID) (UserID, error) {
	if id == uuid.Nil {
		return UserID(uuid.Nil), errors.New("user id can't be nil")
	}

	return UserID(id), nil
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
