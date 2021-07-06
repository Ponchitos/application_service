package logger

import (
	"context"
	"github.com/Ponchitos/application_service/server/pkg/transport"
	"github.com/Ponchitos/application_service/server/tools/logger"
)

type errorHandler struct {
	logger logger.Logger
}

func NewErrorHandler(logger logger.Logger) transport.ErrorHandler {
	return &errorHandler{logger}
}

func (handler *errorHandler) Handle(_ context.Context, err error) {
	handler.logger.Error("Error: ", err)
}
