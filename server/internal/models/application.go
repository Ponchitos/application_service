package models

type Application struct {
	ID int `json:"-"`

	ApplicationUUID string `json:"applicationUId"`

	VersionUUID string `json:"versionUId"`

	EnterpriseID string `json:"enterpriseId"`

	Icon string `json:"icon"`

	UsesFeatures []string `json:"usesFeatures"`

	Permissions []string `json:"permissions"`

	ApplicationSettings []*ManagedProperty `json:"applicationSettings"`

	ShortInfo `json:"shortInfo"`

	Availability `json:"availability"`
}

type ShortInfo struct {
	Name        string `json:"name"`
	PackageName string `json:"packageName"`
	VersionCode int    `json:"versionCode"`
	VersionName string `json:"versionName"`
	MinSDK      int    `json:"minSdk"`
}

type Availability struct {
	Location  string `json:"location"`
	Available string `json:"available"`
}
