package subscriber

type DecodeFunc func([]byte) (interface{}, error)
