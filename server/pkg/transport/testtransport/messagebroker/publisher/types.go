package publisher

type DecoderFunc func(topic string, message interface{}) (interface{}, error)
