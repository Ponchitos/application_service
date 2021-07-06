package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/infrastructure/repositories/postgres/application/sql/get"
	"github.com/Ponchitos/application_service/server/internal/models"
	"github.com/jackc/pgx/v4"
)

func (repo *repository) GetGoogleApplication(ctx context.Context, googleApplicationUUID string) (*models.GoogleApplication, error) {
	var (
		result *models.GoogleApplication
		err    error
	)

	result, err = repo.getBasicGoogleApplicationInfo(ctx, googleApplicationUUID)
	if err != nil {
		return nil, err
	}

	result.Permissions, err = repo.getGoogleAppPermissions(ctx, result.ID)
	if err != nil {
		return nil, err
	}

	result.AppTracks, err = repo.getGoogleApplicationTracks(ctx, result.ID)
	if err != nil {
		return nil, err
	}

	result.ManagedProperties, err = repo.getGoogleApplicationManagedPropertiesByAppID(ctx, result.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *repository) getBasicGoogleApplicationInfo(ctx context.Context, googleApplicationUUID string) (*models.GoogleApplication, error) {
	var result models.GoogleApplication

	row := repo.client.QueryRow(ctx, get.GoogleApplicationSQL, googleApplicationUUID)

	err := row.Scan(&result.ID, &result.UUID, &result.Name, &result.Title)

	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, nil
		default:
			repo.lgr.Error("GetGoogleApplication: row scan error")

			return nil, err
		}
	}

	return &result, nil
}

func (repo *repository) getGoogleAppPermissions(ctx context.Context, googleAppID int) ([]*models.ApplicationPermission, error) {
	result := make([]*models.ApplicationPermission, 0)

	rows, err := repo.client.Query(ctx, get.GoogleAppPermissionsSQL, googleAppID)
	if err != nil {
		repo.lgr.Error("GetGoogleApplication: get rows error")

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var permission models.ApplicationPermission

		err := rows.Scan(&permission.Name, &permission.Description, &permission.PermissionID)
		if err != nil {
			switch err {
			case pgx.ErrNoRows:
				continue
			default:
				repo.lgr.Error("getGoogleAppPermissions: row scan error")

				return nil, err
			}
		}

		result = append(result, &permission)
	}

	return result, nil
}

func (repo *repository) getGoogleApplicationTracks(ctx context.Context, googleAppID int) ([]*models.ApplicationTrack, error) {
	result := make([]*models.ApplicationTrack, 0)

	rows, err := repo.client.Query(ctx, get.GoogleApplicationTracksSQL, googleAppID)
	if err != nil {
		repo.lgr.Error("getGoogleApplicationTracks: get rows error")

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var track models.ApplicationTrack

		err := rows.Scan(&track.TrackID, &track.TrackAlias)
		if err != nil {
			switch err {
			case pgx.ErrNoRows:
				continue
			default:
				repo.lgr.Error("getGoogleApplicationTracks: row scan error")

				return nil, err
			}
		}

		result = append(result, &track)
	}

	return result, nil
}

func (repo *repository) getGoogleApplicationManagedPropertiesByAppID(ctx context.Context, googleAppID int) ([]*models.ManagedProperty, error) {
	var err error

	result := make([]*models.ManagedProperty, 0)

	rows, err := repo.client.Query(ctx, get.GoogleApplicationManagedPropertyByAppIDSQL, googleAppID)
	if err != nil {
		repo.lgr.Error("getGoogleApplicationManagedPropertiesByAppID: get rows error")

		return nil, err
	}

	for rows.Next() {
		var (
			property models.ManagedProperty
			entries  []byte
		)

		err := rows.Scan(
			&property.ManagedPropertyID,
			&property.Key,
			&property.Type,
			&property.Title,
			&property.Description,
			&property.DefaultValue,
			&entries,
		)

		if err != nil {
			switch err {
			case pgx.ErrNoRows:
				continue
			default:
				rows.Close()

				repo.lgr.Error("getGoogleApplicationManagedPropertiesByAppID: row scan error")

				return nil, err
			}
		}

		err = property.ConvertBytesToEntries(entries)
		if err != nil {
			rows.Close()

			repo.lgr.Error("getGoogleApplicationManagedPropertiesByAppID [ConvertBytesToEntries]: ", err)

			return nil, err
		}

		result = append(result, &property)
	}

	rows.Close()

	for _, property := range result {
		property.NestedProperties, err = repo.getGoogleApplicationManagedPropertiesByPropertyID(ctx, property.ManagedPropertyID)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (repo *repository) getGoogleApplicationManagedPropertiesByPropertyID(ctx context.Context, managedPropertyID int) ([]*models.ManagedProperty, error) {
	var err error

	result := make([]*models.ManagedProperty, 0)

	rows, err := repo.client.Query(ctx, get.GoogleApplicationManagedPropertyByPropertyIDSQL, managedPropertyID)
	if err != nil {
		repo.lgr.Error("getGoogleApplicationManagedPropertiesByPropertyID: get rows error")

		return nil, err
	}

	for rows.Next() {
		var (
			property models.ManagedProperty
			entries  []byte
		)

		err := rows.Scan(
			&property.ManagedPropertyID,
			&property.Key,
			&property.Type,
			&property.Title,
			&property.Description,
			&property.DefaultValue,
			&entries,
		)

		if err != nil {
			switch err {
			case pgx.ErrNoRows:
				continue
			default:
				rows.Close()

				repo.lgr.Error("getGoogleApplicationManagedPropertiesByPropertyID: row scan error")

				return nil, err
			}
		}

		err = property.ConvertBytesToEntries(entries)
		if err != nil {
			rows.Close()

			repo.lgr.Error("getGoogleApplicationManagedPropertiesByPropertyID [ConvertBytesToEntries]: ", err)

			return nil, err
		}

		result = append(result, &property)
	}

	rows.Close()

	for _, property := range result {
		property.NestedProperties, err = repo.getGoogleApplicationManagedPropertiesByPropertyID(ctx, property.ManagedPropertyID)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
