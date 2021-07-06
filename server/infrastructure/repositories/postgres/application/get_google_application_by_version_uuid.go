package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/infrastructure/repositories/postgres/application/sql/get"
	"github.com/Ponchitos/application_service/server/internal/models"
)

func (repo *repository) GetGoogleApplicationByVersionUUID(ctx context.Context, applicationUUID, versionUUID, enterpriseID string) (*models.GoogleApplication, error) {
	var result *models.GoogleApplication

	row := repo.client.QueryRow(ctx, get.GoogleApplicationInfoByVersionUUIDSQL, versionUUID, applicationUUID, enterpriseID)

	err := row.Scan(&result)

	return result, err
}
