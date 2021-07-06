package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/internal/services/applications"
	"github.com/Ponchitos/application_service/server/pkg"
)

type UninstallApplicationCompleteRequest struct {
	VersionUUID     string `json:"versionUId"`
	EnterpriseID    string `json:"enterpriseId"`
	ApplicationUUID string `json:"applicationUId"`
}

func MakeUninstallApplicationCompleteEndpoint(srv applications.Service) pkg.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		requestBody := request.(*UninstallApplicationCompleteRequest)

		err := srv.UninstallApplicationComplete(ctx, requestBody.ApplicationUUID, requestBody.VersionUUID, requestBody.EnterpriseID)

		return nil, err
	}
}
