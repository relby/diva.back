package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/relby/diva.back/internal/model"
	"github.com/relby/diva.back/internal/repository"
	"github.com/relby/diva.back/pkg/gensqlc"
)

var _ repository.AdminRepository = (*AdminRepository)(nil)

type AdminRepository struct {
	postgresPool *pgxpool.Pool
	queries      *gensqlc.Queries
}

func NewAdminRepository(postgresPool *pgxpool.Pool, queries *gensqlc.Queries) *AdminRepository {
	return &AdminRepository{
		postgresPool: postgresPool,
		queries:      queries,
	}
}

func (repository *AdminRepository) Save(ctx context.Context, admin *model.Admin) error {
	tx, err := repository.postgresPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queriesTx := repository.queries.WithTx(tx)

	if err = queriesTx.UpsertUser(ctx, gensqlc.UpsertUserParams{
		ID:       uuid.UUID(admin.ID()),
		FullName: string(admin.FullName()),
	}); err != nil {
		return err
	}

	if err = queriesTx.UpsertAdmin(ctx, gensqlc.UpsertAdminParams{
		UserID:         uuid.UUID(admin.ID()),
		Login:          string(admin.Login()),
		HashedPassword: string(admin.HashedPassword()),
	}); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
