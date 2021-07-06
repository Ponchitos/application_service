package app

import (
	"context"
	"fmt"
	"github.com/Ponchitos/application_service/server/api"
	"github.com/Ponchitos/application_service/server/config"
	"github.com/Ponchitos/application_service/server/infrastructure/broker/application"
	testBroker "github.com/Ponchitos/application_service/server/infrastructure/broker/application/test"
	customError "github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Ponchitos/application_service/server/infrastructure/handlers/apk"
	"github.com/Ponchitos/application_service/server/infrastructure/repositories"
	"github.com/Ponchitos/application_service/server/infrastructure/repositories/postgres"
	testRepository "github.com/Ponchitos/application_service/server/infrastructure/repositories/test"
	"github.com/Ponchitos/application_service/server/internal/broker"
	"github.com/Ponchitos/application_service/server/internal/handlers"
	"github.com/Ponchitos/application_service/server/internal/services/applications"
	"github.com/Ponchitos/application_service/server/messagecontroller"
	"github.com/Ponchitos/application_service/server/pkg"
	"github.com/Ponchitos/application_service/server/pkg/mqcoordinator"
	"github.com/Ponchitos/application_service/server/pkg/transport/kafka/decoder"
	"github.com/Ponchitos/application_service/server/pkg/transport/kafka/encoder"
	"github.com/Ponchitos/application_service/server/pkg/transport/kafka/publisher"
	"github.com/Ponchitos/application_service/server/pkg/transport/kafka/subscriber"
	testDecoder "github.com/Ponchitos/application_service/server/pkg/transport/testtransport/messagebroker/decoder"
	testEncoder "github.com/Ponchitos/application_service/server/pkg/transport/testtransport/messagebroker/encoder"
	testPublisher "github.com/Ponchitos/application_service/server/pkg/transport/testtransport/messagebroker/publisher"
	testSubscriber "github.com/Ponchitos/application_service/server/pkg/transport/testtransport/messagebroker/subscriber"
	"github.com/Ponchitos/application_service/server/tools/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	test  = "test"
	dev   = "dev"
	prod  = "prod"
	local = "local"
)

type app struct {
	config            *config.Config
	logger            logger.Logger
	errChan           chan error
	store             repositories.DataBase
	handler           handlers.ApkFileHandler
	publisher         pkg.Publisher
	subscriber        pkg.Subscriber
	broker            broker.ApplicationBroker
	application       applications.Service
	messageController mqcoordinator.MsCoordinator
	testQueue         chan interface{}
}

func Start(conf *config.Config, lgr logger.Logger) error {
	a := &app{config: conf, logger: lgr, errChan: make(chan error), testQueue: make(chan interface{})}

	a.logger.Debug("Application server starting ...")

	if err := a.configureStorage(); err != nil {
		return err
	}

	if err := a.configureApkHandler(); err != nil {
		return err
	}

	if err := a.configureMessageBrokerPublisher(); err != nil {
		return err
	}

	if err := a.configureMessageBrokerSubscription(); err != nil {
		return err
	}

	if err := a.configureApplicationBroker(); err != nil {
		return err
	}

	a.application = applications.NewApplicationService(lgr, a.store.GetApplicationRepository(), a.handler, conf, a.broker)

	if err := a.configureMessageController(); err != nil {
		return err
	}

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

		a.errChan <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		if a.messageController != nil {
			err := a.messageController.Run(context.Background())
			if err != nil {
				a.errChan <- fmt.Errorf("Cannot run message controller: %s ", err)
			}
		}
	}()

	a.launchAPIServer()

	err := <-a.errChan

	lgr.Info("Exit: ", err)

	a.Stop()

	return err
}

func (a *app) configureStorage() error {
	a.logger.Debug("configuring storage...")

	defer a.logger.Debugf("configured database: %s", a.config.Env)

	var err error

	switch a.config.Env {
	case test:
		a.store, err = testRepository.OpenTestConnection(a.logger)
		if a.logger.FailOnError(err, "failed to connect to Test DB") {
			return err
		}

		return nil
	case local:
		a.store, err = postgres.OpenConnection(a.config, a.logger)
		if a.logger.FailOnError(err, "failed to connect to Postgre DB") {
			return err
		}

		return nil
	case dev:
		a.store, err = postgres.OpenConnection(a.config, a.logger)
		if a.logger.FailOnError(err, "failed to connect to Postgre DB") {
			return err
		}

		return nil
	case prod:
		a.store, err = postgres.OpenConnection(a.config, a.logger)
		if a.logger.FailOnError(err, "failed to connect to Postgre DB") {
			return err
		}

		return nil
	default:
		return customError.NewError("unknown db type", "неверный тип БД")
	}
}

func (a *app) configureApkHandler() error {
	a.logger.Debug("configuring apk handler...")
	defer a.logger.Debugf("configured apk handler: %s", a.config.Env)

	var err error

	switch a.config.Env {
	case test:
		a.handler, err = apk.NewApkHandler(a.config, a.logger)
		if a.logger.FailOnError(err, "Cannot create apk handler") {
			return err
		}

		return nil
	case local:
		a.handler, err = apk.NewApkHandler(a.config, a.logger)
		if a.logger.FailOnError(err, "Cannot create apk handler") {
			return err
		}

		return nil
	case dev:
		a.handler, err = apk.NewApkHandler(a.config, a.logger)
		if a.logger.FailOnError(err, "Cannot create apk handler") {
			return err
		}

		return nil
	case prod:
		a.handler, err = apk.NewApkHandler(a.config, a.logger)
		if a.logger.FailOnError(err, "Cannot create apk handler") {
			return err
		}

		return nil
	default:
		return customError.NewError("unknown handler type", "неверный тип apk обработчика")

	}

}

