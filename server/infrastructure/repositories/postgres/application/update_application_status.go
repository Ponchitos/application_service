package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/infrastructure/repositories/postgres/application/sql/edit"
)

func (repo *repository) UpdateApplicationStatusByVersionUUID(ctx context.Context, versionUUID, enterpriseID, status string) error {
	_, err := repo.client.Exec(ctx, edit.ApplicationStatusByVersionUUIDSQL, status, versionUUID, enterpriseID)

	return err
}
