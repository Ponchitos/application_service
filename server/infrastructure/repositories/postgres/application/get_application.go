package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/infrastructure/repositories/postgres/application/sql/get"
	"github.com/Ponchitos/application_service/server/internal/models"
	"github.com/jackc/pgx/v4"
)

type application struct {
	VersionUUID     string `json:"version_uuid"`
	ApplicationUUID string `json:"application_uuid"`

	PackageName  string `json:"package_name"`
	EnterpriseID string `json:"enterprise_id"`
	Available    string `json:"available"`
	Location     string `json:"location"`

	Metadata *models.ApplicationMetadata `json:"metadata"`
	Google   *models.GoogleApplication   `json:"google"`
}

func (repo *repository) GetApplication(ctx context.Context, applicationUUID, applicationVersionUUID, enterpriseID string) (*models.Application, error) {
	applicationInfo, err := repo.getApplication(ctx, applicationUUID, applicationVersionUUID, enterpriseID)
	if err != nil {
		return nil, err
	}

	if applicationInfo.Metadata != nil && len(applicationInfo.Metadata.UUID) > 0 {
		applicationInfo.Metadata, err = repo.GetApplicationMetadata(ctx, applicationInfo.Metadata.UUID)
		if err != nil {
			return nil, err
		}
	} else {
		applicationInfo.Metadata = nil
	}

	if applicationInfo.Google != nil && len(applicationInfo.Google.UUID) > 0 {
		applicationInfo.Google, err = repo.GetGoogleApplication(ctx, applicationInfo.Google.UUID)
		if err != nil {
			return nil, err
		}
	} else {
		applicationInfo.Google = nil
	}

	if applicationInfo == nil {
		return nil, nil
	}

	applicationModel := repo.convertToApplicationModel(applicationInfo)

	return applicationModel, nil
}

func (repo *repository) getApplication(ctx context.Context, applicationUUID, applicationVersionUUID, enterpriseID string) (*application, error) {
	var result *application

	row := repo.client.QueryRow(ctx, get.ApplicationSQL, applicationUUID, enterpriseID, applicationVersionUUID)

	err := row.Scan(&result)

	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, nil
		default:
			repo.lgr.Error("getApplication: row scan error")

			return nil, err
		}
	}

	return result, nil
}

func (repo *repository) convertToApplicationModel(application *application) *models.Application {
	result := &models.Application{
		ApplicationUUID: application.ApplicationUUID,
		VersionUUID:     application.VersionUUID,
		EnterpriseID:    application.EnterpriseID,
		Availability: models.Availability{
			Location:  application.Location,
			Available: application.Available,
		},
		UsesFeatures: make([]string, 0),
		Permissions:  make([]string, 0),

		ApplicationSettings: make([]*models.ManagedProperty, 0),
	}

	if application.Metadata != nil {
		result.UsesFeatures = application.Metadata.UsesFeatures
		result.Permissions = application.Metadata.GetPermissionsAsStrings()
		result.Icon = application.Metadata.IconBase64
		result.ShortInfo.Name = application.Metadata.ApplicationLabel
		result.ShortInfo.PackageName = application.PackageName
		result.ShortInfo.VersionCode = application.Metadata.VersionCode
		result.ShortInfo.VersionName = application.Metadata.VersionName
		result.ShortInfo.MinSDK = application.Metadata.MinimumSDK
	} else if application.Google != nil {
		result.Permissions = application.Google.GetPermissionsAsString()
		result.ShortInfo.Name = application.Google.Name
		result.ShortInfo.PackageName = application.Google.Name
		result.ApplicationSettings = application.Google.ManagedProperties
	}

	return result
}
