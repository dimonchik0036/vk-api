package vkapi

import (
	"net/url"
)

const (
	defaultVersion = "5.65"
	defaultScheme  = "https"
	defaultHost    = "api.vk.com"
	defaultPath    = "method"

	defaultHTTPS    = "1"
	defaultLanguage = "en"
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

var DefaultApiClient = &ApiClient{
	defaultHTTPClient(),
	defaultVersion,
	AccessToken{},
	"",
	defaultHTTPS,
	defaultLanguage,
}

func ApiUrl() (url url.URL) {
	url.Host = defaultHost
	url.Path = defaultPath
	url.Scheme = defaultScheme
	return url
}
