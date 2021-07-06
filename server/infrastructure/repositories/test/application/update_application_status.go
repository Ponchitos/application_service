package application

import "context"

func (repo *testRepository) UpdateApplicationStatusByVersionUUID(ctx context.Context, versionUUID, enterpriseID, status string) error {
	return nil
}
