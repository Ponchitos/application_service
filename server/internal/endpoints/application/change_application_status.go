package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/internal/services/applications"
	"github.com/Ponchitos/application_service/server/pkg"
)

type ChangeApplicationStatusRequest struct {
	VersionUUID  string `json:"versionUId"`
	EnterpriseID string `json:"enterpriseId"`
	Status       string `json:"status"`
}

func MakeChangeApplicationStatusEndpoint(srv applications.Service) pkg.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		requestBody := request.(*ChangeApplicationStatusRequest)

		err := srv.ChangeApplicationStatus(ctx, requestBody.VersionUUID, requestBody.EnterpriseID, requestBody.Status)

		return nil, err
	}
}
