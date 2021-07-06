package google

import (
	"context"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Ponchitos/application_service/server/internal/models"
	"github.com/google/uuid"
	"google.golang.org/api/androidmanagement/v1"
)

func (e *external) GetGoogleApplication(ctx context.Context, enterpriseID, packageName string) (*models.GoogleApplication, error) {
	options, err := e.getOptions(ctx, enterpriseID)

	if err != nil {
		return nil, errors.NewErrorf("Cannot get google options: %v", "Не удалось получить google опции: %v", err)
	}

	service, err := androidmanagement.NewService(ctx, options)
	if err != nil {
		return nil, errors.NewErrorf("Cannot init google service: %v", "Не удалось инициализировать google сервис: %v", err)
	}

	googleApplication, err := service.Enterprises.Applications.Get(e.initApplicationFullName(enterpriseID, packageName)).Do()
	if err != nil {
		return nil, errors.NewErrorf("Cannot get app from google play: %v", "Не удалось получить приложение из google play: %v", err)
	}

	return e.convertGoogleApplicationToInternalApplication(googleApplication), nil
}

func (e *external) convertGoogleApplicationToInternalApplication(googleApplication *androidmanagement.Application) *models.GoogleApplication {
	result := &models.GoogleApplication{
		UUID:              uuid.New().String(),
		Name:              e.parserApplicationName(googleApplication.Name),
		Title:             googleApplication.Title,
		Permissions:       e.convertGooglePermissions(googleApplication.Permissions),
		AppTracks:         e.convertGoogleAppTracks(googleApplication.AppTracks),
		ManagedProperties: e.convertGoogleManagedProperty(googleApplication.ManagedProperties),
	}

	return result
}

func (e *external) convertGooglePermissions(googlePermissions []*androidmanagement.ApplicationPermission) []*models.ApplicationPermission {
	var response []*models.ApplicationPermission

	for _, permission := range googlePermissions {
		if permission != nil {
			response = append(response, &models.ApplicationPermission{
				PermissionID: permission.PermissionId,
				Name:         permission.Name,
				Description:  permission.Description,
			})
		}
	}

	return response
}

func (e *external) convertGoogleAppTracks(googleAppTracks []*androidmanagement.AppTrackInfo) []*models.ApplicationTrack {
	var response []*models.ApplicationTrack

	for _, track := range googleAppTracks {
		if track != nil {
			response = append(response, &models.ApplicationTrack{
				TrackID:    track.TrackId,
				TrackAlias: track.TrackAlias,
			})
		}
	}

	return response
}

func (e *external) convertGoogleManagedProperty(googleManagedProperty []*androidmanagement.ManagedProperty) []*models.ManagedProperty {
	var response []*models.ManagedProperty

	for _, property := range googleManagedProperty {
		if property != nil {
			response = append(response, &models.ManagedProperty{
				Key:              property.Key,
				Type:             property.Type,
				Title:            property.Title,
				Description:      property.Description,
				DefaultValue:     property.DefaultValue,
				Entries:          e.convertGoogleEntries(property.Entries),
				NestedProperties: e.convertGoogleManagedProperty(property.NestedProperties),
			})
		}
	}

	return response
}

func (e *external) convertGoogleEntries(googleEntries []*androidmanagement.ManagedPropertyEntry) []*models.ManagedPropertyEntry {
	var response []*models.ManagedPropertyEntry

	for _, entry := range googleEntries {
		if entry != nil {
			response = append(response, &models.ManagedPropertyEntry{
				Value: entry.Value,
				Name:  entry.Name,
			})
		}
	}

	return response
}
