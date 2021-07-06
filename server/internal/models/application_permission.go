package models

import "github.com/Ponchitos/application_service/server/tools/types"

type ApplicationPermission struct {
	ID                  int `json:"-"`
	GoogleApplicationID int `json:"-"`

	Created types.NullTime `json:"created"`

	PermissionID string `json:"permissionId"`
	Name         string `json:"name"`
	Description  string `json:"description"`
}
