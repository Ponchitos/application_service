package applications

import (
	"context"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Ponchitos/application_service/server/internal/models"
)

const completeUploadFileTopic = "application_service.upload_file"

func (app *application) CompleteUploadFile(ctx context.Context, apkUUID, enterpriseID string) (string, string, error) {
	app.lgr.Info("ApplicationService: CompleteUploadFile - received request")

	var applicationVersionUUID string

	applicationMetadata, err := app.store.GetApplicationMetadata(ctx, apkUUID)
	if err != nil && applicationMetadata == nil {
		app.lgr.Errorf("ApplicationService: CompleteUploadFile [GetApplicationMetadata] - %v", err)

		return "", "", errors.NewErrorf("Cannot get application metadata: %v", "Не удалось получить метаданные приложения: %v", err)
	}

	applicationUUID, latestVersion, err := app.store.CheckExistApplicationByMetadata(ctx, applicationMetadata.ID)
	if err != nil {
		app.lgr.Errorf("ApplicationService: CompleteUploadFile [GetApplicationMetadata] - %v", err)

		return "", "", errors.NewErrorf("Cannot check exist application: %v", "Не удалось проверить наличия приложения: %v", err)
	}

	if latestVersion == 0 && len(applicationUUID) > 0 {
		app.lgr.Error("ApplicationService: CompleteUploadFile [check exist from google play]")

		err := app.CancelApkUpload(ctx, apkUUID)
		if err != nil {
			return "", "", err
		}

		return "", "", errors.NewError("An application with this packageName is already exists. Change packageName.", "Приложение с таким packageName уже имеется в системе. Смените packageName для загрузки.")
	}

	if latestVersion != 0 && applicationMetadata.VersionCode <= latestVersion {
		app.lgr.Error("ApplicationService: CompleteUploadFile [check version]")

		return "", "", errors.NewError("Version of the apk is lower than the existing one", "Версия apk ниже уже существующей")
	}

	if len(applicationUUID) == 0 {
		applicationUUID, applicationVersionUUID, err = app.addBasicApplication(ctx, applicationMetadata, enterpriseID)
		if err != nil {
			app.lgr.Errorf("ApplicationService: CompleteUploadFile [addBasicApplication] - %v", err)

			return "", "", errors.NewErrorf("Cannot create application: %v", "Не удалось добавить приложение: %v", err)
		}
	} else {
		applicationVersionUUID, err = app.store.AddApplicationVersionMetadataBased(ctx, applicationMetadata, applicationUUID, enterpriseID)
		if err != nil {
			app.lgr.Errorf("ApplicationService: CompleteUploadFile [AddApplicationVersionMetadataBased] - %v", err)

			return "", "", errors.NewErrorf("Cannot increase application version: %v", "Не удалось увеличить версию приложения: %v", err)
		}
	}

	err = app.broker.CompleteUploadFileRecord(ctx, completeUploadFileTopic, applicationMetadata.PackageName)
	if err != nil {
		app.lgr.Errorf("ApplicationService: CompleteUploadFile [CompleteUploadFileRecord] - %v", err)
	}

	app.lgr.Info("ApplicationService: CompleteUploadFile - handle completed")

	return applicationUUID, applicationVersionUUID, nil
}

func (app *application) addBasicApplication(ctx context.Context, applicationMetadata *models.ApplicationMetadata, enterpriseID string) (string, string, error) {
	applicationUUID, applicationVersionUUID, err := app.store.AddBasicInfoApplicationMetadataBased(ctx, applicationMetadata, models.Private, models.Internal, enterpriseID)

	if err != nil {
		return "", "", err
	}

	return applicationUUID, applicationVersionUUID, nil
}
