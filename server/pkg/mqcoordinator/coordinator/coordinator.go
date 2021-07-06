package coordinator

import (
	"context"
	"fmt"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Ponchitos/application_service/server/pkg"
	"github.com/Ponchitos/application_service/server/tools/common"
	"github.com/Ponchitos/application_service/server/tools/logger"
	syncinternal "github.com/Ponchitos/application_service/server/tools/sync"
	"runtime/debug"
	"sync"
	"time"
)

type CoordinatorConfig struct {
	CloseTimeout time.Duration
}

func (config *CoordinatorConfig) setDefault() {
	if config.CloseTimeout == 0 {
		config.CloseTimeout = time.Second * 30
	}
}

type Coordinator struct {
	config CoordinatorConfig

	handlers map[string]*handler

	handlersWg        *sync.WaitGroup
	runningHandlersWg *sync.WaitGroup

	closedLock sync.Mutex

	closeCh   chan struct{}
	closedCh  chan struct{}
	runningCh chan struct{}

	closed    bool
	isRunning bool

	logger logger.Logger
}

func NewCoordinator(config CoordinatorConfig, logger logger.Logger) (*Coordinator, error) {
	config.setDefault()

	return &Coordinator{
		config: config,

		handlers: make(map[string]*handler),

		handlersWg:        &sync.WaitGroup{},
		runningHandlersWg: &sync.WaitGroup{},

		closeCh:   make(chan struct{}),
		closedCh:  make(chan struct{}),
		runningCh: make(chan struct{}),

		logger: logger,
	}, nil
}

func (coordinator *Coordinator) Logger() logger.Logger {
	return coordinator.logger
}

func (coordinator *Coordinator) AddHandler(
	handlerName string,
	subscriberTopic string,
	subscriber pkg.Subscriber,
	publisherTopic string,
	publisher pkg.Publisher,
	handlerFunc pkg.Endpoint,
	decoder DecoderFunc,
	encoder EncoderFunc,
) {
	coordinator.logger.Infow("AddHandler: ",
		"handler_name", handlerName,
		"subscriber_topic", subscriberTopic,
		"publisher_topic", publisherTopic,
	)

	if _, ok := coordinator.handlers[handlerName]; ok {
		panic(errors.NewError("Duplicate handler", "Обработчик с таким именем уже существует"))
	}

	publisherName, subscriberName := common.StructName(publisher), common.StructName(subscriber)

	newHandler := &handler{
		name: handlerName,

		logger: coordinator.logger,

		subscriber:      subscriber,
		subscriberTopic: subscriberTopic,
		subscriberName:  subscriberName,

		publisher:      publisher,
		publisherTopic: publisherTopic,
		publisherName:  publisherName,

		handlerFunc: handlerFunc,

		decoder: decoder,
		encoder: encoder,

		runningHandlersWg: coordinator.runningHandlersWg,

		closeCh: coordinator.closeCh,
	}

	coordinator.handlers[handlerName] = newHandler
}

func (coordinator *Coordinator) AddWithoutResponseHandler(
	handlerName string,
	subscriberTopic string,
	subscriber pkg.Subscriber,
	handlerFunc pkg.WithoutResponseEndpoint,
	decoder DecoderFunc,
	encoder EncoderFunc,
) {

	handlerAdapter := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return nil, handlerFunc(ctx, request)
	}

	coordinator.AddHandler(handlerName, subscriberTopic, subscriber, "", nil, handlerAdapter, decoder, encoder)
}

func (coordinator *Coordinator) isClosed() bool {
	coordinator.closedLock.Lock()
	defer coordinator.closedLock.Unlock()

	return coordinator.closed
}

func (coordinator *Coordinator) Close() error {
	coordinator.closedLock.Lock()
	defer coordinator.closedLock.Unlock()

	if coordinator.closed {
		return nil
	}

	coordinator.closed = true

	coordinator.logger.Info("Closing coordinator")
	defer coordinator.logger.Info("Coordinator closed")

	close(coordinator.closeCh)
	close(coordinator.closedCh)

	timeout := syncinternal.WaitGroupTimeout(coordinator.handlersWg, coordinator.config.CloseTimeout)

	if timeout {
		return errors.NewError("Coordinator close timeout", "Истекло время закрытия координатора")
	}

	return nil
}

func (coordinator *Coordinator) closeWhenAllHandlersStopped() {
	coordinator.handlersWg.Wait()

	if coordinator.isClosed() {
		return
	}

	coordinator.logger.Error(errors.NewError("All router handlers stopped", "Все обработчики остановлены"))

	if err := coordinator.Close(); err != nil {
		coordinator.logger.Error(errors.NewError("Cannot close coordinator", "Не удалось закрыть координатор"))
	}

}

func (coordinator *Coordinator) Running() chan struct{} {
	return coordinator.runningCh
}

