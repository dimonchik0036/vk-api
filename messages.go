package vkapi

import (
	"fmt"
	"net/url"
	"strconv"
)

const (
	chatOffset = 2000000000
)

// Dialog describes the structure of the message.
type Dialog struct {
	Unread     int64    `json:"unread"`
	Message    *Message `json:"message"`
	InRead     int64    `json:"in_read"`
	OutRead    int64    `json:"out_read"`
	RealOffset int64    `json:"real_offset"`
}

// Message describes the structure of the message.
type Message struct {
	Id          int64      `json:"id"`
	UserId      int64      `json:"user_id"`
	FromId      int64      `json:"from_id"`
	Date        int64      `json:"date"`
	ReadState   int        `json:"read_state"`
	Out         int        `json:"out"`
	Title       string     `json:"title"`
	Body        string     `json:"body"`
	FwdMessages *[]Message `json:"fwd_messages"`
	Emoji       int        `json:"emoji"`
	Important   int        `json:"important"`
	Deleted     int        `json:"deleted"`
	RandomId    int64      `json:"random_id"`
	ChatId      int64      `json:"chat_id"`
	ChatActive  []int64    `json:"chat_active"`
	UsersCount  int        `json:"users_count"`
	AdminId     int64      `json:"admin_id"`
	Action      string     `json:"action"`
	ActionMid   int64      `json:"action_mid"`   /*идентификатор пользователя (если > 0) или email (если < 0), которого пригласили или исключили (для служебных сообщений с action = chat_invite_user или chat_kick_user). */
	ActionEmail string     `json:"action_email"` /*email, который пригласили или исключили (для служебных сообщений с action = chat_invite_user или chat_kick_user и отрицательным action_mid). */
	ActionText  string     `json:"action_text"`  /*название беседы (для служебных сообщений с action = chat_create или chat_title_update). */
	Photo50     string     `json:"photo_50"`
	Photo100    string     `json:"photo_100"`
	Photo200    string     `json:"photo_200"`
	/*Geo       *Geo {
		type (string) — тип места;
		coordinates (string) — координаты места;
		place (object) — описание места (если оно добавлено), объект с полями:
		id (integer) — идентификатор места (если назначено);
		title (string) — название места (если назначено);
		latitude (number) — географическая широта;
		longitude (number) — географическая долгота;
		created (integer) — дата создания (если назначено);
		icon (string) — URL изображения-иконки;
		country (string) — название страны;
		city (string) — название города;
	} `json:"geo"`*/

	/*Attachments *[]Attachments `json:"attachments"`*/
	/*PushSettings *PushSettings { настройки уведомлений для беседы, если они есть.	} `json:"push_settings"`*/
	/*string	тип действия (если это служебное сообщение). Возможные значения:

	  chat_photo_update — обновлена фотография беседы;
	  chat_photo_remove — удалена фотография беседы;
	  chat_create — создана беседа;
	  chat_title_update — обновлено название беседы;
	  chat_invite_user — приглашен пользователь;
	  chat_kick_user — исключен пользователь.*/

}

// MessageConfig contains the data
// necessary to send a message.
type MessageConfig struct {
	UserID          int64   `json:"user_id"`
	RandomID        int64   `json:"random_id"`
	PeerID          int64   `json:"peer_id"`
	Domain          string  `json:"domain"`
	ChatID          int64   `json:"chat_id"`
	GroupID         int64   `json:"group_id"`
	UserIDs         []int64 `json:"user_ids"`
	Message         string  `json:"message"`
	geo             bool    `json:"-"`
	lat             float64 `json:"lat"`
	long            float64 `json:"long"`
	ForwardMessages []int64 `json:"forward_messages"`
	StickerID       int64   `json:"sticker_id"`
	AccessToken     string  `json:"access_token"`
	//attachment *[]Attachment `json:"attachment"`
}

// SetGeo sets the location.
func (m *MessageConfig) SetGeo(lat float64, long float64) {
	m.geo = true
	m.lat = lat
	m.long = long
}

// NewMessage creates a new message for the user from the text.
func NewMessage(id int64, message string) (config MessageConfig) {
	config.PeerID = id
	config.Message = message
	return
}

// NewMessageToChat creates a new message for the chat from the text.
func NewMessageToChat(id int64, message string) (config MessageConfig) {
	return NewMessage(id+chatOffset, message)
}

// NewMessageToUsers creates a new message for several users from the text.
func NewMessageToUsers(message string, ids ...int64) (config MessageConfig) {
	config.UserIDs = ids
	config.Message = message
	return
}

// SendMessage tries to send a message with the configuration
// from the MessageConfig and returns message ID if it succeeds.
func (client *Client) SendMessage(config MessageConfig) (int64, *Error) {
	var req Request
	req.Token = config.AccessToken
	req.Method = "messages.send"
	v := url.Values{}

	if config.PeerID != 0 {
		v.Add("peer_id", fmt.Sprintf("%d", config.PeerID))
	}

	if config.UserID != 0 {
		v.Add("user_id", fmt.Sprintf("%d", config.UserID))
	}

	if config.Domain != "" {
		v.Add("domain", config.Domain)
	}

	if config.ChatID != 0 {
		v.Add("chat_id", fmt.Sprintf("%d", config.RandomID))
	}

	if len(config.UserIDs) != 0 {
		v.Add("user_ids", ConcatInt64ToString(config.UserIDs...))
	}

	if len(config.ForwardMessages) != 0 {
		v.Add("forward_messages", ConcatInt64ToString(config.ForwardMessages...))
	}

	if config.StickerID != 0 {
		v.Add("sticker_id", fmt.Sprintf("%d", config.StickerID))
	}

	if config.Message != "" {
		v.Add("message", config.Message)
	}

	if config.RandomID != 0 {
		v.Add("random_id", fmt.Sprintf("%d", config.RandomID))
	}

	if config.geo {
		v.Add("lat", strconv.FormatFloat(config.lat, 'f', -1, 64))
		v.Add("long", strconv.FormatFloat(config.long, 'f', -1, 64))
	}

	req.Values = v
	res, err := client.Do(req)
	if err != nil && !err.Code.Is(ErrZero) {
		return 0, err
	}

	answer, error := strconv.ParseInt(res.Response.String(), 10, 64)
	if error != nil {
		return 0, NewError(ErrBadResponseCode, error.Error())
	}

	return answer, nil
}

// SetActivity changes the status of typing by user in the dialog.
// Accepts userID as string or int64, chat as 2000000000 + chatID,
// group as -groupID.
func (client *Client) SetActivity(dst interface{}) *Error {
	values := url.Values{}
	switch dst.(type) {
	case string:
		values.Add("user_id", dst.(string))
	case int64:
		values.Add("peer_id", strconv.FormatInt(dst.(int64), 10))
	case int:
		values.Add("peer_id", strconv.FormatInt(int64(dst.(int)), 10))
	default:
		return NewError(ErrBadCode, "Wrong data")
	}

	values.Add("type", "typing")
	_, err := client.Do(NewRequest("messages.setActivity", "", values))
	if err != nil {
		return err
	}

	return nil
}
