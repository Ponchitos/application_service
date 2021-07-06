package models

import "github.com/Ponchitos/application_service/server/tools/types"

const (
	Private = "PRIVATE"
	Public  = "PUBLIC"

	Internal = "SELF_HOSTED"
	Google   = "GOOGLE_PLAY"
)

const (
	Approved         = "APPROVED"
	Installed        = "INSTALLED"
	Uninstalled      = "UNINSTALLED"
	UpdateAvailable  = "UPDATE_AVAILABLE"
	Updated          = "UPDATED"
	WaitingInstall   = "WAITING_INSTALL"
	WaitingUpdate    = "WAITING_UPDATE"
	WaitingUninstall = "WAITING_UNINSTALL"
)

type BasicApplication struct {
	UUID string `json:"uid"`
	ID   int    `json:"-"`

	Created types.NullTime `json:"created"`

	VersionUUID string `json:"versionUid"`

	Icon        string `json:"icon"`
	Name        string `json:"name"`
	PackageName string `json:"packageName"`
	Available   string `json:"available"`
	Location    string `json:"location"`
	VersionName string `json:"versionName"`
	Status      string `json:"status"`

	VersionCode int `json:"versionCode"`
	MinSDK      int `json:"minSdk"`
}
