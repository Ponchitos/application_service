package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/internal/models"
)

func (repo *testRepository) GetGoogleApplicationByVersionUUID(ctx context.Context, applicationUUID, versionUUID, enterpriseID string) (*models.GoogleApplication, error) {
	return nil, nil
}
