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

func CancelApkUploadRoute(srv applications.Service, lgr logger.Logger) http.Handler {
	cancelApkUploadEndpoint := applicationEndpoints.MakeCancelApkUploadEndpoint(srv)
	cancelApkUploadEndpoint = endpointMiddleware.ContextInit(cancelApkUploadEndpoint, time.Second*10)

	return server.NewHTTPSeverHandler(
		cancelApkUploadEndpoint,
		decoderCancelApkUpload(lgr),
		encoder.Response,
		encoder.Error,
		errorHandler.NewErrorHandler(lgr),
	)
}

func decoderCancelApkUpload(lgr logger.Logger) server.DecodeRequestFunc {
	return func(ctx context.Context, request *http.Request) (interface{}, error) {
		var requestBody applicationEndpoints.CancelApkUploadRequest

		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}

		if err := request.Body.Close(); err != nil {
			lgr.Errorf("decoderCancelApkUpload: %v", err)
		}

		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			return nil, err
		}

		lgr.Info("decoderCancelApkUpload request body: ", string(body))

		return &requestBody, nil
	}
}
