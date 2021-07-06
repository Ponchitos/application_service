package applications

import (
	"context"
	"fmt"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Ponchitos/application_service/server/internal/handlers"
)

func (app *application) CancelApkUpload(ctx context.Context, metadataUUID string) error {
	app.lgr.Info("ApplicationService: CancelApkUpload - received request")

	fullFileName := fmt.Sprintf("%s.apk", metadataUUID)

	err := app.store.DeleteApplicationMetadata(ctx, metadataUUID)
	if err != nil {
		app.lgr.Errorf("ApplicationService: CancelApkUpload [DeleteApplicationMetadata] - %v", err)

		return errors.NewErrorf("Cannot delete application metadata: %v", "Не удалось удалить метаданные приложения: %v", err)
	}

	err = app.apkHandler.DeleteApkFile(ctx, fullFileName, handlers.ExternalMode)
	if err != nil {
		app.lgr.Errorf("ApplicationService: CancelApkUpload [DeleteApkFile] - %v", err)

		return errors.NewErrorf("Cannot delete file from external store: %v", "Не удалось удалить файл из внешенго хранилища: %v", err)
	}

	app.lgr.Info("ApplicationService: CancelApkUpload - handle completed")

	return nil
}
