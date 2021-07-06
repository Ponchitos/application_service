package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/internal/models"
	"github.com/Ponchitos/application_service/server/internal/services/applications"
	"github.com/Ponchitos/application_service/server/pkg"
)

type GetApplicationsRequest struct {
	EnterpriseID string `json:"enterpriseId"`
	Offset       int    `json:"offset"`
	Limit        int    `json:"limit"`
}

type GetApplicationsResponse struct {
	Count        int                        `json:"count"`
	Applications []*models.BasicApplication `json:"applications"`
}

func MakeGetApplicationsEndpoint(srv applications.Service) pkg.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var (
			response GetApplicationsResponse
			err      error
		)

		requestBody := request.(*GetApplicationsRequest)

		response.Count, response.Applications, err = srv.GetApplications(ctx, requestBody.EnterpriseID, requestBody.Offset, requestBody.Limit)

		return &response, err
	}
}
