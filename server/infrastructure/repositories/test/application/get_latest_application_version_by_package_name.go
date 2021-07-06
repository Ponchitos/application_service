package application

import "context"

func (repo *testRepository) GetLatestApplicationVersionByPackageName(ctx context.Context, packageName string) (int, error) {
	return 0, nil
}
