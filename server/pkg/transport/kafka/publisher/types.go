package publisher

import "github.com/Shopify/sarama"

type DecoderFunc func(topic string, message interface{}) (*sarama.ProducerMessage, error)
