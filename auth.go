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
	//defaultClientID = "2274003"                  //VK for Android
	//defaultClientSecret = "hHbZxrka2uZ6jB1inYsH" //VK for Android
	//defaultClientID = "3140623"                  //VK for iPhone
	//defaultClientSecret = "VeWdmVclDCtn6ihuP1nt" //VK for iPhone
	defaultClientID     = "3697615"              //VK for Windows
	defaultClientSecret = "AlVXZFMUqyrnABp8ncuU" //VK for Windows

	paramGrantType    = "grant_type"
	paramClientID     = "client_id"
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

// OAuthURL return standard url for interacting with authentication server.
func OAuthURL() (url url.URL) {
	url.Scheme = oAuthScheme
	url.Host = oAuthHost
	url.Path = oAuthPath
	return url
}

// Application allows you to interact with authentication server.
type Application struct {
	// GrantType - Authorization type, must be equal to `password`
	GrantType string `json:"grant_type"`

	// ClientID - ID of your application
	ClientID string `json:"client_id"`

	// ClientSecret - Secret key of your application
	ClientSecret string `json:"client_secret"`

	// Username - User username
	Username string `json:"username"`

	// Password - User password
	Password string `json:"password"`

	// Scope - Access rights required by the application
	Scope int64 `json:"scope"`
}

// Values returns values from this Application.
func (app *Application) Values() (values url.Values) {
	values = url.Values{}
	values.Set(paramGrantType, app.GrantType)
	values.Set(paramClientID, app.ClientID)
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
	app.ClientID = defaultClientID
	app.ClientSecret = defaultClientSecret

	return
}

// Authenticate authenticates *ApiClient through Application.
// If the outcome is successful, it returns a *AccessToken.
func Authenticate(api *APIClient, app Application) (token *AccessToken, err error) {
	token = new(AccessToken)
	if api.httpClient == nil {
		return nil, errors.New("HTTPClient not found.")
	}
	auth := OAuthURL()

	q := ConcatValues(false, auth.Query(), app.Values(), api.Values())
	//q.Set("test_redirect_uri", "1")
	if q != nil {
		auth.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(oAuthMethod, auth.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("Error: " + res.Status)
	}

	err = json.NewDecoder(res.Body).Decode(token)
	if err != nil {
		return nil, err
	}

	return
}
