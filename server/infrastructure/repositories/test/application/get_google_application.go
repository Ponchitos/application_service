package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/internal/models"
)

func (repo *testRepository) GetGoogleApplication(ctx context.Context, googleApplicationUUID string) (*models.GoogleApplication, error) {
	return nil, nil
}
