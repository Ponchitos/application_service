package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/internal/models"
)

func (repo *testRepository) AddBasicInfoApplicationMetadataBased(ctx context.Context, metadata *models.ApplicationMetadata, available, location, enterpriseID string) (string, string, error) {
	return "", "", nil
}
