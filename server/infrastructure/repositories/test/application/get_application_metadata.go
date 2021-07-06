package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/internal/models"
)

func (repo *testRepository) GetApplicationMetadata(ctx context.Context, metadataUUID string) (*models.ApplicationMetadata, error) {
	return nil, nil
}
