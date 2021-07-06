package applications

import "context"

func (app *application) ChangeApplicationStatus(ctx context.Context, versionUUID, enterpriseID, status string) error {
	app.lgr.Info("ApplicationService: ChangeApplicationStatus - received request")

	err := app.store.UpdateApplicationStatusByVersionUUID(ctx, versionUUID, enterpriseID, status)
	if err != nil {
		app.lgr.Errorf("ApplicationService: ChangeApplicationStatus [UpdateApplicationStatusByVersionUUID] - %v", err)

		return err
	}

	app.lgr.Info("ApplicationService: ChangeApplicationStatus - handle completed")

	return nil
}
