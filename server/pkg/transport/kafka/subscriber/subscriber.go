package subscriber

import (
	"context"
	"errors"
	customError "github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Ponchitos/application_service/server/pkg"
	"github.com/Ponchitos/application_service/server/tools/logger"
	"github.com/Shopify/sarama"
	"sync"
	"time"
)

type kafkaSubscriber struct {
	logger logger.Logger

	config *ConfigSubscriber

	decoder DecodeFunc

	closing chan struct{}

	subscribersWg sync.WaitGroup

	closed bool
}

func NewSubscriber(
	config *ConfigSubscriber,
	logger logger.Logger,
	decoder DecodeFunc,
) (pkg.Subscriber, error) {

	config.setDefaults()

	if err := config.validate(); err != nil {
		return nil, err
	}

	if logger == nil {
		return nil, customError.NewError("Missing logger", "Не указана система логирования")
	}

	return &kafkaSubscriber{
		config:  config,
		logger:  logger,
		decoder: decoder,
		closing: make(chan struct{}),
	}, nil
}

func (subscriber *kafkaSubscriber) Subscribe(ctx context.Context, topic string) (<-chan interface{}, error) {
	if subscriber.closed {
		return nil, customError.NewError("Kafka subscriber already closed", "Kafka подписчик закрыт")
	}

	subscriber.logger.Infow("Subscribing to Kafka topic",
		"topic", topic,
		"consumer_group", subscriber.config.ConsumerGroup,
	)

	subscriber.subscribersWg.Add(1)

	outputCh := make(chan interface{}, 0)

	consumeClosed, err := subscriber.consumeMessages(ctx, topic, outputCh)
	if err != nil {
		subscriber.subscribersWg.Done()

		return nil, err
	}

	go func() {
		subscriber.handleReconnects(ctx, topic, outputCh, consumeClosed)
		close(outputCh)
		subscriber.subscribersWg.Done()
	}()

	return outputCh, nil
}

func (subscriber *kafkaSubscriber) Close() error {
	if subscriber.closed {
		return nil
	}

	subscriber.closed = true

	close(subscriber.closing)

	subscriber.subscribersWg.Wait()

	subscriber.logger.Info("Kafka subscriber closed")

	return nil
}

func (subscriber *kafkaSubscriber) handleReconnects(ctx context.Context, topic string, outputCh chan interface{}, consumeClosed chan struct{}) {
	for {
		if consumeClosed != nil {
			<-consumeClosed
			subscriber.logger.Infow("consumeMessages stopped",
				"topic", topic,
				"consumer_group", subscriber.config.ConsumerGroup,
			)
		} else {
			subscriber.logger.Infow("consumeClosed empty",
				"topic", topic,
				"consumer_group", subscriber.config.ConsumerGroup,
			)
		}

		select {
		case <-subscriber.closing:
			subscriber.logger.Infow("Closing subscriber, no reconnect needed",
				"topic", topic,
				"consumer_group", subscriber.config.ConsumerGroup,
			)

			return
		case <-ctx.Done():
			subscriber.logger.Infow("Ctx cancelled, no reconnect needed",
				"topic", topic,
				"consumer_group", subscriber.config.ConsumerGroup,
			)

			return
		default:
			subscriber.logger.Infow("Not closing, reconnecting",
				"topic", topic,
				"consumer_group", subscriber.config.ConsumerGroup,
			)
		}

		subscriber.logger.Infow("Reconnecting consumer",
			"topic", topic,
			"consumer_group", subscriber.config.ConsumerGroup,
		)

		var err error

		consumeClosed, err = subscriber.consumeMessages(ctx, topic, outputCh)
		if err != nil {
			subscriber.logger.Errorw("Cannot reconnect messages consumer: ", err,
				"topic", topic,
				"consumer_group", subscriber.config.ConsumerGroup,
			)

			if subscriber.config.ReconnectRetry != NoSleep {
				time.Sleep(subscriber.config.ReconnectRetry)
			}

			continue
		}
	}
}

func (subscriber *kafkaSubscriber) consumeMessages(ctx context.Context, topic string, outputCh chan interface{}) (chan struct{}, error) {
	subscriber.logger.Infow("Starting consuming",
		"topic", topic,
		"consumer_group", subscriber.config.ConsumerGroup,
	)
	var err error

	consumeMessagesClosed := make(chan struct{})

	client, err := sarama.NewClient(subscriber.config.Brokers, subscriber.config.SaramaConfig)
	if err != nil {
		return nil, customError.NewErrorf("Cannot create kafka client: %v", "Не удается создать kafka клиент: %v", err)
	}

	ctx, cancel := context.WithCancel(ctx)

	if subscriber.config.ConsumerGroup == "" {
		consumeMessagesClosed, err = subscriber.consumeWithoutConsumerGroups(ctx, client, topic, outputCh)
	} else {
		consumeMessagesClosed, err = subscriber.consumeGroupMessages(ctx, client, topic, outputCh)
	}

	if err != nil {
		subscriber.logger.Error("Starting consume failed, cancelling context. Reason: ", err)
		cancel()

		return nil, err
	}

	go func() {
		select {
		case <-subscriber.closing:
			subscriber.logger.Infow("Closing subscriber, cancelling consumeMessages",
				"topic", topic,
				"consumer_group", subscriber.config.ConsumerGroup,
			)
			cancel()
		case <-ctx.Done():
		}
	}()

	go func() {
		<-consumeMessagesClosed
		err := client.Close()
		if err != nil {
			subscriber.logger.Errorw("Cannot close client: ", err,
				"topic", topic,
				"consumer_group", subscriber.config.ConsumerGroup,
			)
		} else {
			subscriber.logger.Infow("Client closed",
				"topic", topic,
				"consumer_group", subscriber.config.ConsumerGroup,
			)
		}
	}()

	return consumeMessagesClosed, nil
}

