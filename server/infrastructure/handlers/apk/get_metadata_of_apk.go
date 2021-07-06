package apk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Ponchitos/application_service/server/internal/models"
	"os"
	"os/exec"
	"strconv"
)

type scriptResult struct {
	PackageName         string `json:"package_name"`
	ApplicationLabel    string `json:"application_label"`
	VersionName         string `json:"version_name"`
	FileSha1Base64      string `json:"file_sha1_base64"`
	FileSha256Base64    string `json:"file_sha256_base64"`
	IconBase64          string `json:"icon_base64"`
	ExternallyHostedURL string `json:"externally_hosted_url"`
	VersionCode         string `json:"version_code"`
	MinimumSDK          string `json:"minimum_sdk"`

	FileSize int `json:"file_size"`

	NativeCodes        []string `json:"native_codes"`
	CertificateBase64s []string `json:"certificate_base64"`
	UsesFeatures       []string `json:"uses_feature"`

	Permissions []struct {
		Name          string `json:"name"`
		MaxSdkVersion string `json:"max_sdk_version"`
	} `json:"uses_permission"`
}

func (hAPK *apkHandler) GetMetadataOfAPK(filePath, name string) (*models.ApplicationMetadata, error) {
	var (
		commandResult bytes.Buffer
		errorResult   bytes.Buffer
		err           error
		tempResult    scriptResult
		result        *models.ApplicationMetadata
	)

	_, err = os.Stat(filePath)
	if err != nil {
		return nil, errors.NewErrorf("Cannot find apk file: %v", "Не удалось найти apk файл: %v", err)
	}

	args := []string{
		hAPK.scriptPath,
		fmt.Sprintf("--apk=%s", filePath),
		fmt.Sprintf("--externallyHostedUrl=%s", "https://localhost:8000"),
	}

	cmd := exec.Command(hAPK.pythonPath, args...)
	cmd.Stdout = &commandResult
	cmd.Stderr = &errorResult

	errC := cmd.Run()
	if errC != nil {
		if errM := json.Unmarshal(errorResult.Bytes(), &err); errM != nil {
			return nil, errors.NewErrorf("Cannot execute command: %s", "Не удалось выполнить комманду: %s", errC)
		}

		return nil, errors.NewErrorf("Cannot execute command: %s\n(details: %s)", "Не удалось выполнить комманду: %s\n(детали: %s)", errC, err)
	}

	err = json.Unmarshal(commandResult.Bytes(), &tempResult)
	if err != nil {
		return nil, err
	}

	result, err = hAPK.convertToMetadataModel(&tempResult, name)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (hAPK *apkHandler) convertToMetadataModel(temp *scriptResult, nameFile string) (*models.ApplicationMetadata, error) {
	fileSize := strconv.Itoa(temp.FileSize)

	minSDK, err := strconv.Atoi(temp.MinimumSDK)
	if err != nil {
		return nil, err
	}

	versionCode, err := strconv.Atoi(temp.VersionCode)
	if err != nil {
		return nil, err
	}

	return &models.ApplicationMetadata{
		UUID:                nameFile,
		PackageName:         temp.PackageName,
		ApplicationLabel:    temp.ApplicationLabel,
		VersionName:         temp.VersionName,
		FileSize:            fileSize,
		FileSha1Base64:      temp.FileSha1Base64,
		FileSha256Base64:    temp.FileSha256Base64,
		IconBase64:          temp.IconBase64,
		ExternallyHostedURL: temp.ExternallyHostedURL,
		NativeCodes:         temp.NativeCodes,
		CertificateBase64s:  temp.CertificateBase64s,
		UsesFeatures:        temp.UsesFeatures,
		UsesPermissions:     hAPK.convertPermissions(temp),
		VersionCode:         versionCode,
		MinimumSDK:          minSDK,
	}, nil
}

func (hAPK *apkHandler) convertPermissions(temp *scriptResult) []*models.Permission {
	var result []*models.Permission

	for _, element := range temp.Permissions {
		result = append(result, &models.Permission{
			Name:          element.Name,
			MaxSdkVersion: element.MaxSdkVersion,
		})
	}

	return result
}
