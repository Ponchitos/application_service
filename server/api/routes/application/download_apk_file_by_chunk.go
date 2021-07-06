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

	"strconv"
	"strings"
)

func DownloadApkFileByChunkRoute(srv applications.Service, lgr logger.Logger) http.Handler {
	downloadApkFileByChunkEndpoint := applicationEndpoints.MakeDownloadApkFileByChunkEndpoint(srv)
	downloadApkFileByChunkEndpoint = endpointMiddleware.ContextInit(downloadApkFileByChunkEndpoint, time.Second*10)

	return server.NewHTTPSeverHandler(
		downloadApkFileByChunkEndpoint,
		decoderDownloadApkFileByChunkRequest(lgr),
		custome.DownloadApkFileByChunkRoute,
		encoder.Error,
		errorHandler.NewErrorHandler(lgr),
	)
}

func decoderDownloadApkFileByChunkRequest(lgr logger.Logger) server.DecodeRequestFunc {
	return func(_ context.Context, request *http.Request) (interface{}, error) {
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

		requestRange := request.Header.Get("range")
		if len(requestRange) == 0 {
			return nil, errors.NewError("Range header not found", "Range заголовок не найден")
		}

		requestRange = requestRange[6:]

		splitRange := strings.Split(requestRange, "-")

		if len(splitRange) != 2 {
			return nil, errors.NewError("Invalid values for header 'Range'", "Не корректные значения для заголовка 'Range'")
		}

		begin, err := strconv.ParseInt(splitRange[0], 10, 64)
		if err != nil {
			return nil, errors.NewErrorf("Can not get begin value: %v", "Не удалось получить начальное значение: %v", err)
		}

		end, err := strconv.ParseInt(splitRange[1], 10, 64)
		if err != nil {
			return nil, errors.NewErrorf("Can not get end value: %v", "Не удалось получить конечное значение: %v", err)
		}

		if begin > end {
			return nil, errors.NewError("Range header values are invalid", "Значения Range заголовка не корректны")
		}

		return &applicationEndpoints.DownloadApkFileByChunkRequest{
			Start:           begin,
			End:             end,
			ApplicationUUID: applicationUUID,
			VersionUUID:     applicationVersionUUID,
			EnterpriseID:    enterpriseID,
		}, nil
	}
}
