package application

import (
	"context"
	"encoding/json"
	"github.com/Ponchitos/application_service/server/api/encoder"
	errorHandler "github.com/Ponchitos/application_service/server/api/errorhandlers/logger"
	endpointMiddleware "github.com/Ponchitos/application_service/server/infrastructure/middlerware/endpoints"
	applicationEndpoints "github.com/Ponchitos/application_service/server/internal/endpoints/application"
	"github.com/Ponchitos/application_service/server/internal/services/applications"
	"github.com/Ponchitos/application_service/server/pkg/transport/http/server"
	"github.com/Ponchitos/application_service/server/tools/logger"
	"io/ioutil"
	"net/http"
	"time"
)

func UninstallApplicationCompleteRoute(srv applications.Service, lgr logger.Logger) http.Handler {
	uninstallApplicationCompleteEndpoint := applicationEndpoints.MakeUninstallApplicationCompleteEndpoint(srv)
	uninstallApplicationCompleteEndpoint = endpointMiddleware.ContextInit(uninstallApplicationCompleteEndpoint, time.Second*10)

	return server.NewHTTPSeverHandler(
		uninstallApplicationCompleteEndpoint,
		decodeUninstallApplicationCompleteRequest(lgr),
		encoder.Response,
		encoder.Error,
		errorHandler.NewErrorHandler(lgr),
	)
}

func decodeUninstallApplicationCompleteRequest(lgr logger.Logger) server.DecodeRequestFunc {
	return func(ctx context.Context, request *http.Request) (interface{}, error) {
		var requestBody applicationEndpoints.UninstallApplicationCompleteRequest

		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}

		if err := request.Body.Close(); err != nil {
			lgr.Errorf("decodeUninstallApplicationCompleteRequest: %v", err)
		}

		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			return nil, err
		}

		lgr.Info("decodeUninstallApplicationCompleteRequest request body: ", string(body))

		return &requestBody, nil
	}
}
