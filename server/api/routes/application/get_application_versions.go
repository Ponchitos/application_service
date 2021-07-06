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
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func GetApplicationVersions(srv applications.Service, lgr logger.Logger) http.Handler {
	getApplicationVersionsEndpoint := applicationEndpoints.MakeGetApplicationVersionsEndpoint(srv)
	getApplicationVersionsEndpoint = endpointMiddleware.ContextInit(getApplicationVersionsEndpoint, time.Second*10)

	return server.NewHTTPSeverHandler(
		getApplicationVersionsEndpoint,
		decoderGetApplicationVersionsRequest,
		encoder.Response,
		encoder.Error,
		errorHandler.NewErrorHandler(lgr),
	)
}

func decoderGetApplicationVersionsRequest(_ context.Context, request *http.Request) (interface{}, error) {
	applicationUUID := mux.Vars(request)["applicationUId"]
	enterpriseID := request.URL.Query().Get("enterpriseId")

	if len(applicationUUID) < 1 {
		return nil, errors.NewError("applicationUId is wrong", "Некорректное значение applicationUId")
	}

	if len(enterpriseID) < 1 {
		return nil, errors.NewError("enterpriseId is wrong", "Некорректное значение enterpriseId")
	}

	return &applicationEndpoints.GetApplicationsVersionsRequest{
		ApplicationUUID: applicationUUID,
		EnterpriseID:    enterpriseID,
	}, nil
}
