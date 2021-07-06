package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/internal/models"
)

func (repo *testRepository) GetApplicationVersions(ctx context.Context, applicationUUID, enterpriseID string) ([]*models.BasicApplication, error) {
	return nil, nil
}