func (subscriber *kafkaSubscriber) consumeWithoutConsumerGroups(ctx context.Context, client sarama.Client, topic string, outputCh chan interface{}) (chan struct{}, error) {
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, customError.NewErrorf("Cannot create consumer: %v", "Не удалось создать потребителя: %v", err)
	}

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		return nil, customError.NewErrorf("Cannot get partition: %v", "Не удалось получить раздел: %v", err)
	}

	subscriber.logger.Debug("partitions: ", partitions)

	partitionConsumersWg := &sync.WaitGroup{}

	for _, partition := range partitions {
		partitionConsumer, err := consumer.ConsumePartition(topic, partition, subscriber.config.SaramaConfig.Consumer.Offsets.Initial)
		if err != nil {
			if err := client.Close(); err != nil && !errors.Is(err, sarama.ErrClosedClient) {
				subscriber.logger.Errorw("Cannot close kafka client: ", err,
					"topic", topic,
					"consumer_group", subscriber.config.ConsumerGroup,
					"kafka_partition", partition,
				)
			}

			return nil, customError.NewErrorf("Failed to start consumer for partition: %v", "Не удалось запустить потребителя для раздела: %v", err)
		}

		messageHandler := subscriber.createMessagesHandler(outputCh)

		partitionConsumersWg.Add(1)

		go subscriber.consumePartition(ctx, partitionConsumer, messageHandler, partitionConsumersWg, partition)
	}

	closed := make(chan struct{})

	go func() {
		partitionConsumersWg.Wait()
		close(closed)
	}()

	return closed, nil
}

func (subscriber *kafkaSubscriber) consumeGroupMessages(
	ctx context.Context,
	client sarama.Client,
	topic string,
	output chan interface{},
) (chan struct{}, error) {

	group, err := sarama.NewConsumerGroupFromClient(subscriber.config.ConsumerGroup, client)
	if err != nil {
		return nil, customError.NewErrorf("Cannot create consumer group client", "Cannot create consumer group client", err)
	}

	groupClosed := make(chan struct{})

	groupHandlerErrorsCtx, cancelGroupHandlerErrors := context.WithCancel(context.Background())
	groupHandlerErrorsDoneCh := subscriber.groupErrorsHandler(groupHandlerErrorsCtx, group)

	handler := consumerGroupHandler{
		ctx:            ctx,
		messageHandler: subscriber.createMessagesHandler(output),
		logger:         subscriber.logger,
		closing:        subscriber.closing,
	}

	go func() {
		err := group.Consume(ctx, []string{topic}, handler)

		if err != nil {
			if err == sarama.ErrUnknown {
				subscriber.logger.Info("Received unknown Sarama error: ", err)
			} else {
				subscriber.logger.Error("Group consume error: ", err)
			}
		} else {
			subscriber.logger.Debug("Consume stopped without any error")
		}

		cancelGroupHandlerErrors()
		<-groupHandlerErrorsDoneCh

		if err := group.Close(); err != nil {
			subscriber.logger.Info("Group close with error: ", err)
		}

		subscriber.logger.Info("Consuming done")
		close(groupClosed)
	}()

	return groupClosed, nil

}

func (subscriber *kafkaSubscriber) groupErrorsHandler(ctx context.Context, group sarama.ConsumerGroup) chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)

		errs := group.Errors()

		for {
			select {
			case err := <-errs:
				if err == nil {
					continue
				}

				subscriber.logger.Error("Sarama internal error: ", err)

			case <-ctx.Done():

				return
			}
		}
	}()

	return done
}

func (subscriber *kafkaSubscriber) createMessagesHandler(outputCh chan interface{}) *messageHandler {
	return &messageHandler{
		outputCh:  outputCh,
		decoder:   subscriber.decoder,
		nackReset: subscriber.config.NackResend,
		logger:    subscriber.logger,
		closing:   subscriber.closing,
	}
}

func (subscriber *kafkaSubscriber) consumePartition(
	ctx context.Context,
	partitionConsumer sarama.PartitionConsumer,
	messageHandler *messageHandler,
	partitionConsumersWg *sync.WaitGroup,
	partition int32,
) {

	defer func() {
		err := partitionConsumer.Close()
		if err != nil {
			subscriber.logger.Errorf("Cannot close partition consumer: %v, partition: %v", err, partition)
		}

		partitionConsumersWg.Done()

		subscriber.logger.Info("consumePartition stopped. Partition: ", partition)
	}()

	kafkaMessages := partitionConsumer.Messages()

	for {
		select {
		case kafkaMsg := <-kafkaMessages:
			if kafkaMsg == nil {
				subscriber.logger.Info("kafkaMsg is closed, stopping consumePartition. Partition: ", partition)

				return
			}

			err := messageHandler.processMessage(ctx, kafkaMsg)
			if err != nil {
				subscriber.logger.Error(err)

				return
			}
		case <-subscriber.closing:
			subscriber.logger.Debug("Subscriber is closing, stopping consumePartition")

			return

		case <-ctx.Done():
			subscriber.logger.Debug("Ctx was cancelled, stopping consumePartition")

			return
		}
	}
}
