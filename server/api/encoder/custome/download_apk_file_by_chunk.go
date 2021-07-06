package custome

import (
	"context"
	"fmt"
	"github.com/Ponchitos/application_service/server/api/encoder"
	"github.com/Ponchitos/application_service/server/internal/endpoints/application"
	"net/http"
)

func DownloadApkFileByChunkRoute(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if err, ok := response.(error); ok && err != nil {
		encoder.Error(ctx, err, w)

		return nil
	}

	responseData := response.(*application.DownloadApkFileByChunkResponse)

	w.WriteHeader(http.StatusPartialContent)

	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %v-%v/%v", responseData.Start, responseData.End, responseData.Size))
	w.Header().Set("Content-Length", fmt.Sprintf("%v", responseData.End-responseData.Start))

	_, err := w.Write(responseData.Chunk)
	if err != nil {
		encoder.Error(ctx, err, w)

		return nil
	}

	return nil
}
