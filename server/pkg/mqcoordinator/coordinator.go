package mqcoordinator

import (
	"context"
	"github.com/Ponchitos/application_service/server/pkg"
	"github.com/Ponchitos/application_service/server/pkg/mqcoordinator/coordinator"
	"github.com/Ponchitos/application_service/server/tools/logger"
)

type MsCoordinator interface {
	Logger() logger.Logger
	AddHandler(
		handlerName string,
		subscriberTopic string,
		subscriber pkg.Subscriber,
		publisherTopic string,
		publisher pkg.Publisher,
		handlerFunc pkg.Endpoint,
		decoder coordinator.DecoderFunc,
		encoder coordinator.EncoderFunc,
	)
	AddWithoutResponseHandler(
		handlerName string,
		subscriberTopic string,
		subscriber pkg.Subscriber,
		handlerFunc pkg.WithoutResponseEndpoint,
		decoder coordinator.DecoderFunc,
		encoder coordinator.EncoderFunc,
	)
	Close() error
	Running() chan struct{}
	Run(ctx context.Context) error
}
