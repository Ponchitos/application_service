package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/internal/services/applications"
	"github.com/Ponchitos/application_service/server/pkg"
	"mime/multipart"
)

type UploadFileRequest struct {
	APK multipart.File
}

func MakeUploadAPKFileEndpoint(srv applications.Service) pkg.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		requestBody := request.(*UploadFileRequest)

		return srv.UploadApkFile(ctx, requestBody.APK)
	}
}
