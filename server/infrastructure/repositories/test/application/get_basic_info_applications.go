package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/internal/models"
)

func (repo *testRepository) GetBasicInfoApplications(ctx context.Context, enterpriseID string, offset, limit int) (int, []*models.BasicApplication, error) {
	return 0, nil, nil
}
