package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/infrastructure/repositories/postgres/application/sql/get"
	"github.com/jackc/pgx/v4"
)

func (repo *repository) CheckExistApplicationByMetadata(ctx context.Context, applicationMetadataID int) (string, int, error) {
	var (
		applicationUUID string
		latestVersion   int
		err             error
	)

	row := repo.client.QueryRow(ctx, get.CheckExistApplicationByMetadataSQL, applicationMetadataID)

	err = row.Scan(&applicationUUID, &latestVersion)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return applicationUUID, 0, nil
		default:
			repo.lgr.Error("CheckExistApplicationByMetadata: row scan error")

			return applicationUUID, 0, err
		}
	}

	return applicationUUID, latestVersion, err
}
