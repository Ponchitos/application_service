package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/infrastructure/repositories/postgres/application/sql/add"
	"github.com/Ponchitos/application_service/server/internal/models"
)

func (repo *repository) AddApplicationMetadata(ctx context.Context, metadata *models.ApplicationMetadata) error {
	permissions, err := metadata.ConvertPermissionsToBytes()
	if err != nil {
		return err
	}

	_, err = repo.client.Exec(ctx, add.ApplicationMetadataSQL,
		metadata.UUID,
		metadata.Link,
		metadata.PackageName,
		metadata.ApplicationLabel,
		metadata.VersionName,
		metadata.FileSize,
		metadata.FileSha1Base64,
		metadata.FileSha256Base64,
		metadata.IconBase64,
		metadata.ExternallyHostedURL,
		metadata.NativeCodes,
		metadata.CertificateBase64s,
		metadata.UsesFeatures,
		metadata.VersionCode,
		metadata.MinimumSDK,
		permissions,
	)

	return err
}