func (a *app) configureMessageBrokerPublisher() error {
	a.logger.Debug("configuring message broker publisher...")

	defer a.logger.Debugf("configured message broker publisher: %s", a.config.Env)

	var err error

	switch a.config.Env {
	case test:
		a.publisher, err = testPublisher.NewTestPublisher(a.logger, testDecoder.DefaultTestDecoderPublisherFunc, a.testQueue)
		if a.logger.FailOnError(err, "Cannot create message broker publisher") {
			return err
		}

		return nil
	case local:
		a.publisher, err = publisher.NewPublisher(&publisher.ConfigPublisher{Brokers: a.config.KafkaBrokers}, a.logger, decoder.DefaultKafkaDecoderPublisherFunc)
		if a.logger.FailOnError(err, "Cannot create message broker publisher") {
			return err
		}

		return nil
	case dev:
		a.publisher, err = publisher.NewPublisher(&publisher.ConfigPublisher{Brokers: a.config.KafkaBrokers}, a.logger, decoder.DefaultKafkaDecoderPublisherFunc)
		if a.logger.FailOnError(err, "Cannot create message broker publisher") {
			return err
		}

		return nil
	case prod:
		a.publisher, err = publisher.NewPublisher(&publisher.ConfigPublisher{Brokers: a.config.KafkaBrokers}, a.logger, decoder.DefaultKafkaDecoderPublisherFunc)
		if a.logger.FailOnError(err, "Cannot create message broker publisher") {
			return err
		}

		return nil
	default:
		return customError.NewError("unknown message broker publisher type", "неверный тип брокера публикации")

	}
}

func (a *app) configureMessageBrokerSubscription() error {
	a.logger.Debug("configuring message broker subscription...")

	defer a.logger.Debugf("configured message broker subscription: %s", a.config.Env)

	var err error

	switch a.config.Env {
	case test:
		a.subscriber, err = testSubscriber.NewTestSubscriber(a.logger, testEncoder.DefaultTestEncoderSubscriberFunc, a.testQueue)
		if a.logger.FailOnError(err, "Cannot create message subscriber") {
			return err
		}

		return nil
	case local:
		a.subscriber, err = subscriber.NewSubscriber(&subscriber.ConfigSubscriber{Brokers: a.config.KafkaBrokers}, a.logger, encoder.DefaultKafkaEncoderSubscriberFunc)
		if a.logger.FailOnError(err, "Cannot create message subscriber") {
			return err
		}

		return nil
	case dev:
		a.subscriber, err = subscriber.NewSubscriber(&subscriber.ConfigSubscriber{Brokers: a.config.KafkaBrokers}, a.logger, encoder.DefaultKafkaEncoderSubscriberFunc)
		if a.logger.FailOnError(err, "Cannot create message subscriber") {
			return err
		}

		return nil
	case prod:
		a.subscriber, err = subscriber.NewSubscriber(&subscriber.ConfigSubscriber{Brokers: a.config.KafkaBrokers}, a.logger, encoder.DefaultKafkaEncoderSubscriberFunc)
		if a.logger.FailOnError(err, "Cannot create message subscriber") {
			return err
		}

		return nil
	default:
		return customError.NewError("unknown message broker subscription type", "неверный тип брокера подписки")
	}
}

func (a *app) configureApplicationBroker() error {
	a.logger.Debug("configuring application broker...")

	defer a.logger.Debugf("configured application broker: %s", a.config.Env)

	switch a.config.Env {
	case test:
		a.broker = testBroker.NewTestApplicationServiceBroker(a.logger)
		return nil
	case local:
		a.broker = application.NewApplicationServiceBroker(a.publisher, a.logger)
		return nil
	case dev:
		a.broker = application.NewApplicationServiceBroker(a.publisher, a.logger)
		return nil
	case prod:
		a.broker = application.NewApplicationServiceBroker(a.publisher, a.logger)
		return nil
	default:
		return customError.NewError("unknown application broker type", "неверный тип брокера сервиса приложения")
	}
}

func (a *app) configureMessageController() error {
	a.logger.Debug("configuring message controller...")
	defer a.logger.Debugf("configured message controller: %s", a.config.Env)

	var err error

	switch a.config.Env {
	case test:
		return nil
	case local:
		a.messageController, err = messagecontroller.NewMessageController(a.logger, a.subscriber, a.publisher, a.application)
		if a.logger.FailOnError(err, "Cannot create message controller") {
			return err
		}

		return nil
	case dev:
		a.messageController, err = messagecontroller.NewMessageController(a.logger, a.subscriber, a.publisher, a.application)
		if a.logger.FailOnError(err, "Cannot create message controller") {
			return err
		}

		return nil
	case prod:
		a.messageController, err = messagecontroller.NewMessageController(a.logger, a.subscriber, a.publisher, a.application)
		if a.logger.FailOnError(err, "Cannot create message controller") {
			return err
		}

		return nil
	default:
		return customError.NewError("unknown message controller type", "неверный тип контроллера сообщениями")
	}
}

func (a *app) launchAPIServer() {
	a.logger.Debugf("Application service starting on port %v ...", a.config.APIPort)

	httpHandler := api.NewHTTPHandler(a.application, a.logger)

	go func() {
		a.errChan <- http.ListenAndServe(fmt.Sprintf(":%v", a.config.APIPort), httpHandler)
	}()
}

func (a *app) Stop() {
	a.store.Close()
	close(a.errChan)
	close(a.testQueue)
	a.messageController.Close()
	a.publisher.Close()
	a.subscriber.Close()
}
