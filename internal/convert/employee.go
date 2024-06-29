package convert

import (
	"github.com/relby/diva.back/internal/model"
	"github.com/relby/diva.back/pkg/genproto"
	"github.com/relby/diva.back/pkg/gensqlc"
)

func EmployeeFromRowToModel(userRow gensqlc.User, employeeRow gensqlc.Employee) (*model.Employee, error) {
	id, fullName, err := userFromRowToValueObjects(userRow)
	if err != nil {
		return nil, err
	}

	accessKey, err := model.NewEmployeeAccessKey(employeeRow.AccessKey)
	if err != nil {
		return nil, err
	}

	permissionSlice := make([]model.EmployeePermission, len(employeeRow.Permissions))
	for i, permission := range employeeRow.Permissions {
		permissionSlice[i], err = model.NewEmployeePermission(permission)
		if err != nil {
			return nil, err
		}
	}

	permissions, err := model.NewEmployeePermissions(permissionSlice)
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
	var err error
	permissionsSlice := make([]model.EmployeePermission, len(permissionsProto))

	for i, permissionProto := range permissionsProto {
		permissionsSlice[i], err = model.NewEmployeePermission(genproto.EmployeePermission_name[int32(permissionProto)])
		if err != nil {
			return nil, err
		}
	}

	permissionsModel, err := model.NewEmployeePermissions(permissionsSlice)
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
