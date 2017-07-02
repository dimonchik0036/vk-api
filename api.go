package vkapi

import (
	"errors"
	"log"
	"net/url"
	"os"
)

const (
	defaultVersion = "5.65"
	defaultScheme  = "https"
	defaultHost    = "api.vk.com"
	defaultPath    = "method"
	defaultMethod  = "GET"

	defaultHTTPS    = "1"
	defaultLanguage = LangEN

	paramVersion  = "v"
	paramLanguage = "lang"
	paramHTTPS    = "https"
	paramToken    = "access_token"
)

const (
	ErrApiClientNotFound = "ApiClient not found."
)

const (
	LangRU = "ru" //Russian
	LangUA = "ua" //Ukrainian
	LangBE = "be" //Belarusian
	LangEN = "en" //English
	LangES = "es" //Spanish
	LangFI = "fi" //Finnish
	LangDE = "de" //German
	LangIT = "it" //Italian
)

// ApiClient allows you to send requests to API server.
type ApiClient struct {
	httpClient  HTTPClient
	ApiVersion  string
	AccessToken *AccessToken
	secureToken string

	// If Log is true, ApiClient will write logs.
	Log    bool
	Logger *log.Logger

	// HTTPS defines if use https instead of http. 1 - use https. 0 - use http.
	HTTPS string

	// Language define the language in which different data will be returned, for example, names of countries and cities.
	Language string
}

// SetAccessToken sets access token to ApiClient.
func (api *ApiClient) SetAccessToken(token string) {
	api.AccessToken = &AccessToken{token,
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

// Values returns values from this ApiClient.
func (api *ApiClient) Values() (values url.Values) {
	values = url.Values{}
	values.Add(paramVersion, api.ApiVersion)
	values.Add(paramLanguage, api.Language)
	values.Add(paramHTTPS, api.HTTPS)
	return
}

// Authenticate run authentication this ApiClient from Application.
func (api *ApiClient) Authenticate(application Application) (err error) {
	api.AccessToken, err = Authenticate(api, application)
	if err != nil {
		return err
	}

	if api.AccessToken.Error != "" {
		return errors.New(api.AccessToken.Error + " : " + api.AccessToken.ErrorDescription)
	}

	return nil
}

// NewApiClient creates a new *ApiClient instance.
func NewApiClient() *ApiClient {
	client := &ApiClient{
		defaultHTTPClient(),
		defaultVersion,
		nil,
		"",
		false,
		log.New(os.Stdout, "", log.LstdFlags),
		defaultHTTPS,
		defaultLanguage,
	}

	return client
}

// ApiUrl return standard url for interacting with server API.
func ApiUrl() (url url.URL) {
	url.Host = defaultHost
	url.Path = defaultPath
	url.Scheme = defaultScheme
	return url
}

func (api *ApiClient) logPrintf(format string, v ...interface{}) {
	if api.Log {
		api.Logger.Printf(format, v...)
	}
}
