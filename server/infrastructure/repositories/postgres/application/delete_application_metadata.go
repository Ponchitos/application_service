package application

import (
	"context"
	deleteQuery "github.com/Ponchitos/application_service/server/infrastructure/repositories/postgres/application/sql/delete"
)

func (repo *repository) DeleteApplicationMetadata(ctx context.Context, metadataUUID string) error {
	_, err := repo.client.Exec(ctx, deleteQuery.ApplicationMetadataByUUIDSQL, metadataUUID)

	return err
}
