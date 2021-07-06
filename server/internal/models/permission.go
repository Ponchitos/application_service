package models

type Permission struct {
	Name          string `json:"name"`
	MaxSdkVersion string `json:"maxSdkVersion"`
}
