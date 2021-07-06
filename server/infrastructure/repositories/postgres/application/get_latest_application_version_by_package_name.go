package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/infrastructure/repositories/postgres/application/sql/get"
	"github.com/jackc/pgx/v4"
)

func (repo *repository) GetLatestApplicationVersionByPackageName(ctx context.Context, packageName string) (int, error) {
	var version int

	row := repo.client.QueryRow(ctx, get.LatestApplicationNameByPackageNameSQL, packageName)

	err := row.Scan(&version)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return version, nil
		default:
			repo.lgr.Error("GetLatestApplicationVersionByPackageName: row scan error")

			return version, err
		}
	}

	return version, nil
}
