package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/relby/diva.back/internal/api"
	"github.com/relby/diva.back/internal/closer"
	"github.com/relby/diva.back/internal/config"
	"github.com/relby/diva.back/internal/config/env"
	"github.com/relby/diva.back/internal/repository"
	"github.com/relby/diva.back/internal/repository/postgres"
	"github.com/relby/diva.back/internal/service"
	"github.com/relby/diva.back/pkg/gensqlc"
)

type DIContainer struct {
	postgresConfig config.PostgresConfig
	grpcConfig     config.GRPCConfig
	postgresPool   *pgxpool.Pool

	queries *gensqlc.Queries

	customerRepository repository.CustomerRepository
	customerService    *service.CustomerService
	grpcServer         *api.GRPCServer
}

func NewDIContainer() (*DIContainer, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	return &DIContainer{}, nil
}

func (diContainer *DIContainer) GRPCConfig() (config.GRPCConfig, error) {
	if diContainer.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			return nil, err
		}
		diContainer.grpcConfig = cfg
	}

	return diContainer.grpcConfig, nil
}

func (diContainer *DIContainer) PostgresConfig() (config.PostgresConfig, error) {
	if diContainer.postgresConfig == nil {
		cfg, err := env.NewPostgresConfig()
		if err != nil {
			return nil, err
		}

		diContainer.postgresConfig = cfg
	}

	return diContainer.postgresConfig, nil
}

func (diContainer *DIContainer) PostgresPool(ctx context.Context) (*pgxpool.Pool, error) {
	if diContainer.postgresPool == nil {
		postgresConfig, err := env.NewPostgresConfig()
		if err != nil {
			return nil, err
		}
		postgresPool, err := pgxpool.New(ctx, postgresConfig.DSN())
		if err != nil {
			return nil, err
		}
		if err := postgresPool.Ping(ctx); err != nil {
			return nil, err
		}
		closer.Add(func() error {
			postgresPool.Close()
			return nil
		})

		diContainer.postgresPool = postgresPool
	}

	return diContainer.postgresPool, nil
}

func (diContainer *DIContainer) Queries(ctx context.Context) (*gensqlc.Queries, error) {
	if diContainer.queries == nil {
		postgresPool, err := diContainer.PostgresPool(ctx)
		if err != nil {
			return nil, err
		}

		diContainer.queries = gensqlc.New(postgresPool)
	}

	return diContainer.queries, nil
}

func (diContainer *DIContainer) CustomerRepository(ctx context.Context) (repository.CustomerRepository, error) {
	if diContainer.customerRepository == nil {
		queries, err := diContainer.Queries(ctx)
		if err != nil {
			return nil, err
		}

		diContainer.customerRepository = postgres.NewCustomerRepository(queries)
	}

	return diContainer.customerRepository, nil
}

func (diContainer *DIContainer) CustomerService(ctx context.Context) (*service.CustomerService, error) {
	if diContainer.customerService == nil {
		customerRepository, err := diContainer.CustomerRepository(ctx)
		if err != nil {
			return nil, err
		}

		diContainer.customerService = service.NewCustomerService(customerRepository)
	}

	return diContainer.customerService, nil
}

func (diContainer *DIContainer) GRPCServer(ctx context.Context) (*api.GRPCServer, error) {
	if diContainer.grpcServer == nil {
		customerService, err := diContainer.CustomerService(ctx)
		if err != nil {
			return nil, err
		}
		diContainer.grpcServer = api.NewGRPCServer(customerService)
	}

	return diContainer.grpcServer, nil
}
