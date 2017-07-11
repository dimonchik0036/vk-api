package vkapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

const (
	LPMessageFlagUnread = 1 << iota
	LPMessageFlagOutBox
	LPMessageFlagReplied
	LPMessageFlagImportant
	LPMessageFlagChat
	LPMessageFlagFriends
	LPMessageFlagSpam
	LPMessageFlagDeleted
	LPMessageFlagFixed
	LPMessageFlagMedia
	LPMessageFlagHidden = 65536
)

const (
	LPModeAttachments   = 2
	LPModeExtendedEvent = 8
	LPModePts           = 32
	LPModeExtra         = 64
	LPModeRandomID      = 128
)

const (
	LPCodeMessageSetFlags      = 1
	LPCodeMessageAddFlags      = 2
	LPCodeMessageDelFlags      = 3
	LPCodeNewMessage           = 4
	LPCodeReadAllInboxMessage  = 6
	LPCodeReadAllOutboxMessage = 7
	LPCodeFriendOnline         = 8
	LPCodeFriendOffline        = 9
	LPCodeDialogDelFlags       = 10
	LPCodeDialogSetFlags       = 11
	LPCodeDialogAddFlags       = 12
	LPCodeDelAllMessage        = 13
	LPCodeChangeChat           = 51
	LPCodeTypingInDialog       = 61
	LPCodeTypingInChat         = 62
	LPCodeCall                 = 70
	LPCodeUnreadMessage        = 80
	LPCodeChangeNotification   = 114
)

const (
	LPPlatformUndefined = iota
	LPPlatformMobile
	LPPlatformIPhone
	LPPlatformIPad
	LPPlatformAndroid
	LPPlatformWPhone
	LPPlatformWindows
	LPPlatformWeb
)

// Timestamp is the wrapper of int64.
type Timestamp int64

func (ts Timestamp) String() string {
	return time.Unix(int64(ts), 0).Format("15:04:05 02/01/2006")
}

// LongPoll allow you to interact with long poll server.
type LongPoll struct {
	Host      string    `json:"server"`
	Path      string    `json:"path"`
	Key       string    `json:"key"`
	Timestamp Timestamp `json:"ts"`
	LPVersion int       `json:"-"`
	NeedPts   int       `json:"-"`
}

// LPUpdate stores response from a long poll server.
type LPUpdate struct {
	Code               int64
	Update             []interface{}
	Message            *LPMessage
	FriendNotification *LPFriendNotification
}

// Event returns event as a string.
func (update *LPUpdate) Event() (event string) {
	switch update.Code {
	case LPCodeMessageSetFlags:
		event = "Setting the message flags"
	case LPCodeMessageAddFlags:
		event = "Adding the message flags"
	case LPCodeMessageDelFlags:
		event = "Deleting the message flags"
	case LPCodeNewMessage:
		event = "New message"
	case LPCodeReadAllInboxMessage:
		event = "You read inbox message"
	case LPCodeReadAllOutboxMessage:
		event = "You read outbox message"
	case LPCodeFriendOnline:
		event = "Friend online"
	case LPCodeFriendOffline:
		event = "Friend offline"
	case LPCodeDialogDelFlags:
		event = "Deleting the dialog flags"
	case LPCodeDialogSetFlags:
		event = "Setting the dialog flags"
	case LPCodeDialogAddFlags:
		event = "Adding the dialog flags"
	case LPCodeDelAllMessage:
		event = "Deleting all messages in the dialog"
	case LPCodeChangeChat:
		event = "Chat parameters changed"
	case LPCodeTypingInDialog:
		event = "User typing text in the dialog"
	case LPCodeTypingInChat:
		event = "User typing text in the chat"
	case LPCodeCall:
		event = "User called"
	case LPCodeUnreadMessage:
		event = "Number of unread messages"
	case LPCodeChangeNotification:
		event = "Notifications setting changed"
	default:
		event = fmt.Sprintf("Undefined event (%d)", update.Code)
	}

	return
}

// IsMessageSetFlags will return true if the message flags have been replaced.
func (update *LPUpdate) IsMessageSetFlags() bool {
	return update.Code == LPCodeMessageSetFlags
}

// IsMessageAddFlags will return true if the message flags have been added.
func (update *LPUpdate) IsMessageAddFlags() bool {
	return update.Code == LPCodeMessageAddFlags
}

// IsMessageDelFlags will return true if the message flags have been deleted.
func (update *LPUpdate) IsMessageDelFlags() bool {
	return update.Code == LPCodeMessageDelFlags
}

