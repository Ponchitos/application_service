package redis

import (
	"github.com/Ponchitos/application_service/server/config"
	"github.com/Ponchitos/application_service/server/infrastructure/keymanager/provider"
	"github.com/Ponchitos/application_service/server/tools/logger"
	"github.com/go-redis/redis/v8"
)

type client struct {
	conn *redis.Client
	lgr  logger.Logger
	conf *config.Config
}

func NewRedisClient(conf *config.Config, lgr logger.Logger) provider.SomeProviderCache {
	conn := redis.NewClient(&redis.Options{})

	return &client{
		conn: conn,
		lgr:  lgr,
		conf: conf,
	}
}
