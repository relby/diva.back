package model

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/relby/diva.back/internal/domainerrors"
)

type EmployeeAccessKey string

func NewEmployeeAccessKey(accessKey string) (EmployeeAccessKey, error) {
	accessKeyRegexp := regexp.MustCompile(`^[0-9]{6}$`)
	if !accessKeyRegexp.MatchString(accessKey) {
		return "", domainerrors.NewValidationError(fmt.Sprintf("employee access key must match regexp: `%s`", accessKeyRegexp.String()))
	}

	return EmployeeAccessKey(accessKey), nil
}

type EmployeePermission string

const (
	EmployeePermissionCreate EmployeePermission = "CREATE"
	EmployeePermissionUpdate EmployeePermission = "UPDATE"
	EmployeePermissionDelete EmployeePermission = "DELETE"
)

func NewEmployeePermission[T ~string](permission T) (EmployeePermission, error) {
	permissions := []EmployeePermission{
		EmployeePermissionCreate,
		EmployeePermissionUpdate,
		EmployeePermissionDelete,
	}

	for _, p := range permissions {
		if string(p) == string(permission) {
			return p, nil
		}
	}

	return "", domainerrors.NewValidationError(fmt.Sprintf("employee permission must be one of: %v", permissions))
}

type EmployeePermissions []EmployeePermission

var errEmployeePermissionsHaveDuplicates = domainerrors.NewValidationError("employee permissions can't have duplicates")

func NewEmployeePermissions(permissions []EmployeePermission) (EmployeePermissions, error) {
	permissionsSet := make(map[EmployeePermission]struct{}, len(permissions))

	for _, permission := range permissions {
		if _, found := permissionsSet[permission]; found {
			return nil, errEmployeePermissionsHaveDuplicates
		}
		permissionsSet[permission] = struct{}{}
	}

	return EmployeePermissions(permissions), nil
}
func NewEmployeePermissionsPanic(permissions []EmployeePermission) EmployeePermissions {
	permissions, err := NewEmployeePermissions(permissions)
	if err != nil {
		panic(errEmployeePermissionsHaveDuplicates)
	}

	return permissions
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

func (employee *Employee) HasPermissions(permissions EmployeePermissions) bool {
	for _, permission := range permissions {
		found := false
		for _, employeePermission := range employee.permissions {
			if permission == employeePermission {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}
