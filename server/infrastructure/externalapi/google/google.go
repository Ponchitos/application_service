package google

import (
	"context"
	"fmt"
	"github.com/Ponchitos/application_service/server/config"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Ponchitos/application_service/server/infrastructure/keymanager"
	"github.com/Ponchitos/application_service/server/internal/externalapi"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	htransport "google.golang.org/api/transport/http"
	"net/http"
	"net/url"
)

type external struct {
	manager keymanager.Manager
	config  *config.Config
}

func NewExternalGoogleProvider(manager keymanager.Manager, config *config.Config) externalapi.ExternalAPI {
	return &external{
		manager: manager,
		config:  config,
	}
}

func (e *external) getOptions(ctx context.Context, enterpriseID string) (option.ClientOption, error) {
	token, err := e.manager.GetKey(ctx, enterpriseID)
	if err != nil {
		return nil, err
	}

	if e.config.UseProxy {
		return e.initOptionsWithProxy(ctx, token)
	}

	return e.getOptionsByAccessToken(token), nil
}

func (e *external) initOptionsWithProxy(ctx context.Context, token string) (option.ClientOption, error) {
	if e.config.UserProxy == "" || e.config.PassProxy == "" {
		return nil, errors.NewError("Not valid proxy user or password", "Не корректный логин или пароль для прокси сервера")
	}

	if e.config.SchemaProxy != "http" && e.config.SchemaProxy != "https" {
		return nil, errors.NewError("Not valid schema proxy (valid: http or https)", "Не корректная схема запроса прокси сервера (поддерживается: http или https)")
	}

	scopes := option.WithScopes("https://www.googleapis.com/auth/androidmanagement")
	accessTokenOption := e.getOptionsByAccessToken(token)

	transport, err := htransport.NewTransport(ctx, e.initTransport(), scopes, accessTokenOption)
	if err != nil {
		return nil, errors.NewErrorf("Cannot create new proxy transport: %v", "Не удалось создать прокси транспорт: %v", err)
	}

	return option.WithHTTPClient(&http.Client{
		Transport: transport,
	}), nil
}

func (e *external) initTransport() http.RoundTripper {
	transport := http.DefaultTransport.(*http.Transport).Clone()

	transport.MaxIdleConnsPerHost = 100

	transport.Proxy = http.ProxyURL(&url.URL{
		Scheme: e.config.SchemaProxy,
		User:   url.UserPassword(e.config.UserProxy, e.config.PassProxy),
		Host:   fmt.Sprintf("%v:%v", e.config.HostProxy, e.config.PortProxy),
	})

	return transport
}

func (e *external) getOptionsByAccessToken(accessToken string) option.ClientOption {
	token := &oauth2.Token{
		AccessToken: accessToken,
	}

	return option.WithTokenSource(oauth2.StaticTokenSource(token))
}
