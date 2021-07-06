package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/infrastructure/repositories/postgres/application/sql/get"
	"github.com/Ponchitos/application_service/server/internal/models"
	"github.com/jackc/pgx/v4"
)

func (repo *repository) GetApplicationMetadata(ctx context.Context, metadataUUID string) (*models.ApplicationMetadata, error) {
	var (
		applicationMetadata models.ApplicationMetadata
		permissions         []byte
	)

	row := repo.client.QueryRow(ctx, get.ApplicationMetadataSQL, metadataUUID)

	err := row.Scan(
		&applicationMetadata.UUID,
		&applicationMetadata.ID,
		&applicationMetadata.Link,
		&applicationMetadata.PackageName,
		&applicationMetadata.ApplicationLabel,
		&applicationMetadata.VersionName,
		&applicationMetadata.FileSize,
		&applicationMetadata.FileSha1Base64,
		&applicationMetadata.FileSha256Base64,
		&applicationMetadata.IconBase64,
		&applicationMetadata.ExternallyHostedURL,
		&applicationMetadata.NativeCodes,
		&applicationMetadata.CertificateBase64s,
		&applicationMetadata.UsesFeatures,
		&applicationMetadata.VersionCode,
		&applicationMetadata.MinimumSDK,
		&applicationMetadata.Created,
		&permissions,
	)

	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, nil
		default:
			repo.lgr.Error("GetApplicationMetadata: row scan error")

			return nil, err
		}
	}

	err = applicationMetadata.ConvertBytesToPermissions(permissions)
	if err != nil {
		repo.lgr.Error("GetApplicationMetadata: convert")

		return nil, err
	}

	return &applicationMetadata, nil
}
