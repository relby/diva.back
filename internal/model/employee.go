package model

import (
	"errors"

	"github.com/google/uuid"
)

type EmployeeAccessKey string

func NewEmployeeAccessKey(accessKey string) (EmployeeAccessKey, error) {
	if accessKey == "" {
		return "", errors.New("employee access key is empty")
	}

	return EmployeeAccessKey(accessKey), nil
}

type employeePermission string

const (
	employeePermissionCreate employeePermission = "CREATE"
	employeePermissionRead   employeePermission = "READ"
	employeePermissionUpdate employeePermission = "UPDATE"
	employeePermissionDelete employeePermission = "DELETE"
)

func newEmployeePermission(permission string) (employeePermission, error) {
	permissions := []employeePermission{
		employeePermissionCreate,
		employeePermissionRead,
		employeePermissionUpdate,
		employeePermissionDelete,
	}

	for _, p := range permissions {
		if string(p) == permission {
			return p, nil
		}
	}

	return "", errors.New("employee permission invalid")
}

type EmployeePermissions []employeePermission

func NewEmployeePermissions(p []string) (EmployeePermissions, error) {
	var err error

	permissions := make([]employeePermission, len(p))
	for i, permission := range p {
		permissions[i], err = newEmployeePermission(permission)
		if err != nil {
			return nil, err
		}
	}

	permissionSet := make(map[employeePermission]struct{}, len(permissions))

	for _, permission := range permissions {
		if _, found := permissionSet[permission]; found {
			return nil, errors.New("employee permissions can't have duplicates")
		}
		permissionSet[permission] = struct{}{}
	}

	return EmployeePermissions(permissions), nil
}

type Employee struct {
	user
	accessKey   EmployeeAccessKey
	permissions EmployeePermissions
}

func NewEmployee(id UserID, fullName UserFullName, accessKey EmployeeAccessKey, permissions EmployeePermissions) (*Employee, error) {
	return &Employee{
		user: user{
			id:       id,
			fullName: fullName,
		},
		accessKey:   accessKey,
		permissions: permissions,
	}, nil
}

func NewEmployeeWithRandomID(fullName UserFullName, accessKey EmployeeAccessKey, permissions EmployeePermissions) (*Employee, error) {
	id, err := NewUserID(uuid.New())
	if err != nil {
		return nil, err
	}
	employee, err := NewEmployee(id, fullName, accessKey, permissions)
	if err != nil {
		return nil, err
	}

	return employee, nil
}

func (employee *Employee) User() user {
	return employee.user
}

func (employee *Employee) AccessKey() EmployeeAccessKey {
	return employee.accessKey
}

func (employee *Employee) Permissions() EmployeePermissions {
	return employee.permissions
}

func (employee *Employee) SetAccessKey(accessKey EmployeeAccessKey) {
	employee.accessKey = accessKey
}

func (employee *Employee) SetPermissions(permisssions EmployeePermissions) {
	employee.permissions = permisssions
}
