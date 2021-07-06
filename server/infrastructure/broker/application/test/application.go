package test

import (
	service "github.com/Ponchitos/application_service/server/internal/broker"
	"github.com/Ponchitos/application_service/server/tools/logger"
)

type testBroker struct {
	lgr logger.Logger
}

func NewTestApplicationServiceBroker(lgr logger.Logger) service.ApplicationBroker {
	return &testBroker{lgr}
}
