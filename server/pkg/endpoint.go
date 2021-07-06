package pkg

import "context"

type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)

type WithoutResponseEndpoint func(ctx context.Context, request interface{}) error

type Middleware func(Endpoint) Endpoint

type EncoderFunc func(context.Context, interface{}) (result interface{}, err error)
type DecoderFunc func(context.Context, interface{}) (result interface{}, err error)
