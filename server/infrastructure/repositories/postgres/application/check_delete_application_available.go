package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Ponchitos/application_service/server/infrastructure/repositories/postgres/application/sql/get"
	"github.com/jackc/pgx/v4"
)

func (repo *repository) CheckDeleteApplicationAvailable(ctx context.Context, applicationUUID, versionUUID, enterpriseID string) (string, string, error) {
	var location, status string

	row := repo.client.QueryRow(ctx, get.CheckDeleteApplicationAvailableSQL, versionUUID, applicationUUID, enterpriseID)

	err := row.Scan(&location, &status)

	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return location, status, errors.NewError("Application not found", "Приложение не найдено")
		default:
			repo.lgr.Error("CheckDeleteApplicationAvailable: row scan error")

			return location, status, err
		}
	}

	return location, status, nil
}
