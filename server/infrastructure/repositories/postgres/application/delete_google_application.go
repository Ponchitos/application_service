package application

import (
	"context"
)

func (repo *repository) DeleteGoogleApplicationByVersionUUID(ctx context.Context, applicationUUID, versionUUID, enterpriseID string) error {
	tx, err := repo.client.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = repo.deleteVersionApplicationTx(ctx, tx, applicationUUID, versionUUID, enterpriseID)
	if err != nil {
		return err
	}

	err = repo.deleteLatestVersionTx(ctx, tx, applicationUUID, enterpriseID)
	if err != nil {
		return err
	}

	err = repo.deleteApplicationTx(ctx, tx, applicationUUID, enterpriseID)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil

}
