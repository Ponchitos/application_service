package repositories

import (
	"context"
	"github.com/Ponchitos/application_service/server/internal/models"
)

type ApplicationRepository interface {
	AddApplicationMetadata(ctx context.Context, metadata *models.ApplicationMetadata) error

	GetApplicationMetadata(ctx context.Context, metadataUUID string) (*models.ApplicationMetadata, error)

	DeleteApplicationMetadata(ctx context.Context, metadataUUID string) error

	AddBasicInfoApplicationMetadataBased(ctx context.Context, metadata *models.ApplicationMetadata, available, location, enterpriseID string) (applicationUUID string, applicationVersionUUID string, err error)

	GetBasicInfoApplications(ctx context.Context, enterpriseID string, offset, limit int) (count int, applications []*models.BasicApplication, err error)

	GetApplicationVersions(ctx context.Context, applicationUUID, enterpriseID string) ([]*models.BasicApplication, error)

	AddApplicationVersionMetadataBased(ctx context.Context, metadata *models.ApplicationMetadata, applicationUUID, enterpriseID string) (applicationVersionUUID string, err error)

	GetApplication(ctx context.Context, applicationUUID, applicationVersionUUID, enterpriseID string) (*models.Application, error)

	CheckExistApplicationByMetadata(ctx context.Context, applicationMetadataID int) (applicationUUID string, latestVersion int, err error)

	CheckDeleteApplicationAvailable(ctx context.Context, applicationUUID, versionUUID, enterpriseID string) (location, status string, err error)

	GetLatestApplicationVersionByPackageName(ctx context.Context, packageName string) (version int, err error)

	DeleteVersionApplication(ctx context.Context, applicationUUID, versionUUID, enterpriseID string) (metadataUUID string, err error)

	UpdateApplicationStatusByVersionUUID(ctx context.Context, versionUUID, enterpriseID, status string) error

	DeleteGoogleApplicationByVersionUUID(ctx context.Context, applicationUUID, versionUUID, enterpriseID string) error

	GetApplicationMetadataByVersionUUID(ctx context.Context, applicationUUID, versionUUID, enterpriseID string) (*models.ApplicationMetadata, error)

	GetGoogleApplicationByVersionUUID(ctx context.Context, applicationUUID, versionUUID, enterpriseID string) (*models.GoogleApplication, error)
}
