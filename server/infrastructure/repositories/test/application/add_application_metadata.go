package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/internal/models"
)

func (repo *testRepository) AddApplicationMetadata(ctx context.Context, metadata *models.ApplicationMetadata) error {
	return nil
}
