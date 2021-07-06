package applications

import (
	"context"
	"fmt"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Ponchitos/application_service/server/internal/handlers"
	"github.com/Ponchitos/application_service/server/internal/models"
)

const deleteVersionApplicationTopic = "application_service.delete_version_application"

func (app *application) DeleteVersionApplication(ctx context.Context, applicationUUID, versionUUID, enterpriseID string) error {
	app.lgr.Info("ApplicationService: DeleteVersionApplication - received request")

	location, status, err := app.store.CheckDeleteApplicationAvailable(ctx, applicationUUID, versionUUID, enterpriseID)
	if err != nil {
		app.lgr.Errorf("ApplicationService: DeleteVersionApplication [CheckDeleteApplicationAvailable] - %v", err)

		return err
	}

	if location == models.Google {
		app.lgr.Error("ApplicationService: DeleteVersionApplication [check application location]")

		return errors.NewError("Cannot delete an application with type location", "Приложение с таким типом нельзя удалить.")
	}

	if status != models.Approved {
		app.lgr.Error("ApplicationService: DeleteVersionApplication [check application status]")

		return errors.NewError("Removal is not possible. You must first uninstall the application from all devices.", "Удаление невозможно. Необходимо предварительно удалить приложение со всех устройств.")
	}

	applicationMetadataUUID, err := app.store.DeleteVersionApplication(ctx, applicationUUID, versionUUID, enterpriseID)
	if err != nil {
		app.lgr.Errorf("ApplicationService: DeleteVersionApplication [DeleteVersionApplication] - %v", err)

		return errors.NewErrorf("Cannot delete application: %v", "Не удалось удалить приложение: %v", err)
	}

	err = app.apkHandler.DeleteApkFile(ctx, fmt.Sprintf("%s.apk", applicationMetadataUUID), handlers.ExternalMode)
	if err != nil {
		app.lgr.Errorf("ApplicationService: DeleteVersionApplication [DeleteApkFile] - %v", err)
	}

	err = app.broker.DeleteApplicationVersionRecord(ctx, deleteVersionApplicationTopic, versionUUID)
	if err != nil {
		app.lgr.Errorf("ApplicationService: DeleteApplicationVersionRecord [DeleteApplicationVersionRecord] - %v", err)
	}

	app.lgr.Info("ApplicationService: DeleteVersionApplication - handle completed")

	return nil
}
