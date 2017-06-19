package vkapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	oAuthScheme = "https"
	oAuthHost   = "oauth.vk.com"
	oAuthPath   = "token"
	oAuthMethod = "GET"
	//defaultClientId = "2274003"                  //VK for Android
	//defaultClientSecret = "hHbZxrka2uZ6jB1inYsH" //VK for Android
	//defaultClientId = "3140623"                  //VK for iPhone
	//defaultClientSecret = "VeWdmVclDCtn6ihuP1nt" //VK for iPhone
	defaultClientId     = "3697615"              //VK for Windows
	defaultClientSecret = "AlVXZFMUqyrnABp8ncuU" //VK for Windows

	paramGrantType    = "grant_type"
	paramClientId     = "client_id"
	paramClientSecret = "client_secret"
	paramUsername     = "username"
	paramPassword     = "password"
	paramScope        = "scope"

	paramPhoneMask      = "phone_mask"
	paramValidationType = "validation_type"
	paramCaptchaSid     = "captcha_sid"
	paramCode           = "code"
	paramForceSms       = "force_sms"
	paramNeedCaptcha    = "need_captcha"
	paramNeedValidation = "need_validation"
	paramInvalidClient  = "invalid_client"
)

func OAuthUrl() (url url.URL) {
	url.Scheme = oAuthScheme
	url.Host = oAuthHost
	url.Path = oAuthPath
	return url
}

type Application struct {
	// GrantType - Authorization type, must be equal to `password`
	GrantType string `json:"grant_type"`

	// ClientId - Id of your application
	ClientId string `json:"client_id"`

	// ClientSecret - Secret key of your application
	ClientSecret string `json:"client_secret"`

	// Username - User username
	Username string `json:"username"`

	// Password - User password
	Password string `json:"password"`

	// Scope - Access rights required by the application
	Scope int64 `json:"scope"`
}

type AccessToken struct {
	AccessToken      string        `json:"access_token"`
	ExpiresIn        time.Duration `json:"expires_in"`
	UserID           int           `json:"user_id"`
	Error            string        `json:"error"`
	ErrorDescription string        `json:"error_description"`
	RedirectUri      string        `json:"redirect_uri"`
	CaptchaSid       string        `json:"captcha_sid"`
	CaptchaImg       string        `json:"captcha_img"`
	ValidationType   string        `json:"validation_type"` //2fa_sms 2fa_app
	PhoneMask        string        `json:"phone_mask"`
}

func DefaultApplication(username string, password string, scope int64) (application Application) {
	application.GrantType = "password"
	application.Username = username
	application.Password = password
	application.Scope = scope
	application.ClientId = defaultClientId
	application.ClientSecret = defaultClientSecret

	return
}

func Authenticate(client *ApiClient, application Application) (token AccessToken, err error) {
	if client.httpClient == nil {
		return AccessToken{}, errors.New("HttpClient not found")
	}

	auth := OAuthUrl()
	q := auth.Query()
	q.Set(paramGrantType, application.GrantType)
	q.Set(paramClientId, application.ClientId)
	q.Set(paramClientSecret, application.ClientSecret)
	q.Set(paramUsername, application.Username)
	q.Set(paramPassword, application.Password)
	q.Set(paramScope, strconv.FormatInt(application.Scope, 10))
	//q.Set("test_redirect_uri", "1")

	q.Set(paramVersion, client.ApiVersion)
	q.Set(paramLanguage, client.Language)
	q.Set(paramHTTPS, client.HTTPS)
	auth.RawQuery = q.Encode()

	req, err := http.NewRequest(oAuthMethod, auth.String(), nil)
	if err != nil {
		return AccessToken{}, err
	}

	res, err := client.httpClient.Do(req)
	if err != nil {
		return AccessToken{}, err
	}

	/*if res.StatusCode != http.StatusOK {
		return AccessToken{}, errors.New("StatusCode != StatusOK")
	}*/

	err = json.NewDecoder(res.Body).Decode(&token)
	if err != nil {
		return AccessToken{}, err
	}

	return
}
