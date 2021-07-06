package externalapi

import (
	"context"
	"github.com/Ponchitos/application_service/server/internal/models"
)

type ExternalAPI interface {
	GetGoogleApplication(ctx context.Context, enterpriseID, packageName string) (*models.GoogleApplication, error)
}
