package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/internal/models"
)

func (repo *testRepository) GetApplication(ctx context.Context, applicationUUID, applicationVersionUUID, enterpriseID string) (*models.Application, error) {
	return nil, nil
}
