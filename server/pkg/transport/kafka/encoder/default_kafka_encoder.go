package encoder

import "github.com/Shopify/sarama"

func DefaultKafkaEncoderSubscriberFunc(kafkaMessage *sarama.ConsumerMessage) (interface{}, error) {
	return kafkaMessage.Value, nil
}
