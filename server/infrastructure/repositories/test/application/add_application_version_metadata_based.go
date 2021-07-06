package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/internal/models"
)

func (repo *testRepository) AddApplicationVersionMetadataBased(ctx context.Context, metadata *models.ApplicationMetadata, applicationUUID, enterpriseID string) (string, error) {
	return "", nil
}
