package provider

import (
	"context"
	"fmt"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Ponchitos/application_service/server/tools/common"
	httpHelper "github.com/Ponchitos/application_service/server/tools/http"
	"net/http"
	"time"
)

func (sp *someProvider) GetKey(ctx context.Context, enterpriseID string) (string, error) {
	token, needUpdate, err := sp.cache.GetValue(ctx, enterpriseID)
	if err != nil {
		return sp.getAuthToken(ctx, enterpriseID)
	}

	if len(token) > 0 && !needUpdate {
		return token, nil
	}

	if len(token) > 0 && needUpdate {
		locked := sp.cache.CheckLock(ctx, enterpriseID)
		if locked {
			return token, nil
		}

		return sp.updateCacheValue(ctx, enterpriseID)
	}

	locked := sp.cache.CheckLock(ctx, enterpriseID)
	if !locked {
		return sp.updateCacheValue(ctx, enterpriseID)
	}

	time.Sleep(time.Second)

	token, _, err = sp.cache.GetValue(ctx, enterpriseID)

	if err != nil {
		return token, err
	}

	if len(token) == 0 {
		return token, errors.NewError("Auth token don't available", "Авторизационный токен недоступен")
	}

	return token, nil
}

func (sp *someProvider) getAuthToken(ctx context.Context, enterpriseID string) (string, error) {
	response, status, err := httpHelper.ExecuteHTTPRequest(ctx, sp.getAuthTokenURL(enterpriseID), "GET", nil, time.Second*10)
	if err != nil {
		return "", err
	}

	if status != http.StatusOK {
		response, err := sp.responseHandler(response)
		if err != nil {
			return "", err
		}

		return "", errors.NewError(response, response)
	}

	token, err := sp.responseHandler(response)

	return token, err
}

func (sp *someProvider) getAuthTokenURL(enterpriseID string) string {
	return fmt.Sprintf("%v?enterpriseId=%v", sp.config.KeyManagerServiceURL, enterpriseID)
}

func (sp *someProvider) updateCacheValue(ctx context.Context, enterpriseID string) (string, error) {
	lockValue := common.RandStringBytes(20)

	sp.Lock()
	defer sp.Unlock()
	defer sp.cache.Unlock(ctx, enterpriseID, lockValue)

	err := sp.cache.Lock(ctx, enterpriseID, lockValue)

	if err != nil {
		return "", err
	}

	token, err := sp.getAuthToken(ctx, enterpriseID)
	if err != nil {
		return token, err
	}

	err = sp.cache.SetValue(ctx, enterpriseID, token)
	if err != nil {
		return token, err
	}

	return token, nil
}
