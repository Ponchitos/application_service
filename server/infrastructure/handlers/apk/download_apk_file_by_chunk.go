package apk

import "context"

func (hAPK *apkHandler) DownloadApkFileByChunk(ctx context.Context, key string, offset, whence int64, mode string) ([]byte, int64, error) {
	store, err := hAPK.getStore(mode)
	if err != nil {
		return nil, 0, err
	}

	return store.GetByChunk(ctx, key, offset, whence)
}