// IsNewMessage will return true if it is a new message.
func (update *LPUpdate) IsNewMessage() bool {
	return update.Code == LPCodeNewMessage
}

// IsFriendOnline will return true if a friend became online.
func (update *LPUpdate) IsFriendOnline() bool {
	return update.Code == LPCodeFriendOnline
}

// IsFriendOffline will return true if a friend became offline.
func (update *LPUpdate) IsFriendOffline() bool {
	return update.Code == LPCodeFriendOffline
}

// IsDialogDelFlags will return true if the dialog flags have been deleted.
func (update *LPUpdate) IsDialogDelFlags() bool {
	return update.Code == LPCodeDialogDelFlags
}

// IsDialogSetFlags will return true if the dialog flags have been replaced.
func (update *LPUpdate) IsDialogSetFlags() bool {
	return update.Code == LPCodeDialogSetFlags
}

// IsDialogAddFlags will return true if the dialog flags have been added.
func (update *LPUpdate) IsDialogAddFlags() bool {
	return update.Code == LPCodeDialogAddFlags
}

func unescaped(string string) string {
	string = html.UnescapeString(string)
	reg, err := regexp.Compile("<br>")
	if err != nil {
		return string
	}

	return reg.ReplaceAllLiteralString(string, "\n")
}

// UnmarshalUpdate unmarshal a LPUpdate.
func (update *LPUpdate) UnmarshalUpdate(mode int) error {
	update.Code = int64(update.Update[0].(float64))
	updateLen := len(update.Update)
	switch update.Code {
	case LPCodeMessageSetFlags, LPCodeMessageAddFlags, LPCodeMessageDelFlags, LPCodeNewMessage:
		message := new(LPMessage)
		message.Type = update.Code
		message.ID = int64(update.Update[1].(float64))
		message.Flags = int64(update.Update[2].(float64))
		if updateLen == 3 {
			update.Message = message
			break
		}

		message.FromID = int64(update.Update[3].(float64))
		if updateLen == 4 {
			update.Message = message
			break
		}

		message.Timestamp = Timestamp(update.Update[4].(float64))
		message.Text = unescaped(update.Update[5].(string))

		if updateLen == 6 {
			update.Message = message
			break
		}

		if mode&LPModeAttachments != 0 {
			message.Attachments = make(map[string]string)
			for key, value := range update.Update[6].(map[string]interface{}) {
				message.Attachments[key] = value.(string)
			}
		} else {
			if mode&LPModeRandomID != 0 {
				message.RandomID = int64(update.Update[6].(float64))
			}
		}

		if updateLen == 7 {
			update.Message = message
			break
		}

		if mode&LPModeRandomID&LPModeAttachments != 0 {
			message.RandomID = int64(update.Update[7].(float64))
		}

		update.Message = message
	case LPCodeDialogDelFlags, LPCodeDialogSetFlags, LPCodeDialogAddFlags:
		message := new(LPMessage)
		message.Type = update.Code
		message.FromID = int64(update.Update[1].(float64))
		message.Flags = int64(update.Update[2].(float64))

		update.Message = message
	case LPCodeFriendOnline, LPCodeFriendOffline:
		if len(update.Update) < 3 {
			return errors.New("(" + string(update.Code) + ") invalid update size.")
		}

		friend := new(LPFriendNotification)
		friend.Activity = update.Code
		friend.ID = -int64(update.Update[1].(float64))
		friend.Arg = int64(update.Update[2].(float64)) & 0xFF
		friend.Timestamp = Timestamp(update.Update[3].(float64))

		update.FriendNotification = friend
	case LPCodeReadAllInboxMessage, LPCodeReadAllOutboxMessage, LPCodeDelAllMessage:
		message := new(LPMessage)
		message.Type = update.Code
		message.FromID = int64(update.Update[1].(float64))
		message.ID = int64(update.Update[2].(float64))

		update.Message = message
	case LPCodeTypingInDialog:
		message := new(LPMessage)
		message.Type = update.Code
		message.FromID = int64(update.Update[1].(float64))

		update.Message = message
	}

	return nil
}

// LPMessage is new messages
// that come from long poll server.
type LPMessage struct {
	Type        int64
	ID          int64
	Flags       int64
	FromID      int64
	Timestamp   Timestamp
	Text        string
	Attachments map[string]string
	RandomID    int64
}

