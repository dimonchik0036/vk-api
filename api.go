package vkapi

import (
	"errors"
	"net/url"
)

const (
	defaultVersion = "5.65"
	defaultScheme  = "https"
	defaultHost    = "api.vk.com"
	defaultPath    = "method"

	defaultHTTPS    = "1"
	defaultLanguage = "en"

	paramVersion  = "v"
	paramLanguage = "lang"
	paramHTTPS    = "https"
)

type ApiClient struct {
	httpClient  HTTPClient
	ApiVersion  string
	AccessToken AccessToken
	secureToken string

	// HTTPS defines if use https instead of http. 1 - use https. 0 - use http
	HTTPS string

	// Language defines the language in which different data will be returned, for example, names of countries and cities
	// ru — Russian
	// ua — Ukrainian
	// be — Belarusian
	// en — English
	// es — Spanish
	// fi — Finnish
	// de — German
	// it — Italian
	Language string
}

func (api *ApiClient) SetAccessToken(token string) {
	api.AccessToken = AccessToken{token,
		0,
		0,
		"",
		"",
		"",
		"",
		"",
		"",
		""}
}

func (api *ApiClient) Authenticate(application Application) (err error) {
	api.AccessToken, err = Authenticate(api, application)
	if err != nil {
		return err
	}

	if api.AccessToken.Error != "" {
		return errors.New(api.AccessToken.Error + ":" + api.AccessToken.ErrorDescription)
	}

	return nil
}

func DefaultApiClient() *ApiClient {
	client := &ApiClient{
		defaultHTTPClient(),
		defaultVersion,
		AccessToken{},
		"",
		defaultHTTPS,
		defaultLanguage,
	}

	return client
}

func ApiUrl() (url url.URL) {
	url.Host = defaultHost
	url.Path = defaultPath
	url.Scheme = defaultScheme
	return url
}
