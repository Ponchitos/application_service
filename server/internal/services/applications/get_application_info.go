package applications

import (
	"context"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Ponchitos/application_service/server/internal/models"
)

func (app *application) GetApplicationInfo(ctx context.Context, applicationUUID, applicationVersionUUID, enterpriseID string) (*models.Application, error) {
	app.lgr.Info("ApplicationService: GetApplicationInfo - received request")

	application, err := app.store.GetApplication(ctx, applicationUUID, applicationVersionUUID, enterpriseID)
	if err != nil {
		app.lgr.Errorf("ApplicationService: GetApplicationInfo [GetApplication] - %v", err)

		return nil, errors.NewErrorf("Cannot get application information: %v", "Не удалось получить информацию о приложении: %v", err)
	}

	app.lgr.Info("ApplicationService: GetApplicationInfo - handle completed")

	return application, nil
}
