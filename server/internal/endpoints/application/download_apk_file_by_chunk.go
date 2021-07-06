package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/internal/services/applications"
	"github.com/Ponchitos/application_service/server/pkg"
)

type DownloadApkFileByChunkRequest struct {
	Start           int64
	End             int64
	ApplicationUUID string
	VersionUUID     string
	EnterpriseID    string
}

type DownloadApkFileByChunkResponse struct {
	Chunk    []byte
	Size     int64
	FileName string
	Start    int64
	End      int64
}

func MakeDownloadApkFileByChunkEndpoint(srv applications.Service) pkg.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		requestBody := request.(*DownloadApkFileByChunkRequest)

		chunk, size, fileName, err := srv.DownloadApkFileByChunk(ctx, requestBody.ApplicationUUID, requestBody.VersionUUID, requestBody.EnterpriseID, requestBody.Start, requestBody.End)

		return &DownloadApkFileByChunkResponse{
			Chunk:    chunk,
			Size:     size,
			FileName: fileName,
			Start:    requestBody.Start,
			End:      requestBody.End,
		}, err
	}
}
