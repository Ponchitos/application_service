package applications

import (
	"context"
	"fmt"
	"github.com/Ponchitos/application_service/server/config"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Ponchitos/application_service/server/internal/handlers"
	"github.com/google/uuid"
	"mime/multipart"
)

func (app *application) UploadApkFile(ctx context.Context, apk multipart.File) (string, error) {
	app.lgr.Info("ApplicationService: UploadApkFile - received request")

	defer func() {
		if err := apk.Close(); err != nil {
			app.lgr.Errorf("Cannot close io stream: %v", err)
		}
	}()

	applicationName := uuid.New().String()
	fullApplicationFileName := fmt.Sprintf("%s.apk", applicationName)

	filePath, err := app.apkHandler.ApkFileSave(ctx, apk, fullApplicationFileName, handlers.LocalMode)
	if err != nil {
		app.lgr.Errorf("ApplicationService: UploadApkFile [ApkFileSave] - %v", err)

		return "", errors.NewErrorf("Cannot local save apk file: %v", "Не удалось локально сохранить apk файл: %v", err)
	}

	applicationMetadata, err := app.apkHandler.GetMetadataOfAPK(filePath, applicationName)
	if err != nil {
		app.lgr.Errorf("ApplicationService: UploadApkFile [GetMetadataOfAPK] - %v", err)

		if err := app.apkHandler.DeleteApkFile(ctx, filePath, handlers.LocalMode); err != nil {
			app.lgr.Errorf("ApplicationService: UploadApkFile [DeleteApkFile] - %v", err)
		}

		return "", errors.NewErrorf("Cannot get apk metafile: %v", "Не удалось получить метаданные apk файла: %v", err)
	}

	err = app.checkVersion(ctx, applicationMetadata.PackageName, applicationMetadata.VersionCode)
	if err != nil {
		app.lgr.Errorf("ApplicationService: UploadApkFile [checkVersion] - %v", err)

		if err := app.apkHandler.DeleteApkFile(ctx, filePath, handlers.LocalMode); err != nil {
			app.lgr.Errorf("ApplicationService: UploadApkFile [DeleteApkFile] - %v", err)
		}

		return "", err
	}

	linkPath, err := app.apkHandler.ApkFileSave(ctx, apk, fullApplicationFileName, handlers.ExternalMode)
	if err != nil {
		app.lgr.Errorf("ApplicationService: UploadApkFile [ApkFileSave] - %v", err)

		if err := app.apkHandler.DeleteApkFile(ctx, filePath, handlers.LocalMode); err != nil {
			app.lgr.Errorf("ApplicationService: UploadApkFile [DeleteApkFile] - %v", err)
		}

		return "", errors.NewErrorf("Cannot external save apk file: %v", "Не удалось сохранить apk файл на стороннем хранилище: %v", err)
	}

	applicationMetadata.Link = linkPath

	if app.config.Env != config.Local {
		err = app.apkHandler.DeleteApkFile(ctx, filePath, handlers.LocalMode)
		if err != nil {
			app.lgr.Errorf("ApplicationService: UploadApkFile [DeleteApkFile] - %v", err)
		}
	}

	err = app.store.AddApplicationMetadata(ctx, applicationMetadata)
	if err != nil {
		app.lgr.Errorf("ApplicationService: UploadApkFile [AddApplicationMetadata] - %v", err)

		err = app.apkHandler.DeleteApkFile(ctx, fullApplicationFileName, handlers.ExternalMode)
		if err != nil {
			app.lgr.Errorf("ApplicationService: UploadApkFile [DeleteApkFile] - %v", err)
		}

		return "", errors.NewErrorf("Cannot save metadata of apk file to database: %v", "Не удалось сохранить метаданные apk файла в базе данных: %v", err)
	}

	app.lgr.Info("ApplicationService: UploadApkFile - handle completed")

	return applicationMetadata.UUID, nil
}

func (app *application) checkVersion(ctx context.Context, packageName string, version int) error {
	latestVersion, err := app.store.GetLatestApplicationVersionByPackageName(ctx, packageName)
	if err != nil {
		return errors.NewErrorf("Cannot get application latest version: %v", "Не удалось получить последнюю версию приложения: %v", err)
	}

	if version <= latestVersion {
		return errors.NewError("Version of the apk is lower than the existing one", "Версия apk ниже уже существующей")
	}

	return nil
}
