package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/infrastructure/repositories/postgres/application/sql/get"
	"github.com/Ponchitos/application_service/server/internal/models"
	"github.com/jackc/pgx/v4"
)

func (repo *repository) GetBasicInfoApplications(ctx context.Context, enterpriseID string, offset, limit int) (int, []*models.BasicApplication, error) {
	var (
		count        int
		applications []*models.BasicApplication
	)

	row := repo.client.QueryRow(ctx, get.BasicApplicationsSQL, enterpriseID, offset, limit)

	err := row.Scan(&count, &applications)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return 0, make([]*models.BasicApplication, 0), nil
		default:
			repo.lgr.Error("GetBasicInfoApplications: row scan error")

			return 0, nil, err
		}
	}

	if len(applications) == 0 {
		return 0, make([]*models.BasicApplication, 0), nil
	}

	return count, applications, nil
}
