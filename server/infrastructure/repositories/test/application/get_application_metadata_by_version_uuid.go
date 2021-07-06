package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/internal/models"
)

func (repo *testRepository) GetApplicationMetadataByVersionUUID(ctx context.Context, applicationUUID, versionUUID, enterpriseID string) (*models.ApplicationMetadata, error) {
	return nil, nil
}
