package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/relby/diva.back/internal/convert"
	"github.com/relby/diva.back/internal/model"
	"github.com/relby/diva.back/internal/repository"
	"github.com/relby/diva.back/pkg/gensqlc"
)

var _ repository.EmployeeRepository = (*EmployeeRepository)(nil)

type EmployeeRepository struct {
	postgresPool *pgxpool.Pool
	queries      *gensqlc.Queries
}

func NewEmployeeRepository(postgresPool *pgxpool.Pool, queries *gensqlc.Queries) *EmployeeRepository {
	return &EmployeeRepository{
		postgresPool: postgresPool,
		queries:      queries,
	}
}

func (repository *EmployeeRepository) GetByID(ctx context.Context, id model.UserID) (*model.Employee, error) {
	employeeRow, err := repository.queries.SelectEmployeeByID(ctx, uuid.UUID(id))
	if err != nil {
		return nil, err
	}

	employeeModel, err := convert.EmployeeFromRowToModel(employeeRow.User, employeeRow.Employee)
	if err != nil {
		return nil, err
	}

	return employeeModel, nil
}

func (repository *EmployeeRepository) GetByAccessKey(ctx context.Context, accessKey model.EmployeeAccessKey) (*model.Employee, error) {
	employeeRow, err := repository.queries.SelectEmployeeByAccessKey(ctx, string(accessKey))
	if err != nil {
		return nil, err
	}

	employeeModel, err := convert.EmployeeFromRowToModel(employeeRow.User, employeeRow.Employee)
	if err != nil {
		return nil, err
	}

	return employeeModel, nil
}

func (repository *EmployeeRepository) List(ctx context.Context) ([]*model.Employee, error) {
	employeesRow, err := repository.queries.SelectEmployees(ctx)
	if err != nil {
		return nil, err
	}

	employeesModel := make([]*model.Employee, len(employeesRow))
	for i, row := range employeesRow {
		employeesModel[i], err = convert.EmployeeFromRowToModel(row.User, row.Employee)
		if err != nil {
			return nil, err
		}
	}

	return employeesModel, nil
}

func (repository *EmployeeRepository) Save(ctx context.Context, employee *model.Employee) error {
	tx, err := repository.postgresPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queriesTx := repository.queries.WithTx(tx)

	permissions := make([]gensqlc.EmployeePermission, len(employee.Permissions()))
	for i, permission := range employee.Permissions() {
		permissions[i] = gensqlc.EmployeePermission(permission)
	}

	if err = queriesTx.UpsertUser(ctx, gensqlc.UpsertUserParams{
		ID:       uuid.UUID(employee.ID()),
		FullName: string(employee.FullName()),
	}); err != nil {
		return err
	}

	if err = queriesTx.UpsertEmployee(ctx, gensqlc.UpsertEmployeeParams{
		UserID:      uuid.UUID(employee.ID()),
		AccessKey:   string(employee.AccessKey()),
		Permissions: permissions,
	}); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (repository *EmployeeRepository) Delete(ctx context.Context, employee *model.Employee) error {
	if err := repository.queries.DeleteUser(ctx, uuid.UUID(employee.ID())); err != nil {
		return err
	}

	return nil
}
