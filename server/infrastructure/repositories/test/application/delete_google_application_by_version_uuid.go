package application

import "context"

func (repo *testRepository) DeleteGoogleApplicationByVersionUUID(ctx context.Context, applicationUUID, versionUUID, enterpriseID string) error {
	return nil
}
