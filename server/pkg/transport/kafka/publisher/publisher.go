package publisher

import (
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Ponchitos/application_service/server/pkg"
	"github.com/Ponchitos/application_service/server/tools/logger"
	"github.com/Shopify/sarama"
)

type kafkaPublisher struct {
	producer sarama.SyncProducer

	logger logger.Logger

	decoder DecoderFunc

	closed bool
}

func NewPublisher(config *ConfigPublisher, logger logger.Logger, decoder DecoderFunc) (pkg.Publisher, error) {
	config.setDefaults()

	err := config.validate()
	if err != nil {
		return nil, err
	}

	if logger == nil {
		return nil, errors.NewError("Missing logger", "Не указана система логирования")
	}

	producer, err := sarama.NewSyncProducer(config.Brokers, config.SaramaConfig)
	if err != nil {
		return nil, errors.NewErrorf("Cannot create kafka producer: %v", "Не удалось создать kafka поставщика: %v", err)
	}

	return &kafkaPublisher{producer: producer, logger: logger, decoder: decoder}, nil
}

func (publisher *kafkaPublisher) Publish(topic string, messages ...interface{}) error {
	if publisher.closed {
		return errors.NewError("Kafka publisher already closed", "Kafka поставщик закрыт")
	}

	for _, msg := range messages {
		request, err := publisher.decoder(topic, msg)
		if err != nil {
			return errors.NewErrorf("Cannot decode message - %v, reason - %v", "Не удалось декодировать сообщений - %v, причина - %v", msg, err)
		}

		partition, offset, err := publisher.producer.SendMessage(request)
		if err != nil {
			return errors.NewErrorf("Cannot produce message - %v, reason - %v", "Не удалось обработать сообщение для отправки - %v, причина - %v", request, err)
		}

		publisher.logger.Infow("Message sent to Kafka",
			"message", msg,
			"kafka_partition", partition,
			"kafka_offset", offset,
		)
	}

	return nil
}

func (publisher *kafkaPublisher) Close() error {
	if publisher.closed {
		return nil
	}

	publisher.closed = true

	err := publisher.producer.Close()
	if err != nil {
		return errors.NewErrorf("Cannot close kafka producer: %v", "Не удалось закрыть поставщика kafka: %v", err)
	}

	return nil
}
