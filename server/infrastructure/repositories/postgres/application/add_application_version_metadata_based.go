package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/infrastructure/repositories/postgres/application/sql/edit"
	"github.com/Ponchitos/application_service/server/internal/models"
	"github.com/jackc/pgx/v4"
)

func (repo *repository) AddApplicationVersionMetadataBased(ctx context.Context, metadata *models.ApplicationMetadata, applicationUUID, enterpriseID string) (string, error) {

	tx, err := repo.client.Begin(ctx)
	if err != nil {
		return "", err
	}

	err = repo.setUpdateAvailableStatusToApplicationTx(ctx, tx, applicationUUID, enterpriseID)
	if err != nil {
		return "", err
	}

	applicationVersionID, applicationVersionUUID, err := repo.addApplicationVersionMetadataBasedTx(ctx, tx, applicationUUID, metadata)
	if err != nil {
		return "", err
	}

	err = repo.updateLatestApplicationVersionTx(ctx, tx, applicationUUID, enterpriseID, applicationVersionID)
	if err != nil {
		return "", err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", err
	}

	return applicationVersionUUID, nil
}

func (repo *repository) setUpdateAvailableStatusToApplicationTx(ctx context.Context, tx pgx.Tx, applicationUUID, enterpriseID string) error {
	_, err := tx.Exec(ctx, edit.SetUpdateAvailableStatusToApplicationSQL, applicationUUID, enterpriseID)

	if err != nil {
		repo.lgr.Error("setUpdateAvailableStatusToApplicationTx: execute sql query")

		repo.txRollback(ctx, tx)

		return err
	}

	return nil
}

func (repo *repository) updateLatestApplicationVersionTx(ctx context.Context, tx pgx.Tx, applicationUUID, enterpriseID string, applicationVersionID int) error {
	_, err := tx.Exec(ctx, edit.LatestApplicationSQL, applicationVersionID, applicationUUID, enterpriseID)
	if err != nil {
		repo.lgr.Error("updateLatestApplicationVersionTx: execute sql query")

		repo.txRollback(ctx, tx)

		return err
	}

	return nil
}