func (coordinator *Coordinator) Run(ctx context.Context) error {
	if coordinator.isRunning {
		return errors.NewError("Coordinator is already running", "Координатор уже запущен")
	}

	coordinator.isRunning = true

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, handler := range coordinator.handlers {
		coordinator.logger.Debugw("Subscribing to topic",
			"subscriber_name", handler.name,
			"topic", handler.subscriberTopic,
		)

		messages, err := handler.subscriber.Subscribe(ctx, handler.subscriberTopic)
		if err != nil {
			return errors.NewErrorf("Cannot subscribe topic %s: %s", "Не удается подписаться к топику %s: %s", handler.subscriberTopic, err)
		}

		handler.messageCh = messages
	}

	for _, handler := range coordinator.handlers {
		coordinator.handlersWg.Add(1)

		targetHandler := handler
		go func() {
			targetHandler.run()

			coordinator.handlersWg.Done()
			coordinator.logger.Infow("Subscriber stopped",
				"subscriber_name", targetHandler.name,
				"topic", targetHandler.subscriberTopic,
			)
		}()
	}

	close(coordinator.runningCh)

	go coordinator.closeWhenAllHandlersStopped()

	<-coordinator.closeCh
	cancel()

	coordinator.logger.Infow("Waiting for messages",
		"timeout", coordinator.config.CloseTimeout,
	)

	<-coordinator.closedCh

	coordinator.logger.Info("All messages processed", nil)

	return nil
}

type handler struct {
	name string

	logger logger.Logger

	subscriber      pkg.Subscriber
	subscriberTopic string
	subscriberName  string

	publisher      pkg.Publisher
	publisherTopic string
	publisherName  string

	handlerFunc pkg.Endpoint

	decoder DecoderFunc
	encoder EncoderFunc

	runningHandlersWg *sync.WaitGroup

	messageCh <-chan interface{}

	closeCh chan struct{}
}

func (hr *handler) run() {
	hr.logger.Infow("Starting handler",
		"name", hr.name,
		"subscriber_name", hr.subscriberName,
		"subscriber_topic", hr.subscriberTopic,
		"publisher_name", hr.publisherName,
		"publisher_topic", hr.publisherTopic,
	)

	go hr.handleClose()

	for msg := range hr.messageCh {
		hr.runningHandlersWg.Add(1)
		go hr.handleMessage(msg, hr.handlerFunc, hr.decoder, hr.encoder)
	}

	if hr.publisher != nil {
		hr.logger.Debug("Waiting for publisher to close")
		if err := hr.publisher.Close(); err != nil {
			hr.logger.Error("Failed to close publisher: ", err)
		}
		hr.logger.Debug("Publisher closed")
	}

	hr.logger.Debug("Router handler stopped", nil)
}

func (hr *handler) handleClose() {
	<-hr.closeCh

	hr.logger.Info("Waiting for subscriber to close")

	if err := hr.subscriber.Close(); err != nil {
		hr.logger.Error("Failed to close subscriber: ", err)
	}

	hr.logger.Debug("Subscriber closed")
}

func (hr *handler) handleMessage(msg interface{}, endpoint pkg.Endpoint, decoder DecoderFunc, encoder EncoderFunc) {
	defer hr.runningHandlersWg.Done()

	ctx := context.Background()

	defer func() {
		if recovered := recover(); recovered != nil {
			hr.logger.Error(
				"Panic recovered in handler. Stack: "+string(debug.Stack()),
				errors.NewErrorf("%s", "%s", recovered),
				msg,
			)
		}
	}()

	hr.logger.Info("Received message: ", msg)

	request, err := decoder(ctx, msg)
	if err != nil {
		hr.logger.Error("Decoder returned error: ", err)

		return
	}

	response, err := endpoint(ctx, request)
	if err != nil {
		hr.logger.Error("Handler returned error: ", err)

		return
	}

	messages, err := encoder(ctx, response)
	if err != nil {
		hr.logger.Error("Encoder returned error: ", err)

		return
	}

	if err := hr.publishProducedMessages(messages); err != nil {
		hr.logger.Error("Publishing produced messages failed: ", err)

		return
	}

	hr.logger.Info("Message handled: ", msg)
}

func (hr *handler) publishProducedMessages(messages []interface{}) error {
	if len(messages) == 0 {
		return nil
	}

	if hr.publisher == nil {
		return errors.NewError("Returned output messages in a handler without publisher", "Не определен отправитель")
	}

	hr.logger.Infow("Sending produced messages",
		"produced_messages_count", len(messages),
		"publish_topic", hr.publisherTopic,
	)

	for _, msg := range messages {
		if err := hr.publisher.Publish(hr.publisherTopic, msg); err != nil {
			hr.logger.Errorw("Cannot publish message: ", err,
				"not_sent_message", fmt.Sprintf("%#v", msg),
			)

			return err
		}
	}

	return nil
}
