package application

import "context"

func (repo *testRepository) CheckDeleteApplicationAvailable(ctx context.Context, applicationUUID, versionUUID, enterpriseID string) (string, string, error) {
	return "", "", nil
}
