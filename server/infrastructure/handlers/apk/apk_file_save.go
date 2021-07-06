package apk

import (
	"context"
	"mime/multipart"
)

func (hAPK *apkHandler) ApkFileSave(ctx context.Context, apkFile multipart.File, name, mode string) (string, error) {
	store, err := hAPK.getStore(mode)
	if err != nil {
		return "", err
	}

	return store.Save(ctx, apkFile, name)
}
