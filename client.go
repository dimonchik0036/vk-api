package vkapi

import (
	"errors"
	"log"
)

// Client allows you to transparently send requests to API server.
type Client struct {
	apiClient *APIClient
	User      Users
	LongPoll  *LongPoll
	//	Group
	//	Wall
	//	WallComment
	//  Message
	//  Chat
	//	Note
	//	Page
	//	Board
	//	BoardComment
}

// SetLanguage sets the language in which different data will be returned,
// for example, names of countries and cities.
func (client *Client) SetLanguage(lang string) error {
	if client.apiClient == nil {
		return errors.New(ErrApiClientNotFound)
	}

	client.apiClient.Language = lang
	return nil
}

// SetLogger sets logger.
func (client *Client) SetLogger(logger *log.Logger) error {
	if client.apiClient == nil {
		return errors.New(ErrApiClientNotFound)
	}

	client.apiClient.Logger = logger
	return nil
}

// Log allow write log.
func (client *Client) Log(flag bool) error {
	if client.apiClient == nil {
		return errors.New(ErrApiClientNotFound)
	}

	client.apiClient.Log = flag
	return nil
}

// NewClientFromToken creates a new *Client instance.
func NewClientFromToken(token string) (client *Client, err error) {
	client = new(Client)
	client.apiClient = NewApiClient()
	client.apiClient.SetAccessToken(token)
	return
}

// NewClientFromLogin creates a new *Client instance
// and allows you to pass a authentication.
func NewClientFromLogin(username string, password string, scope int64) (client *Client, err error) {
	client = new(Client)
	client.apiClient = NewApiClient()
	err = client.apiClient.Authenticate(NewApplication(username, password, scope))
	if err != nil {
		return nil, err
	}

	return
}

// Do makes a request to a specific endpoint with our request
// and returns response.
func (client *Client) Do(request Request) (response *Response, err *Error) {
	if client.apiClient == nil {
		return nil, NewError(ErrBadCode, ErrApiClientNotFound)
	}

	if request.Token == "" && client.apiClient.AccessToken != nil {
		request.Token = client.apiClient.AccessToken.AccessToken
	}

	return client.apiClient.Do(request)
}
