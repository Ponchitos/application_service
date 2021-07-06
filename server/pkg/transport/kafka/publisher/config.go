package publisher

import (
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Shopify/sarama"
	"time"
)

func DefaultSaramaConfig() *sarama.Config {
	config := sarama.NewConfig()

	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true
	config.Version = sarama.V2_6_0_0
	config.Metadata.Retry.Backoff = time.Second * 2
	config.ClientID = "default_client"
	config.Producer.MaxMessageBytes = 2000000

	return config
}

type ConfigPublisher struct {
	Brokers []string

	SaramaConfig *sarama.Config
}

func (config *ConfigPublisher) setDefaults() {
	if config.SaramaConfig == nil {
		config.SaramaConfig = DefaultSaramaConfig()
	}
}

func (config *ConfigPublisher) validate() error {
	if len(config.Brokers) == 0 {
		return errors.NewError("Missing brokers", "Не указан брокер")
	}

	return nil
}
