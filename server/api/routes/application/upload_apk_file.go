package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/api/encoder"
	errorHandler "github.com/Ponchitos/application_service/server/api/errorhandlers/logger"
	errorConstants "github.com/Ponchitos/application_service/server/errors"
	endpointMiddleware "github.com/Ponchitos/application_service/server/infrastructure/middlerware/endpoints"
	applicationEndpoints "github.com/Ponchitos/application_service/server/internal/endpoints/application"
	"github.com/Ponchitos/application_service/server/internal/services/applications"
	"github.com/Ponchitos/application_service/server/pkg/transport/http/server"
	"github.com/Ponchitos/application_service/server/tools/logger"
	"net/http"
	"time"
)

func UploadAPKFileRoute(srv applications.Service, lgr logger.Logger) http.Handler {
	uploadAPKEndpoint := applicationEndpoints.MakeUploadAPKFileEndpoint(srv)
	uploadAPKEndpoint = endpointMiddleware.ContextInit(uploadAPKEndpoint, time.Second*10)

	return server.NewHTTPSeverHandler(
		uploadAPKEndpoint,
		decoderUploadAPKFileRequest(lgr),
		encoder.Response,
		encoder.Error,
		errorHandler.NewErrorHandler(lgr),
	)
}

func decoderUploadAPKFileRequest(lgr logger.Logger) server.DecodeRequestFunc {
	return func(ctx context.Context, request *http.Request) (interface{}, error) {
		const (
			fileKey     = "file"
			contentType = "application/vnd.android.package-archive"
		)

		err := request.ParseMultipartForm(10 << 20)
		if err != nil {
			lgr.Error("decoderUploadAPKFileRequest [ParseMultipartForm]: ", err)

			return nil, err
		}

		file, fileHeader, err := request.FormFile(fileKey)
		if err != nil {
			lgr.Error("decoderUploadAPKFileRequest [FormFile]: ", err)

			return nil, err
		}

		if fileHeader.Header.Get("Content-Type") != contentType {
			return nil, errorConstants.BadRequest
		}

		return &applicationEndpoints.UploadFileRequest{
			APK: file,
		}, nil
	}
}
