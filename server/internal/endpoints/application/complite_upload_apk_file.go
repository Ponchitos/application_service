package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/internal/services/applications"
	"github.com/Ponchitos/application_service/server/pkg"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type CompleteUploadAPKFileRequest struct {
	UUID         string `json:"uid"`
	EnterpriseID string `json:"enterpriseId"`
}

func (request *CompleteUploadAPKFileRequest) Validate() error {
	return validation.ValidateStruct(
		request,
		validation.Field(&request.UUID, validation.Required, is.UUID),
		validation.Field(&request.EnterpriseID, validation.Required, validation.Length(1, 100)),
	)
}

type CompleteUploadAPKFileResponse struct {
	ApplicationUUID string `json:"applicationUId"`
	VersionUUID     string `json:"versionUId"`
}

func MakeCompleteUploadAPKFileEndpoint(srv applications.Service) pkg.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		requestBody := request.(*CompleteUploadAPKFileRequest)

		applicationUUID, applicationVersionUUID, err := srv.CompleteUploadFile(ctx, requestBody.UUID, requestBody.EnterpriseID)
		if err != nil {
			return nil, err
		}

		return &CompleteUploadAPKFileResponse{
			ApplicationUUID: applicationUUID,
			VersionUUID:     applicationVersionUUID,
		}, nil
	}
}
