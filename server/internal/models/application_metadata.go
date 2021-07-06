package models

import (
	"encoding/json"
	"github.com/Ponchitos/application_service/server/tools/types"
)

type ApplicationMetadata struct {
	UUID string `json:"uid"`

	ID int `json:"-"`

	Link string `json:"link"`

	PackageName         string `json:"packageName"`
	ApplicationLabel    string `json:"applicationLabel"`
	VersionName         string `json:"versionName"`
	FileSize            string `json:"fileSize"`
	FileSha1Base64      string `json:"fileSha1Base64"`
	FileSha256Base64    string `json:"fileSha256Base64"`
	IconBase64          string `json:"iconBase64"`
	ExternallyHostedURL string `json:"externallyHostedUrl"`

	NativeCodes        []string `json:"nativeCodes"`
	CertificateBase64s []string `json:"certificateBase64s"`
	UsesFeatures       []string `json:"usesFeatures"`

	UsesPermissions []*Permission `json:"usesPermissions"`

	VersionCode int `json:"versionCode"`
	MinimumSDK  int `json:"minimumSdk"`

	Created types.NullTime `json:"created"`
}

func (metadata *ApplicationMetadata) GetPermissionsAsStrings() []string {
	var result []string

	for _, permission := range metadata.UsesPermissions {
		if permission != nil {
			result = append(result, permission.Name)
		}

	}

	return result
}

func (metadata *ApplicationMetadata) ConvertPermissionsToBytes() ([]byte, error) {
	return json.Marshal(metadata.UsesPermissions)
}

func (metadata *ApplicationMetadata) ConvertBytesToPermissions(data []byte) error {
	return json.Unmarshal(data, &metadata.UsesPermissions)
}
