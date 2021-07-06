package application

import "context"

func (repo *testRepository) DeleteVersionApplication(ctx context.Context, applicationUUID, versionUUID, enterpriseID string) (string, error) {
	return "", nil
}
