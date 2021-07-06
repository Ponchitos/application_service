package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/internal/repositories"
	"github.com/Ponchitos/application_service/server/tools/logger"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type repository struct {
	client *pgxpool.Pool
	lgr    logger.Logger
}

func (repo *repository) txRollback(ctx context.Context, tx pgx.Tx) {
	if errRollback := tx.Rollback(ctx); errRollback != nil {
		repo.lgr.Error(errRollback)
	}
}

func NewApplicationRepository(client *pgxpool.Pool, lgr logger.Logger) repositories.ApplicationRepository {
	return &repository{client, lgr}
}