func (message *LPMessage) String() string {
	return fmt.Sprintf("Message (%d):`%s` from (%d) at %s", message.ID, message.Text, message.FromID, message.Timestamp)
}

func (message *LPMessage) LastMessage() int64 {
	return message.ID
}

// Unread will return true if the message is not read.
func (message *LPMessage) Unread() bool {
	return message.Flags&LPMessageFlagUnread != 0
}

// Outbox will return true if this is an outgoing message.
func (message *LPMessage) Outbox() bool {
	return message.Flags&LPMessageFlagOutBox != 0
}

// Replied will be returned true if an answer was created to the message.
func (message *LPMessage) Replied() bool {
	return message.Flags&LPMessageFlagReplied != 0
}

// Important will return true if this is a marked message.
func (message *LPMessage) Important() bool {
	return message.Flags&LPMessageFlagImportant != 0
}

// FromChat will return true if this message was sent via chat.
func (message *LPMessage) FromChat() bool {
	return message.Flags&LPMessageFlagChat != 0
}

// FromFriends will return true if this message was sent from friends.
// Not applicable for messages from group conversations.
func (message *LPMessage) FromFriends() bool {
	return message.Flags&LPMessageFlagFriends != 0
}

// IsSpam will return true if it is spam.
func (message *LPMessage) IsSpam() bool {
	return message.Flags&LPMessageFlagSpam != 0
}

// Deleted will return true if the message was deleted (in the Recycle Bin).
func (message *LPMessage) Deleted() bool {
	return message.Flags&LPMessageFlagDeleted != 0
}

// Fixed will return true if the message has been scanned by the user for spam.
func (message *LPMessage) Fixed() bool {
	return message.Flags&LPMessageFlagFixed != 0
}

// ContainsMedia will return true if the message contains multimedia content.
func (message *LPMessage) ContainsMedia() bool {
	return message.Flags&LPMessageFlagMedia != 0
}

// IsHidden will return true if it is a welcome message from the community.
func (message *LPMessage) IsHidden() bool {
	return message.Flags&LPMessageFlagHidden != 0
}

// LPFriendNotification is a notification
// that a friend has become online or offline.
type LPFriendNotification struct {
	ID int64

	// If friend is online,
	// then Arg is equal to platform.
	//
	// If the friend offline, then
	// 0 - friend logout,
	// 1 - offline by timeout.
	Arg       int64
	Timestamp Timestamp
	Activity  int64
}

// Status returns event as a string.
func (friend *LPFriendNotification) Status() (status string) {
	switch friend.Activity {
	case LPCodeFriendOnline:
		status = "Online"
	case LPCodeFriendOffline:
		status = "Offline"
	default:
		status = "Undefined event"
	}

	return
}

func (friend *LPFriendNotification) String() string {
	return fmt.Sprintf("Friend (%d) was %s at %s", friend.ID, friend.Status(), friend.Timestamp)
}

// Platform returns the name of the platform.
func (friend *LPFriendNotification) Platform() (platform string) {
	switch friend.Arg % 0xFF {
	case LPPlatformMobile:
		platform = "Mobile"
	case LPPlatformIPhone:
		platform = "IPhone"
	case LPPlatformIPad:
		platform = "IPad"
	case LPPlatformAndroid:
		platform = "Android"
	case LPPlatformWPhone:
		platform = "Windows Phone"
	case LPPlatformWindows:
		platform = "Windows"
	case LPPlatformWeb:
		platform = "Web"
	default:
		platform = "Undefined platform"
	}

	return
}

// LPAnswer is response from long poll server.
type LPAnswer struct {
	Failed    int64           `json:"failed"`
	Timestamp Timestamp       `json:"ts"`
	Updates   [][]interface{} `json:"updates"`
}

// LPChan allows to receive new LPUpdate.
type LPChan <-chan LPUpdate

// InitLongPoll establishes a new connection
// to long poll server.
func (client *Client) InitLongPoll(needPts int, lpVersion int) *Error {
	values := url.Values{}
	values.Add("need_pts", strconv.FormatInt(int64(needPts), 10))
	values.Add("lp_version", strconv.FormatInt(int64(lpVersion), 10))

	res, err := client.Do(NewRequest("messages.getLongPollServer", "", values))
	if err != nil {
		return err
	}

	client.LongPoll = new(LongPoll)
	if err := res.To(&client.LongPoll); err != nil {
		return NewError(ErrBadCode, err.Error())
	}

	u, error := url.Parse(client.LongPoll.Host)
	if error != nil {
		return NewError(ErrBadCode, error.Error())
	}

	client.LongPoll.Host = u.Host
	client.LongPoll.Path = u.Path
	client.LongPoll.LPVersion = lpVersion
	client.LongPoll.NeedPts = needPts

	return nil
}

