package applications

import (
	"context"
	"github.com/Ponchitos/application_service/server/internal/models"
)

func (app *application) UninstallApplicationComplete(ctx context.Context, applicationUUID, versionUUID, enterpriseID string) error {
	app.lgr.Info("ApplicationService: UninstallApplicationComplete - received request")

	location, _, err := app.store.CheckDeleteApplicationAvailable(ctx, applicationUUID, versionUUID, enterpriseID)
	if err != nil {
		app.lgr.Errorf("ApplicationService: UninstallApplicationComplete [CheckDeleteApplicationAvailable] - %v", err)

		return err
	}

	if location == models.Google {
		err = app.googleApplicationDeleteProcessHandler(ctx, applicationUUID, versionUUID, enterpriseID)
		if err != nil {
			app.lgr.Errorf("ApplicationService: UninstallApplicationComplete [googleApplicationDeleteProcessHandler] - %v", err)

			return err
		}
	} else {
		err := app.store.UpdateApplicationStatusByVersionUUID(ctx, versionUUID, enterpriseID, models.Approved)
		if err != nil {
			app.lgr.Errorf("ApplicationService: UninstallApplicationComplete [UpdateApplicationStatusByVersionUUID] - %v", err)

			return err
		}
	}

	app.lgr.Info("ApplicationService: UninstallApplicationComplete - handle completed")

	return nil
}

func (app *application) googleApplicationDeleteProcessHandler(ctx context.Context, applicationUUID, versionUUID, enterpriseID string) error {
	googleApplicationInfo, err := app.store.GetGoogleApplicationByVersionUUID(ctx, applicationUUID, versionUUID, enterpriseID)
	if err != nil {
		return err
	}

	if googleApplicationInfo.Status == models.UnapprovedGoogleStatus {
		return app.store.DeleteGoogleApplicationByVersionUUID(ctx, applicationUUID, versionUUID, enterpriseID)
	}

	return app.store.UpdateApplicationStatusByVersionUUID(ctx, versionUUID, enterpriseID, models.Approved)
}
