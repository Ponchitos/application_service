package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/internal/services/applications"
	"github.com/Ponchitos/application_service/server/pkg"
)

type CancelApkUploadRequest struct {
	UUID string `json:"uid"`
}

func MakeCancelApkUploadEndpoint(srv applications.Service) pkg.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		requestBody := request.(*CancelApkUploadRequest)

		return "OK", srv.CancelApkUpload(ctx, requestBody.UUID)
	}
}
