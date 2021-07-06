package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/internal/services/applications"
	"github.com/Ponchitos/application_service/server/pkg"
)

type GetApplicationInfoRequest struct {
	ApplicationUUID        string
	ApplicationVersionUUID string
	EnterpriseID           string
}

func MakeGetApplicationInfoEndpoint(srv applications.Service) pkg.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		requestBody := request.(*GetApplicationInfoRequest)

		return srv.GetApplicationInfo(ctx, requestBody.ApplicationUUID, requestBody.ApplicationVersionUUID, requestBody.EnterpriseID)
	}
}
