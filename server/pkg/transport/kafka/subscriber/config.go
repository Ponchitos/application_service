package subscriber

import (
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Shopify/sarama"
	"time"
)

const (
	NoSleep               time.Duration = -1
	DefaultNackResend                   = time.Millisecond * 100
	DefaultReconnectRetry               = time.Second
)

func DefaultSaramaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Consumer.Return.Errors = true
	config.ClientID = "default_client"
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	return config
}

type ConfigSubscriber struct {
	Brokers []string

	SaramaConfig *sarama.Config

	ConsumerGroup string

	NackResend     time.Duration
	ReconnectRetry time.Duration

	TopicDetails *sarama.TopicDetail
}

func (config *ConfigSubscriber) setDefaults() {
	if config.SaramaConfig == nil {
		config.SaramaConfig = DefaultSaramaConfig()
	}

	if config.NackResend == 0 {
		config.NackResend = DefaultNackResend
	}

	if config.ReconnectRetry == 0 {
		config.ReconnectRetry = DefaultReconnectRetry
	}
}

func (config *ConfigSubscriber) validate() error {
	if len(config.Brokers) == 0 {
		return errors.NewError("Missing brokers", "Не указаны брокеры")
	}

	return nil
}
