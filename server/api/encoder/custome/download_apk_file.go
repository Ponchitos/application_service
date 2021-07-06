package custome

import (
	"context"
	"fmt"
	"github.com/Ponchitos/application_service/server/api/encoder"
	"github.com/Ponchitos/application_service/server/internal/endpoints/application"
	"net/http"
)

func DownloadApkFileResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if err, ok := response.(error); ok && err != nil {
		encoder.Error(ctx, err, w)

		return nil
	}

	responseData := response.(*application.DownloadApkFileResponse)

	w.Header().Set("Content-Type", "application/vnd.android.package-archive")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, responseData.FileName))

	_, err := w.Write(responseData.File)
	if err != nil {
		encoder.Error(ctx, err, w)

		return nil
	}

	return nil
}
