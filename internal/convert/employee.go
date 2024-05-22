package convert

import (
	"github.com/relby/diva.back/internal/model"
	"github.com/relby/diva.back/pkg/genproto"
	"github.com/relby/diva.back/pkg/gensqlc"
)

func EmployeeFromRowToModel(userRow gensqlc.User, employeeRow gensqlc.Employee) (*model.Employee, error) {
	id, err := model.NewUserID(userRow.ID)
	if err != nil {
		return nil, err
	}

	fullName, err := model.NewUserFullName(userRow.FullName)
	if err != nil {
		return nil, err
	}

	accessKey, err := model.NewEmployeeAccessKey(employeeRow.AccessKey)
	if err != nil {
		return nil, err
	}

	permissionsString := make([]string, len(employeeRow.Permissions))
	for i, permission := range employeeRow.Permissions {
		permissionsString[i] = string(permission)
	}

	permissions, err := model.NewEmployeePermissions(permissionsString)
	if err != nil {
		return nil, err
	}

	employee, err := model.NewEmployee(
		id,
		fullName,
		accessKey,
		permissions,
	)
	if err != nil {
		return nil, err
	}

	return employee, nil
}

func EmployeePermissionsFromModelToProto(permissionsModel model.EmployeePermissions) []genproto.EmployeePermission {
	permissionsProto := make([]genproto.EmployeePermission, len(permissionsModel))

	for i, permissionModel := range permissionsModel {
		permissionsProto[i] = genproto.EmployeePermission(genproto.EmployeePermission_value[string(permissionModel)])
	}

	return permissionsProto
}

func EmployeePermissionsFromProtoToModel(permissionsProto []genproto.EmployeePermission) (model.EmployeePermissions, error) {
	permissionsString := make([]string, len(permissionsProto))

	for i, permissionProto := range permissionsProto {
		permissionsString[i] = genproto.EmployeePermission_name[int32(permissionProto)]
	}

	permissionsModel, err := model.NewEmployeePermissions(permissionsString)
	if err != nil {
		return nil, err
	}

	return permissionsModel, nil
}

func EmployeeFromModelToProto(employee *model.Employee) *genproto.Employee {
	id := &genproto.UUID{
		Value: employee.ID().String(),
	}

	return &genproto.Employee{
		Id:          id,
		FullName:    string(employee.FullName()),
		AccessKey:   string(employee.AccessKey()),
		Permissions: EmployeePermissionsFromModelToProto(employee.Permissions()),
	}
}
