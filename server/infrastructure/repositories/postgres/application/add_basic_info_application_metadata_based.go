package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/infrastructure/repositories/postgres/application/sql/add"
	"github.com/Ponchitos/application_service/server/internal/models"
	"github.com/jackc/pgx/v4"
)

func (repo *repository) AddBasicInfoApplicationMetadataBased(ctx context.Context, metadata *models.ApplicationMetadata, available, location, enterpriseID string) (string, string, error) {
	tx, err := repo.client.Begin(ctx)
	if err != nil {
		return "", "", err
	}

	applicationID, applicationUUID, err := repo.addBasicInfoApplication(ctx, tx, enterpriseID, metadata.PackageName, metadata.ApplicationLabel, available, location)
	if err != nil {
		return "", "", err
	}

	applicationVersionID, applicationVersionUUID, err := repo.addApplicationVersionMetadataBasedTx(ctx, tx, applicationUUID, metadata)
	if err != nil {
		return "", "", err
	}

	err = repo.addLatestApplicationVersionMetadataBasedTx(ctx, tx, applicationID, applicationVersionID, enterpriseID)
	if err != nil {
		return "", "", err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", "", err
	}

	return applicationUUID, applicationVersionUUID, nil
}

func (repo *repository) addBasicInfoApplication(ctx context.Context, tx pgx.Tx, enterpriseID, packageName, name, available, location string) (int, string, error) {
	var (
		applicationID   int
		applicationUUID string
	)

	row := tx.QueryRow(ctx, add.ApplicationSQL, enterpriseID, packageName, available, location, name, models.Approved)

	err := row.Scan(&applicationID, &applicationUUID)

	if err != nil {
		repo.lgr.Error("addBasicInfoApplication: row scan error")

		repo.txRollback(ctx, tx)

		return applicationID, applicationUUID, err
	}

	return applicationID, applicationUUID, nil
}

func (repo *repository) addApplicationVersionMetadataBasedTx(ctx context.Context, tx pgx.Tx, applicationUUID string, metadata *models.ApplicationMetadata) (int, string, error) {
	var (
		applicationVersionUUID string
		applicationVersionID   int
	)
	row := tx.QueryRow(ctx, add.ApplicationVersionMetadataBasedSQL, metadata.ID, applicationUUID, metadata.VersionCode, metadata.VersionName, metadata.MinimumSDK, metadata.IconBase64)

	err := row.Scan(&applicationVersionID, &applicationVersionUUID)

	if err != nil {
		repo.lgr.Error("addApplicationVersionMetadataBasedTx: row scan error")

		repo.txRollback(ctx, tx)

		return applicationVersionID, applicationVersionUUID, err
	}

	return applicationVersionID, applicationVersionUUID, nil
}

func (repo *repository) addLatestApplicationVersionMetadataBasedTx(ctx context.Context, tx pgx.Tx, applicationID, applicationVersionID int, enterpriseID string) error {
	_, err := tx.Exec(ctx, add.LatestApplicationVersionMetadataBasedSQL, applicationID, applicationVersionID, enterpriseID)
	if err != nil {
		repo.lgr.Error("addLatestApplicationVersionMetadataBasedTx: execute sql query")

		repo.txRollback(ctx, tx)

		return err
	}

	return nil
}
