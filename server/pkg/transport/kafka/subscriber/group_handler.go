package subscriber

import (
	"context"
	"github.com/Ponchitos/application_service/server/tools/logger"
	"github.com/Shopify/sarama"
)

type consumerGroupHandler struct {
	ctx context.Context

	messageHandler *messageHandler

	logger logger.Logger

	closing chan struct{}
}

func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (handler consumerGroupHandler) ConsumeClaim(consumerGroupSession sarama.ConsumerGroupSession, consumerGroupClaim sarama.ConsumerGroupClaim) error {
	kafkaMessages := consumerGroupClaim.Messages()
	handler.logger.Debugw("Consume claimed",
		"kafka_partition", consumerGroupClaim.Partition(),
		"kafka_initial_offset", consumerGroupClaim.InitialOffset(),
	)

	for {
		select {
		case kafkaMsg, ok := <-kafkaMessages:
			if !ok {
				handler.logger.Debugw("kafkaMessages is closed, stopping consumerGroupHandler",
					"kafka_partition", consumerGroupClaim.Partition(),
					"kafka_initial_offset", consumerGroupClaim.InitialOffset(),
				)

				return nil
			}
			if err := handler.messageHandler.processMessage(handler.ctx, kafkaMsg); err != nil {
				return err
			}

		case <-handler.closing:
			handler.logger.Debugw("Subscriber is closing, stopping consumerGroupHandler",
				"kafka_partition", consumerGroupClaim.Partition(),
				"kafka_initial_offset", consumerGroupClaim.InitialOffset())

			return nil

		case <-handler.ctx.Done():
			handler.logger.Debugw("Ctx was cancelled, stopping consumerGroupHandler",
				"kafka_partition", consumerGroupClaim.Partition(),
				"kafka_initial_offset", consumerGroupClaim.InitialOffset())

			return nil
		}
	}
}
