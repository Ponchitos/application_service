package application

import (
	"context"
	"database/sql"
	deleteQuery "github.com/Ponchitos/application_service/server/infrastructure/repositories/postgres/application/sql/delete"
	"github.com/Ponchitos/application_service/server/infrastructure/repositories/postgres/application/sql/get"
	"github.com/jackc/pgx/v4"
)

func (repo *repository) DeleteVersionApplication(ctx context.Context, applicationUUID, versionUUID, enterpriseID string) (string, error) {
	tx, err := repo.client.Begin(ctx)
	if err != nil {
		return "", err
	}

	applicationMetadataID, err := repo.deleteVersionApplicationTx(ctx, tx, applicationUUID, versionUUID, enterpriseID)
	if err != nil {
		return "", err
	}

	applicationMetadataUUID, err := repo.deleteApplicationMetadataTx(ctx, tx, applicationMetadataID)
	if err != nil {
		return "", err
	}

	latestVersionID, err := repo.checkApplicationLatestVersionExistTx(ctx, tx, applicationUUID, enterpriseID)
	if err != nil {
		return "", err
	}

	if latestVersionID > 0 {
		err := repo.updateLatestApplicationVersionTx(ctx, tx, applicationUUID, enterpriseID, latestVersionID)
		if err != nil {
			return "", err
		}
	} else {
		err := repo.deleteLatestVersionTx(ctx, tx, applicationUUID, enterpriseID)
		if err != nil {
			return "", err
		}

		err = repo.deleteApplicationTx(ctx, tx, applicationUUID, enterpriseID)
		if err != nil {
			return "", err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", err
	}

	return applicationMetadataUUID, nil
}

func (repo *repository) deleteVersionApplicationTx(ctx context.Context, tx pgx.Tx, applicationUUID, versionUUID, enterpriseID string) (int, error) {
	var applicationMetadataID sql.NullInt64

	row := tx.QueryRow(ctx, deleteQuery.VersionApplicationSQL, versionUUID, applicationUUID, enterpriseID)

	err := row.Scan(&applicationMetadataID)
	if err != nil {
		repo.lgr.Error("deleteVersionApplicationTx: row scan error")

		repo.txRollback(ctx, tx)

		return int(applicationMetadataID.Int64), err
	}

	return int(applicationMetadataID.Int64), nil
}

func (repo *repository) deleteApplicationMetadataTx(ctx context.Context, tx pgx.Tx, applicationMetadataID int) (string, error) {
	var applicationMetadataUUID string

	row := tx.QueryRow(ctx, deleteQuery.ApplicationMetadataByIDSQL, applicationMetadataID)

	err := row.Scan(&applicationMetadataUUID)
	if err != nil {
		repo.lgr.Error("deleteApplicationMetadataTx: row scan error")

		repo.txRollback(ctx, tx)

		return applicationMetadataUUID, err
	}

	return applicationMetadataUUID, nil
}

func (repo *repository) checkApplicationLatestVersionExistTx(ctx context.Context, tx pgx.Tx, applicationUUID, enterpriseID string) (int, error) {
	var latestVersionID int

	row := tx.QueryRow(ctx, get.CheckApplicationLatestVersionExistSQL, applicationUUID, enterpriseID)

	err := row.Scan(&latestVersionID)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return latestVersionID, nil
		default:
			repo.lgr.Error("checkApplicationLatestVersionExistTx: row scan error")

			repo.txRollback(ctx, tx)

			return latestVersionID, err
		}
	}

	return latestVersionID, nil
}

func (repo *repository) deleteLatestVersionTx(ctx context.Context, tx pgx.Tx, applicationUUID, enterpriseID string) error {
	_, err := tx.Exec(ctx, deleteQuery.LatestVersionSQL, applicationUUID, enterpriseID)
	if err != nil {
		repo.lgr.Error("deleteLatestVersionTx: execute sql query")

		repo.txRollback(ctx, tx)

		return err
	}

	return nil
}

func (repo *repository) deleteApplicationTx(ctx context.Context, tx pgx.Tx, applicationUUID, enterpriseID string) error {
	_, err := tx.Exec(ctx, deleteQuery.ApplicationSQL, applicationUUID, enterpriseID)
	if err != nil {
		repo.lgr.Error("deleteApplicationTx: execute sql query")

		repo.txRollback(ctx, tx)

		return err
	}

	return nil
}
