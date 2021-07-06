package models

import "github.com/Ponchitos/application_service/server/tools/types"

type ApplicationGroupSet struct {
	SetID int `json:"-"`

	ApplicationID int `json:"-"`
	VersionID     int `json:"-"`

	GroupUUID string `json:"groupUId"`

	GroupName      string `json:"groupName"`
	EnterpriseID   string `json:"enterpriseId"`
	Status         string `json:"status"`
	PreviousStatus string `json:"previousStatus,omitempty"`

	Created  types.NullTime `json:"created"`
	Modified types.NullTime `json:"modified"`
	Deleted  types.NullTime `json:"deleted"`
}
