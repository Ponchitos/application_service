package applications

import (
	"context"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Ponchitos/application_service/server/internal/models"
)

func (app *application) GetApplicationVersions(ctx context.Context, applicationUUID, enterpriseID string) ([]*models.BasicApplication, error) {
	app.lgr.Info("ApplicationService: GetApplicationVersions - received request")

	basicApplications, err := app.store.GetApplicationVersions(ctx, applicationUUID, enterpriseID)
	if err != nil {
		app.lgr.Errorf("ApplicationService: GetApplicationVersions [GetApplicationVersions] - %v", err)

		return nil, errors.NewErrorf("Cannot get application versions: %v", "Не удалось загрузить список версий приложения: %v", err)
	}

	app.lgr.Info("ApplicationService: GetApplicationVersions - handle completed")

	return basicApplications, nil
}
