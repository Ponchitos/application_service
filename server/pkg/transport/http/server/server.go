package server

import (
	"context"
	"github.com/Ponchitos/application_service/server/pkg"
	"github.com/Ponchitos/application_service/server/pkg/transport"
	"net/http"
)

type Server struct {
	endpoint pkg.Endpoint

	decoder DecodeRequestFunc
	encoder EncodeResponseFunc

	errorEncoder ErrorEncoder

	errorHandler transport.ErrorHandler

	before []RequestFunc
	after  []ResponseFunc
}

func NewHTTPSeverHandler(
	endpoint pkg.Endpoint,
	decoder DecodeRequestFunc,
	encoder EncodeResponseFunc,
	errorEncoder ErrorEncoder,
	errorHandler transport.ErrorHandler,
	options ...Option) *Server {
	server := &Server{
		endpoint:     endpoint,
		decoder:      decoder,
		encoder:      encoder,
		errorEncoder: errorEncoder,
		errorHandler: errorHandler,
	}

	for _, option := range options {
		option(server)
	}

	return server
}

func (serv *Server) serveHTTPErrorHandlerFunction(ctx context.Context, err error, w http.ResponseWriter) {
	serv.errorHandler.Handle(ctx, err)
	serv.errorEncoder(ctx, err, w)
}

func (serv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	for _, handler := range serv.before {
		ctx = handler(ctx, r)
	}

	requestObject, err := serv.decoder(ctx, r)
	if err != nil {
		serv.serveHTTPErrorHandlerFunction(ctx, err, w)

		return
	}

	response, err := serv.endpoint(ctx, requestObject)
	if err != nil {
		serv.serveHTTPErrorHandlerFunction(ctx, err, w)

		return
	}

	for _, handler := range serv.after {
		ctx = handler(ctx, w)
	}

	err = serv.encoder(ctx, w, response)
	if err != nil {
		serv.serveHTTPErrorHandlerFunction(ctx, err, w)

		return
	}
}
