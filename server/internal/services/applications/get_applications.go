package applications

import (
	"context"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Ponchitos/application_service/server/internal/models"
)

func (app *application) GetApplications(ctx context.Context, enterpriseID string, offset, limit int) (int, []*models.BasicApplication, error) {
	app.lgr.Info("ApplicationService: GetApplications - received request")

	count, basicApplications, err := app.store.GetBasicInfoApplications(ctx, enterpriseID, offset, limit)
	if err != nil {
		app.lgr.Errorf("ApplicationService: GetApplicationInfo [GetApplication] - %v", err)

		return 0, nil, errors.NewErrorf("Cannot get applications: %v", "Не удалось загрузить список приложений: %v", err)
	}

	app.lgr.Info("ApplicationService: GetApplications - handle completed")

	return count, basicApplications, nil
}
