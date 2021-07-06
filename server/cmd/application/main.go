package main

import (
	"github.com/Ponchitos/application_service/server/app"
	"github.com/Ponchitos/application_service/server/config"
	"github.com/Ponchitos/application_service/server/infrastructure/logger"
)

func main() {
	conf := config.NewConfig()

	lgr := logger.NewLogger(conf.LogLevel, conf.LogLevel == -1)

	err := conf.ReadAllSettings()
	if err != nil {
		lgr.Fatal(err)
	}

	if err := app.Start(conf, lgr); err != nil {
		lgr.Fatal(err)
	}
}