// LPConfig stores data to connect to long poll server.
type LPConfig struct {
	Wait int
	Mode int
}

// GetLPAnswer makes a query with parameters
// from LPConfig to long poll server
// and returns a LPAnswer in case of success.
func (client *Client) GetLPAnswer(config LPConfig) (LPAnswer, error) {
	if client.apiClient == nil {
		return LPAnswer{}, errors.New(ErrApiClientNotFound)
	}

	if client.LongPoll == nil {
		return LPAnswer{}, errors.New("A long poll was not initialized")
	}

	values := url.Values{}
	values.Add("act", "a_check")
	values.Add("key", client.LongPoll.Key)
	values.Add("ts", strconv.FormatInt(int64(client.LongPoll.Timestamp), 10))
	values.Add("wait", strconv.FormatInt(int64(config.Wait), 10))
	values.Add("mode", strconv.FormatInt(int64(config.Mode), 10))
	values.Add("version", strconv.FormatInt(int64(client.LongPoll.LPVersion), 10))

	if client.apiClient.Log {
		client.apiClient.Logger.Printf("Request: %s", NewRequest("getLongPoll", "", values).JS())
	}

	u := url.URL{}
	u.Host = client.LongPoll.Host
	u.Path = client.LongPoll.Path
	u.Scheme = "https"
	u.RawQuery = values.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return LPAnswer{}, err
	}

	res, err := client.apiClient.httpClient.Do(req)
	if err != nil {
		client.apiClient.logPrintf("Response error: %s", err.Error())
		return LPAnswer{}, err
	}

	var reader io.Reader
	reader = res.Body

	if client.apiClient.Log {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		client.apiClient.Logger.Printf("Response: %s", string(b))
		reader = bytes.NewReader(b)
	}

	if res.StatusCode != http.StatusOK {
		client.apiClient.Logger.Printf("Response error: %s", res.Status)
		return LPAnswer{}, errors.New(res.Status)
	}

	var answer LPAnswer
	if err = json.NewDecoder(reader).Decode(&answer); err != nil {
		return LPAnswer{}, err
	}

	return answer, nil
}

// GetLPUpdates makes a query with parameters
// from LPConfig to long poll server
// and returns array LPUpdate in case of success.
func (client *Client) GetLPUpdates(config LPConfig) ([]LPUpdate, error) {
	answer, err := client.GetLPAnswer(config)
	if err != nil {
		return []LPUpdate{}, err
	}

	var LPUpdates []LPUpdate

	switch answer.Failed {
	case 0:
		for i := len(answer.Updates) - 1; i >= 0; i-- {
			var LPUpdate LPUpdate
			LPUpdate.Update = answer.Updates[i]
			if err := LPUpdate.UnmarshalUpdate(config.Mode); err != nil {
				client.apiClient.logPrintf("%s", err.Error())
			}

			LPUpdates = append(LPUpdates, LPUpdate)
		}

		client.LongPoll.Timestamp = answer.Timestamp
		return LPUpdates, nil
	case 1:
		client.LongPoll.Timestamp = answer.Timestamp
		client.apiClient.logPrintf("Timestamp updated")
	case 2, 3:
		if err := client.InitLongPoll(client.LongPoll.NeedPts, client.LongPoll.LPVersion); err != nil {
			client.apiClient.logPrintf("Long poll update error: %s", err.Error())

			return []LPUpdate{}, err
		}

		client.apiClient.logPrintf("Long poll config updated")
	}

	return []LPUpdate{}, nil
}

// GetLPUpdatesChan makes a query with parameters
// from LPConfig to long poll server
// and returns LPChan in case of success.
func (client *Client) GetLPUpdatesChan(bufSize int, config LPConfig) (LPChan, *bool, error) {
	ch := make(chan LPUpdate, bufSize)
	run := true

	go func() {
		for run {
			updates, err := client.GetLPUpdates(config)
			if err != nil {
				log.Print("Failed to get updates, retrying in 3 seconds...")
				time.Sleep(time.Second * 3)

				continue
			}

			for _, u := range updates {
				ch <- u
			}
		}

		close(ch)
	}()

	return ch, &run, nil
}
