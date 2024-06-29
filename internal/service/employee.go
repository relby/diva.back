package service

import (
	"context"

	"github.com/relby/diva.back/internal/model"
	"github.com/relby/diva.back/internal/repository"
)

type EmployeeService struct {
	employeeRepository repository.EmployeeRepository
}

func NewEmployeeService(repository repository.EmployeeRepository) *EmployeeService {
	return &EmployeeService{
		employeeRepository: repository,
	}
}

func (service *EmployeeService) GetEmployeeByID(ctx context.Context, id model.UserID) (*model.Employee, error) {
	employee, err := service.employeeRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return employee, nil
}

func (service *EmployeeService) GetEmployeeByAccessKey(ctx context.Context, accessKey model.EmployeeAccessKey) (*model.Employee, error) {
	employee, err := service.employeeRepository.GetByAccessKey(ctx, accessKey)
	if err != nil {
		return nil, err
	}

	return employee, nil
}

func (service *EmployeeService) ListEmployees(ctx context.Context) ([]*model.Employee, error) {
	employees, err := service.employeeRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (service *EmployeeService) AddEmployee(ctx context.Context, fullName model.UserFullName, accessKey model.EmployeeAccessKey, permissions model.EmployeePermissions) (*model.Employee, error) {
	employee, err := model.NewEmployeeWithRandomID(
		fullName,
		accessKey,
		permissions,
	)
	if err != nil {
		return nil, err
	}

	if err := service.employeeRepository.Save(ctx, employee); err != nil {
		return nil, err
	}

	return employee, nil
}

type UpdateEmployeeValues struct {
	FullName    *model.UserFullName
	AccessKey   *model.EmployeeAccessKey
	Permissions *model.EmployeePermissions
}

func (service *EmployeeService) UpdateEmployee(ctx context.Context, employee *model.Employee, values *UpdateEmployeeValues) error {
	if values.FullName != nil {
		employee.SetFullName(*values.FullName)
	}

	if values.AccessKey != nil {
		employee.SetAccessKey(*values.AccessKey)
	}

	if values.Permissions != nil {
		employee.SetPermissions(*values.Permissions)
	}

	err := service.employeeRepository.Save(ctx, employee)
	if err != nil {
		return err
	}

	return nil
}

func (service *EmployeeService) DeleteEmployee(ctx context.Context, employee *model.Employee) error {
	if err := service.employeeRepository.Delete(ctx, employee); err != nil {
		return err
	}

	return nil
}
