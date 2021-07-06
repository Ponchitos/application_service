package apk

import "context"

func (hAPK *apkHandler) DownloadApkFile(ctx context.Context, name, mode string) ([]byte, error) {
	store, err := hAPK.getStore(mode)
	if err != nil {
		return nil, err
	}

	return store.Get(ctx, name)
}
