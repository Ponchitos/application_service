package server

import "github.com/Ponchitos/application_service/server/pkg/transport"

func ServerErrorHandler(errorHandler transport.ErrorHandler) Option {
	return func(serv *Server) { serv.errorHandler = errorHandler }
}
