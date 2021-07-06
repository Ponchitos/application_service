package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/internal/services/applications"
	"github.com/Ponchitos/application_service/server/pkg"
)

type DeleteApkVersionRequest struct {
	ApplicationUUID string `json:"applicationUId"`
	VersionUUID     string `json:"versionUId"`
	EnterpriseID    string `json:"enterpriseId"`
}

func MakeDeleteApkVersionFileEndpoint(srv applications.Service) pkg.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		requestBody := request.(*DeleteApkVersionRequest)

		err := srv.DeleteVersionApplication(ctx, requestBody.ApplicationUUID, requestBody.VersionUUID, requestBody.EnterpriseID)

		return "OK", err
	}
}
