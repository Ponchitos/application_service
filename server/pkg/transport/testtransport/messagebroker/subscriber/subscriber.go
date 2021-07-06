package subscriber

import (
	"context"
	customError "github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Ponchitos/application_service/server/pkg"
	"github.com/Ponchitos/application_service/server/tools/logger"
)

type testSubscriber struct {
	logger  logger.Logger
	decoder DecodeFunc
	queue   chan interface{}
	closed  bool
}

func NewTestSubscriber(logger logger.Logger, decoder DecodeFunc, queue chan interface{}) (pkg.Subscriber, error) {
	if logger == nil {
		return nil, customError.NewError("Missing logger", "Не указана система логирования")
	}

	return &testSubscriber{
		logger:  logger,
		decoder: decoder,
		queue:   queue,
	}, nil
}

func (subscriber *testSubscriber) Subscribe(ctx context.Context, topic string) (<-chan interface{}, error) {
	if subscriber.closed {
		return nil, customError.NewError("Test subscriber already closed", "Test подписчик закрыт")
	}

	subscriber.logger.Infow("Subscribing to Test topic",
		"topic", topic,
	)

	outputCh := make(chan interface{}, 0)

	for msg := range subscriber.queue {
		outputCh <- msg
	}

	return outputCh, nil
}

func (subscriber *testSubscriber) Close() error {
	if subscriber.closed {
		return nil
	}

	subscriber.closed = true

	subscriber.logger.Info("Test subscriber closed")

	return nil
}
