package apk

import "context"

func (hAPK *apkHandler) DeleteApkFile(ctx context.Context, filePath, mode string) error {
	store, err := hAPK.getStore(mode)
	if err != nil {
		return err
	}

	return store.Delete(ctx, filePath)
}
