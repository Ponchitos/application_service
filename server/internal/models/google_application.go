package models

import "github.com/Ponchitos/application_service/server/tools/types"

const (
	ApprovedGoogleStatus   = "APPROVED"
	UnapprovedGoogleStatus = "UNAPPROVED"
)

type GoogleApplication struct {
	UUID string `json:"uid"`

	Name  string `json:"name"`
	Title string `json:"title"`

	Status string `json:"status"`

	Permissions       []*ApplicationPermission `json:"permissions"`
	AppTracks         []*ApplicationTrack      `json:"appTracks"`
	ManagedProperties []*ManagedProperty       `json:"managedProperties"`

	ID int `json:"-"`

	Created types.NullTime `json:"created"`
}

func (google *GoogleApplication) GetPermissionsAsString() []string {
	var result []string

	for _, permission := range google.Permissions {
		if permission != nil {
			result = append(result, permission.PermissionID)
		}
	}

	if len(result) == 0 {
		return make([]string, 0)
	}

	return result
}
