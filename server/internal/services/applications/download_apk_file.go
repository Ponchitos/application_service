package applications

import (
	"context"
	"fmt"
	"github.com/Ponchitos/application_service/server/config"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Ponchitos/application_service/server/internal/handlers"
)

func (app *application) DownloadApkFile(ctx context.Context, applicationUUID, versionUUID, enterpriseID string) ([]byte, string, error) {
	app.lgr.Info("ApplicationService: DownloadApkFile - received request")

	var apkFile []byte

	metadata, err := app.store.GetApplicationMetadataByVersionUUID(ctx, applicationUUID, versionUUID, enterpriseID)
	if err != nil {
		app.lgr.Errorf("ApplicationService: DownloadApkFile [GetApplicationMetadataByVersionUUID] - %v", err)

		return nil, "", errors.NewErrorf("Cannot get application metadata: %s", "Не удалось получить метаданные: %s", err)
	}

	if app.config.Env == config.Local {
		apkFile, err = app.apkHandler.DownloadApkFile(ctx, fmt.Sprintf("%s.apk", metadata.UUID), handlers.LocalMode)
	} else {
		apkFile, err = app.apkHandler.DownloadApkFile(ctx, fmt.Sprintf("%s.apk", metadata.UUID), handlers.ExternalMode)
	}

	if err != nil {
		app.lgr.Errorf("ApplicationService: DownloadApkFile [DownloadApkFile] - %v", err)

		return nil, "", errors.NewErrorf("Cannot download file from store: %v", "Не удалось скачать файл из хранилища: %v", err)
	}

	app.lgr.Info("ApplicationService: DownloadApkFile - handle completed")

	return apkFile, fmt.Sprintf("%s.apk", metadata.UUID), nil
}
