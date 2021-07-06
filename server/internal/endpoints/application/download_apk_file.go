package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/internal/services/applications"
	"github.com/Ponchitos/application_service/server/pkg"
)

type DownloadApkFileRequest struct {
	VersionUUID     string
	ApplicationUUID string
	EnterpriseID    string
}

type DownloadApkFileResponse struct {
	File     []byte
	FileName string
}

func MakeDownloadApkFileEndpoint(srv applications.Service) pkg.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		requestBody := request.(*DownloadApkFileRequest)

		file, fileName, err := srv.DownloadApkFile(ctx, requestBody.ApplicationUUID, requestBody.VersionUUID, requestBody.EnterpriseID)
		if err != nil {
			return nil, err
		}

		return &DownloadApkFileResponse{
			File:     file,
			FileName: fileName,
		}, nil
	}
}
