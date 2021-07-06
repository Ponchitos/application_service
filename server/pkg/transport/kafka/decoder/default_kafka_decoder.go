package decoder

import (
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Shopify/sarama"
)

func DefaultKafkaDecoderPublisherFunc(topic string, message interface{}) (*sarama.ProducerMessage, error) {
	messageOfBytes, ok := message.([]byte)
	if !ok {
		return nil, errors.NewError("Don't valid message type", "Не корректный тип сообщения")
	}

	return &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(messageOfBytes),
	}, nil
}
