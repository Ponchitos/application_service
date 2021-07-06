package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/infrastructure/repositories/postgres/application/sql/get"
	"github.com/Ponchitos/application_service/server/internal/models"
	"github.com/jackc/pgx/v4"
)

func (repo *repository) GetApplicationVersions(ctx context.Context, applicationUUID, enterpriseID string) ([]*models.BasicApplication, error) {
	var result []*models.BasicApplication

	row := repo.client.QueryRow(ctx, get.ApplicationVersionsSQL, applicationUUID, enterpriseID)

	err := row.Scan(&result)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return make([]*models.BasicApplication, 0), nil
		default:
			repo.lgr.Error("GetApplicationVersions: row scan error")

			return nil, err
		}
	}

	return result, nil
}
