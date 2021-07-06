package subscriber

import "github.com/Shopify/sarama"

type DecodeFunc func(*sarama.ConsumerMessage) (interface{}, error)
