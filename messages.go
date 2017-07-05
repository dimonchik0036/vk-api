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
	ID          int64          `json:"id"`
	UserID      int64          `json:"user_id"`
	FromID      int64          `json:"from_id"`
	Date        Timestamp      `json:"date"`
	ReadState   int            `json:"read_state"`
	Out         int            `json:"out"`
	Title       string         `json:"title"`
	Body        string         `json:"body"`
	FwdMessages *[]Message     `json:"fwd_messages"`
	Emoji       int            `json:"emoji"`
	Important   int            `json:"important"`
	Deleted     int            `json:"deleted"`
	RandomID    int64          `json:"random_id"`
	ChatID      int64          `json:"chat_id"`
	ChatActive  []int64        `json:"chat_active"`
	UsersCount  int            `json:"users_count"`
	AdminID     int64          `json:"admin_id"`
	Action      string         `json:"action"`
	ActionMid   int64          `json:"action_mid"`
	ActionEmail string         `json:"action_email"`
	ActionText  string         `json:"action_text"`
	Photo50     string         `json:"photo_50"`
	Photo100    string         `json:"photo_100"`
	Photo200    string         `json:"photo_200"`
	Attachments *[]interface{} `json:"attachments"`
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

	/*PushSettings *PushSettings { настройки уведомлений для беседы, если они есть.	} `json:"push_settings"`*/
	/*string	тип действия (если это служебное сообщение). Возможные значения:

	  chat_photo_update — обновлена фотография беседы;
	  chat_photo_remove — удалена фотография беседы;
	  chat_create — создана беседа;
	  chat_title_update — обновлено название беседы;
	  chat_invite_user — приглашен пользователь;
	  chat_kick_user — исключен пользователь.*/

}

// IsDeleted will return true if the message was deleted (in the Recycle Bin).
func (message *Message) IsDeleted() bool {
	return message.Deleted != 0
}

// IsOutbox will return true if this is an outgoing message.
func (message *Message) IsOutbox() bool {
	return message.Out != 0
}

func (message *Message) String() string {
	if !message.IsDeleted() {
		return fmt.Sprintf("Message (%d):`%s` from (%d) at %s", message.ID, message.Body, message.FromID, message.Date)
	} else {
		return fmt.Sprintf("Message (%d) was deleted.", message.ID)
	}
}

// MessageConfig contains the data
// necessary to send a message.
type MessageConfig struct {
	Destination     Destination    `json:"destination"`
	RandomID        int64          `json:"random_id"`
	Message         string         `json:"message"`
	ForwardMessages []int64        `json:"forward_messages"`
	StickerID       int64          `json:"sticker_id"`
	AccessToken     string         `json:"access_token"`
	Attachment      *[]interface{} `json:"attachment"`
	geo             bool           `json:"-"`
	lat             float64        `json:"lat"`
	long            float64        `json:"long"`
}

// SetGeo sets the location.
func (m *MessageConfig) SetGeo(lat float64, long float64) {
	m.geo = true
	m.lat = lat
	m.long = long
}

// NewMessage creates a new message for the user from the text.
func NewMessage(dst Destination, message string) (config MessageConfig) {
	config.Destination = dst
	config.Message = message
	return
}

// SendMessage tries to send a message with the configuration
// from the MessageConfig and returns message ID if it succeeds.
func (client *Client) SendMessage(config MessageConfig) (int64, *Error) {
	values := config.Destination.Values()

	if len(config.ForwardMessages) != 0 {
		values.Add("forward_messages", ConcatInt64ToString(config.ForwardMessages...))
	}

	if config.StickerID != 0 {
		values.Add("sticker_id", fmt.Sprintf("%d", config.StickerID))
	}

	if config.Message != "" {
		values.Add("message", config.Message)
	}

	if config.RandomID != 0 {
		values.Add("random_id", fmt.Sprintf("%d", config.RandomID))
	}

	if config.geo {
		values.Add("lat", strconv.FormatFloat(config.lat, 'f', -1, 64))
		values.Add("long", strconv.FormatFloat(config.long, 'f', -1, 64))
	}

	res, err := client.Do(NewRequest("messages.send", config.AccessToken, values))
	if err != nil {
		return 0, err
	}

	answer, error := strconv.ParseInt(res.Response.String(), 10, 64)
	if error != nil {
		return 0, NewError(ErrBadResponseCode, error.Error())
	}

	return answer, nil
}

// SetActivity changes the status of typing by user in the dialog.
func (client *Client) SetActivity(dst Destination) *Error {
	values := url.Values{}

	if dst.GroupID != 0 {
		values.Add("peer_id", strconv.FormatInt(-dst.GroupID, 10))
	} else {
		values = dst.Values()
	}

	values.Add("type", "typing")
	_, err := client.Do(NewRequest("messages.setActivity", "", values))
	if err != nil {
		return err
	}

	return nil
}

// GetMessagesByID returns messages by ID.
func (client *Client) GetMessagesByID(previewLength int64, ids ...int64) ([]Message, *Error) {
	values := url.Values{}
	if previewLength != 0 {
		values.Add("preview_length", strconv.FormatInt(previewLength, 10))
	}

	values.Add("message_ids", ConcatInt64ToString(ids...))

	res, err := client.Do(NewRequest("messages.getById", "", values))
	if err != nil {
		return []Message{}, err
	}

	answer := struct {
		Items []Message `json:"items"`
	}{}

	if err := res.To(&answer); err != nil {
		return []Message{}, NewError(ErrBadCode, err.Error())
	}

	return answer.Items, nil
}
