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

func CompleteUploadAPKFileRoute(srv applications.Service, lgr logger.Logger) http.Handler {
	completeUploadAPKFileEndpoint := applicationEndpoints.MakeCompleteUploadAPKFileEndpoint(srv)
	completeUploadAPKFileEndpoint = endpointMiddleware.ContextInit(completeUploadAPKFileEndpoint, time.Second*10)

	return server.NewHTTPSeverHandler(
		completeUploadAPKFileEndpoint,
		decoderCompleteUploadAPKFileRequest(lgr),
		encoder.Response,
		encoder.Error,
		errorHandler.NewErrorHandler(lgr),
	)
}

func decoderCompleteUploadAPKFileRequest(lgr logger.Logger) server.DecodeRequestFunc {
	return func(ctx context.Context, request *http.Request) (interface{}, error) {
		var requestBody applicationEndpoints.CompleteUploadAPKFileRequest

		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}

		if err := request.Body.Close(); err != nil {
			lgr.Errorf("decoderCompleteUploadAPKFileRequest: %v", err)
		}

		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			return nil, err
		}

		err = requestBody.Validate()
		if err != nil {
			return nil, err
		}

		lgr.Info("decoderCompleteUploadAPKFileRequest request body: ", string(body))

		return &requestBody, nil
	}
}
