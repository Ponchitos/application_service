package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/api/encoder"
	errorHandler "github.com/Ponchitos/application_service/server/api/errorhandlers/logger"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	endpointMiddleware "github.com/Ponchitos/application_service/server/infrastructure/middlerware/endpoints"
	applicationEndpoints "github.com/Ponchitos/application_service/server/internal/endpoints/application"
	"github.com/Ponchitos/application_service/server/internal/services/applications"
	"github.com/Ponchitos/application_service/server/pkg/transport/http/server"
	"github.com/Ponchitos/application_service/server/tools/logger"
	"net/http"
	"strconv"
	"time"
)

func GetApplicationsRoute(srv applications.Service, lgr logger.Logger) http.Handler {
	getApplicationsEndpoint := applicationEndpoints.MakeGetApplicationsEndpoint(srv)
	getApplicationsEndpoint = endpointMiddleware.ContextInit(getApplicationsEndpoint, time.Second*10)

	return server.NewHTTPSeverHandler(
		getApplicationsEndpoint,
		decoderGetApplicationsRequest,
		encoder.Response,
		encoder.Error,
		errorHandler.NewErrorHandler(lgr),
	)
}

func decoderGetApplicationsRequest(_ context.Context, request *http.Request) (interface{}, error) {
	enterpriseID := request.URL.Query().Get("enterpriseId")
	voffset := request.URL.Query().Get("offset")
	vlimit := request.URL.Query().Get("limit")

	if len(enterpriseID) < 1 {
		return nil, errors.NewError("enterpriseId is wrong", "Некорректное значение enterpriseId")
	}

	offset, err := strconv.Atoi(voffset)
	if err != nil {
		return nil, errors.NewError("offset is wrong", "некорректное значение offset")
	}

	limit, err := strconv.Atoi(vlimit)
	if err != nil {
		return nil, errors.NewError("limit is wrong", "некорректное значение limit")
	}

	return &applicationEndpoints.GetApplicationsRequest{
		EnterpriseID: enterpriseID,
		Offset:       offset,
		Limit:        limit,
	}, nil
}
