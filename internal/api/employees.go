package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/relby/diva.back/internal/convert"
	"github.com/relby/diva.back/internal/model"
	"github.com/relby/diva.back/internal/service"
	"github.com/relby/diva.back/pkg/genproto"
)

func (server *GRPCServer) GetEmployees(ctx context.Context, req *genproto.GetEmployeesRequest) (*genproto.GetEmployeesResponse, error) {
	employees, err := server.employeeService.ListEmployees(ctx)
	if err != nil {
		return nil, err
	}

	employeesProto := make([]*genproto.Employee, len(employees))
	for i, employee := range employees {
		employeesProto[i] = convert.EmployeeFromModelToProto(employee)
	}

	return &genproto.GetEmployeesResponse{
		Employees: employeesProto,
	}, nil
}

func (server *GRPCServer) AddEmployee(ctx context.Context, req *genproto.AddEmployeeRequest) (*genproto.AddEmployeeResponse, error) {
	fullName, err := model.NewUserFullName(req.FullName)
	if err != nil {
		return nil, err
	}

	accessKey, err := model.NewEmployeeAccessKey(req.AccessKey)
	if err != nil {
		return nil, err
	}

	employee, err := server.employeeService.AddEmployee(ctx, fullName, accessKey, nil)
	if err != nil {
		return nil, err
	}

	return &genproto.AddEmployeeResponse{
		Employee: convert.EmployeeFromModelToProto(employee),
	}, nil
}

func (server *GRPCServer) UpdateEmployee(ctx context.Context, req *genproto.UpdateEmployeeRequest) (*genproto.UpdateEmployeeResponse, error) {
	id, err := uuid.Parse(req.Id.Value)
	if err != nil {
		return nil, err
	}

	userID, err := model.NewUserID(id)
	if err != nil {
		return nil, err
	}

	employee, err := server.employeeService.GetEmployeeByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var fullName *model.UserFullName
	var accessKey *model.EmployeeAccessKey
	var permissions *model.EmployeePermissions
	if req.FullName != nil {
		value, err := model.NewUserFullName(*req.FullName)
		if err != nil {
			return nil, err
		}
		fullName = &value
	}

	if req.AccessKey != nil {
		value, err := model.NewEmployeeAccessKey(*req.AccessKey)
		if err != nil {
			return nil, err
		}
		accessKey = &value
	}

	if req.Permissions != nil {
		value, err := convert.EmployeePermissionsFromProtoToModel(req.Permissions.Permissions)
		if err != nil {
			return nil, err
		}

		permissions = &value
	}

	if err := server.employeeService.UpdateEmployee(ctx, employee, &service.UpdateEmployeeValues{
		FullName:    fullName,
		AccessKey:   accessKey,
		Permissions: permissions,
	}); err != nil {
		return nil, err
	}

	return &genproto.UpdateEmployeeResponse{
		Employee: convert.EmployeeFromModelToProto(employee),
	}, nil
}

func (server *GRPCServer) DeleteEmployee(ctx context.Context, req *genproto.DeleteEmployeeRequest) (*genproto.DeleteEmployeeResponse, error) {
	id, err := uuid.Parse(req.Id.Value)
	if err != nil {
		return nil, err
	}

	employeeId, err := model.NewUserID(id)
	if err != nil {
		return nil, err
	}

	employee, err := server.employeeService.GetEmployeeByID(ctx, employeeId)
	if err != nil {
		return nil, err
	}

	if err = server.employeeService.DeleteEmployee(ctx, employee); err != nil {
		return nil, err
	}

	return &genproto.DeleteEmployeeResponse{
		Employee: convert.EmployeeFromModelToProto(employee),
	}, nil
}
