package provider

import (
	"context"
	"encoding/json"
	"github.com/Ponchitos/application_service/server/config"
	"github.com/Ponchitos/application_service/server/infrastructure/keymanager"
	"github.com/Ponchitos/application_service/server/tools/logger"
	"sync"
)

type SomeProviderCache interface {
	SetValue(ctx context.Context, key string, value interface{}) error
	GetValue(ctx context.Context, key string) (value string, needUpdateValue bool, err error)
	Lock(ctx context.Context, key, lock string) error
	CheckLock(ctx context.Context, key string) bool
	Unlock(ctx context.Context, key, lock string)
}

type someProvider struct {
	sync.Mutex

	config *config.Config
	lgr    logger.Logger
	cache  SomeProviderCache
}

type httpResponse struct {
	Status   string      `json:"status"`
	Response interface{} `json:"response"`
}

func NewKeyManagerProvider(config *config.Config, lgr logger.Logger, cache SomeProviderCache) keymanager.Manager {
	return &someProvider{config: config, lgr: lgr, cache: cache}
}

func (sp *someProvider) responseHandler(data []byte) (string, error) {
	var result *httpResponse

	err := json.Unmarshal(data, &result)
	if err != nil {
		return "", err
	}

	return result.Response.(string), nil
}
