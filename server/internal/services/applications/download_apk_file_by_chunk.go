package applications

import (
	"context"
	"fmt"
	"github.com/Ponchitos/application_service/server/config"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Ponchitos/application_service/server/internal/handlers"
)

func (app *application) DownloadApkFileByChunk(ctx context.Context, applicationUUID, versionUUID, enterpriseID string, offset, whence int64) ([]byte, int64, string, error) {
	app.lgr.Info("ApplicationService: DownloadApkFileByChunk - received request")

	var (
		chunk []byte
		size  int64
	)

	metadata, err := app.store.GetApplicationMetadataByVersionUUID(ctx, applicationUUID, versionUUID, enterpriseID)
	if err != nil {
		app.lgr.Errorf("ApplicationService: DownloadApkFileByChunk [GetApplicationMetadataByVersionUUID] - %v", err)

		return nil, 0, "", errors.NewErrorf("Cannot get application metadata: %s", "Не удалось получить метаданные: %s", err)
	}

	if app.config.Env == config.Local {
		chunk, size, err = app.apkHandler.DownloadApkFileByChunk(ctx, fmt.Sprintf("%s.apk", metadata.UUID), offset, whence, handlers.LocalMode)
	} else {
		chunk, size, err = app.apkHandler.DownloadApkFileByChunk(ctx, fmt.Sprintf("%s.apk", metadata.UUID), offset, whence, handlers.ExternalMode)
	}

	if err != nil {
		app.lgr.Errorf("ApplicationService: DownloadApkFileByChunk [DownloadApkFileByChunk] - %v", err)

		return nil, 0, "", errors.NewErrorf("Cannot download chunk from store: %v", "Не удалось скачать фрагмент из хранилища: %v", err)
	}

	app.lgr.Info("ApplicationService: DownloadApkFileByChunk - handle completed")

	return chunk, size, fmt.Sprintf("%s.apk", metadata.UUID), nil
}
