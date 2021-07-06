package application

import (
	service "github.com/Ponchitos/application_service/server/internal/broker"
	"github.com/Ponchitos/application_service/server/pkg"
	"github.com/Ponchitos/application_service/server/tools/logger"
)

type broker struct {
	publisher pkg.Publisher
	lgr       logger.Logger
}

func NewApplicationServiceBroker(publisher pkg.Publisher, lgr logger.Logger) service.ApplicationBroker {
	return &broker{publisher, lgr}
}
