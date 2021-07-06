package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/internal/services/applications"
	"github.com/Ponchitos/application_service/server/pkg"
)

type GetApplicationsVersionsRequest struct {
	ApplicationUUID string
	EnterpriseID    string
}

func MakeGetApplicationVersionsEndpoint(srv applications.Service) pkg.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		requestBody := request.(*GetApplicationsVersionsRequest)

		return srv.GetApplicationVersions(ctx, requestBody.ApplicationUUID, requestBody.EnterpriseID)
	}
}
