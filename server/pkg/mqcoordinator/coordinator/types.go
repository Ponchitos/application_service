package coordinator

import "context"

type EncoderFunc func(context.Context, interface{}) ([]interface{}, error)
type DecoderFunc func(context.Context, interface{}) (interface{}, error)
