package application

import "context"

func (repo *testRepository) CheckExistApplicationByMetadata(ctx context.Context, applicationMetadataID int) (string, int, error) {
	return "", 0, nil
}
