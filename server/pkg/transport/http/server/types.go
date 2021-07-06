package server

import (
	"context"
	"net/http"
)

type DecodeRequestFunc func(context.Context, *http.Request) (request interface{}, err error)

type EncodeResponseFunc func(context.Context, http.ResponseWriter, interface{}) error

type ErrorEncoder func(ctx context.Context, err error, w http.ResponseWriter)

type RequestFunc func(context.Context, *http.Request) context.Context

type ResponseFunc func(context.Context, http.ResponseWriter) context.Context

type Option func(*Server)
