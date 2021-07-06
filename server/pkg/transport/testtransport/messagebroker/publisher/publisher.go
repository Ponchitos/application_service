package publisher

import (
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Ponchitos/application_service/server/pkg"
	"github.com/Ponchitos/application_service/server/tools/logger"
)

type testPublisher struct {
	logger  logger.Logger
	decoder DecoderFunc
	closed  bool
	queue   chan interface{}
}

func NewTestPublisher(logger logger.Logger, decoder DecoderFunc, queue chan interface{}) (pkg.Publisher, error) {
	if logger == nil {
		return nil, errors.NewError("Missing logger", "Не указана система логирования")
	}

	return &testPublisher{logger: logger, decoder: decoder, queue: queue}, nil
}

func (publisher *testPublisher) Publish(topic string, messages ...interface{}) error {
	if publisher.closed {
		return errors.NewError("Test publisher already closed", "Test поставщик закрыт")
	}

	for _, msg := range messages {
		request, err := publisher.decoder(topic, msg)
		if err != nil {
			return errors.NewErrorf("Cannot decode message - %v, reason - %v", "Не удалось декодировать сообщений - %v, причина - %v", msg, err)
		}

		publisher.queue <- request

		publisher.logger.Infow("Message sent to Test broker",
			"message", msg,
			"topic", topic,
		)
	}

	return nil
}

func (publisher *testPublisher) Close() error {
	if publisher.closed {
		return nil
	}

	publisher.closed = true

	return nil
}
