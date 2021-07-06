package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/api/encoder"
	"github.com/Ponchitos/application_service/server/api/encoder/custome"
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

func DownloadApkFileRoute(srv applications.Service, lgr logger.Logger) http.Handler {
	downloadApkFileEndpoint := applicationEndpoints.MakeDownloadApkFileEndpoint(srv)
	downloadApkFileEndpoint = endpointMiddleware.ContextInit(downloadApkFileEndpoint, time.Second*10)

	return server.NewHTTPSeverHandler(
		downloadApkFileEndpoint,
		decoderDownloadApkFileRequest,
		custome.DownloadApkFileResponse,
		encoder.Error,
		errorHandler.NewErrorHandler(lgr),
	)
}

func decoderDownloadApkFileRequest(_ context.Context, request *http.Request) (interface{}, error) {
	applicationUUID := mux.Vars(request)["applicationUId"]
	applicationVersionUUID := mux.Vars(request)["versionUId"]
	enterpriseID := request.URL.Query().Get("enterpriseId")

	if len(applicationUUID) < 1 {
		return nil, errors.NewError("applicationUId is wrong", "Некорректное значение applicationUId")
	}

	if len(applicationVersionUUID) < 1 {
		return nil, errors.NewError("versionUId is wrong", "Некорректное значение versionUId")
	}

	if len(enterpriseID) < 1 {
		return nil, errors.NewError("enterpriseId is wrong", "Некорректное значение enterpriseId")
	}

	return &applicationEndpoints.DownloadApkFileRequest{
		VersionUUID:     applicationVersionUUID,
		ApplicationUUID: applicationUUID,
		EnterpriseID:    enterpriseID,
	}, nil
}
