package vkapi

import (
	"errors"
	"log"
	"net/url"
	"strconv"
	"strings"
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

// Destination describes the final destination.
type Destination struct {
	UserID      int64    `json:"user_id"`
	PeerID      int64    `json:"peer_id"`
	Domain      string   `json:"domain"`
	ChatID      int64    `json:"chat_id"`
	GroupID     int64    `json:"group_id"`
	UserIDs     []int64  `json:"user_ids"`
	ScreenName  string   `json:"user_id"`
	ScreenNames []string `json:"user_ids"`
}

func (dst Destination) Values() (values url.Values) {
	values = url.Values{}

	switch {
	case dst.UserID != 0:
		values.Add("user_id", strconv.FormatInt(dst.UserID, 10))
	case dst.PeerID != 0:
		values.Add("peer_id", strconv.FormatInt(dst.PeerID, 10))
	case dst.Domain != "":
		values.Add("domain", dst.Domain)
	case dst.ChatID != 0:
		values.Add("chat_id", strconv.FormatInt(dst.ChatID, 10))
	case dst.GroupID != 0:
		values.Add("group_id", strconv.FormatInt(dst.GroupID, 10))
	case len(dst.UserIDs) != 0:
		values.Add("user_ids", ConcatInt64ToString(dst.UserIDs...))
	case dst.ScreenName != "":
		values.Add("user_id", dst.ScreenName)
	case len(dst.ScreenNames) > 0:
		values.Add("user_ids", strings.Join(dst.ScreenNames, ","))
	}

	return
}

// NewDstFromUserID creates a new MessageConfig instance from userID.
func NewDstFromUserID(userIDs ...int64) (dst Destination) {
	if len(userIDs) == 1 {
		dst.UserID = userIDs[0]
	} else {
		dst.UserIDs = userIDs
	}
	return
}

// NewDstFromScreenName creates a new MessageConfig instance from userID.
func NewDstFromScreenName(screenNames ...string) (dst Destination) {
	if len(screenNames) == 1 {
		dst.ScreenName = screenNames[0]
	} else {
		dst.ScreenNames = screenNames
	}
	return
}

// NewDstFromPeerID creates a new MessageConfig instance from peerID.
func NewDstFromPeerID(peerID int64) (dst Destination) {
	dst.PeerID = peerID
	return
}

// NewDstFromChatID creates a new MessageConfig instance from chatID.
func NewDstFromChatID(chatID int64) (dst Destination) {
	dst.ChatID = chatID
	return
}

// NewDstFromGroupID creates a new MessageConfig instance from groupID.
func NewDstFromGroupID(groupID int64) (dst Destination) {
	dst.GroupID = groupID
	return
}

// NewDstFromDomain creates a new MessageConfig instance from domain.
func NewDstFromDomain(domain string) (dst Destination) {
	dst.Domain = domain
	return
}
