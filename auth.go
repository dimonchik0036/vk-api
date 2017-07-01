package vkapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
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

// OAuthUrl return standard url for interacting with authentication server.
func OAuthUrl() (url url.URL) {
	url.Scheme = oAuthScheme
	url.Host = oAuthHost
	url.Path = oAuthPath
	return url
}

// Application allows you to interact with authentication server.
type Application struct {
	// GrantType - Authorization type, must be equal to `password`
	GrantType string `json:"grant_type" url:"grant_type,omitempty"`

	// ClientId - Id of your application
	ClientId string `json:"client_id" url:"client_id,omitempty"`

	// ClientSecret - Secret key of your application
	ClientSecret string `json:"client_secret" url:"client_secret,omitempty"`

	// Username - User username
	Username string `json:"username" url:"username,omitempty"`

	// Password - User password
	Password string `json:"password" url:"password,omitempty"`

	// Scope - Access rights required by the application
	Scope int64 `json:"scope" url:"scope,omitempty"`
}

// Values returns values from this Application.
func (app *Application) Values() (values url.Values) {
	values = url.Values{}
	values.Set(paramGrantType, app.GrantType)
	values.Set(paramClientId, app.ClientId)
	values.Set(paramClientSecret, app.ClientSecret)
	values.Set(paramUsername, app.Username)
	values.Set(paramPassword, app.Password)
	values.Set(paramScope, strconv.FormatInt(app.Scope, 10))

	return
}

//  AccessToken allows you to interact with API methods.
type AccessToken struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int64  `json:"expires_in"`
	UserID           int    `json:"user_id"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	RedirectUri      string `json:"redirect_uri"`
	CaptchaSid       string `json:"captcha_sid"`
	CaptchaImg       string `json:"captcha_img"`
	ValidationType   string `json:"validation_type"` //2fa_sms 2fa_app
	PhoneMask        string `json:"phone_mask"`
}

// NewApplication creates a new Application instance.
func NewApplication(username string, password string, scope int64) (app Application) {
	app.GrantType = "password"
	app.Username = username
	app.Password = password
	app.Scope = scope
	app.ClientId = defaultClientId
	app.ClientSecret = defaultClientSecret

	return
}

// Authenticate authenticates *ApiClient through Application.
// If the outcome is successful, it returns a *AccessToken.
func Authenticate(client *ApiClient, app Application) (token *AccessToken, err error) {
	token = new(AccessToken)
	if client.httpClient == nil {
		return nil, errors.New("HttpClient not found")
	}
	auth := OAuthUrl()

	q := ConcatValues(false, auth.Query(), app.Values(), client.Values())
	//q.Set("test_redirect_uri", "1")
	if q != nil {
		auth.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(oAuthMethod, auth.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	/*if res.StatusCode != http.StatusOK {
		return AccessToken{}, errors.New("StatusCode != StatusOK")
	}*/

	err = json.NewDecoder(res.Body).Decode(token)
	if err != nil {
		return nil, err
	}

	return
}
