package subscriber

import (
	"context"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Ponchitos/application_service/server/tools/logger"
	"github.com/Shopify/sarama"
	"time"
)

type messageHandler struct {
	outputCh chan<- interface{}

	decoder DecodeFunc

	nackReset time.Duration

	logger logger.Logger

	closing chan struct{}
}

func (handler *messageHandler) processMessage(
	ctx context.Context,
	msg *sarama.ConsumerMessage,
) error {

	handler.logger.Infow("Received message from Kafka",
		"kafka_partition_offset", msg.Offset,
		"kafka_partition", msg.Partition,
	)

	encodeMsg, err := handler.decoder(msg)
	if err != nil {
		return errors.NewErrorf("Message decode failed: %v", "Ошибка декодирования сообщения: %v", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for {
		select {
		case handler.outputCh <- encodeMsg:
			handler.logger.Info("Message sent to consumer")

			return nil

		case <-handler.closing:
			handler.logger.Info("Closing, message discarded")

			return nil

		case <-ctx.Done():
			handler.logger.Info("Closing, ctx cancelled before sent to consumer")

			return nil
		}
	}
}
