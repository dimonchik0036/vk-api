package vkapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

type Keyboard struct {
	OneTime bool       `json:"one_time"`
	Buttons [][]Button `json:"buttons"`
}

type Button struct {
	Action Action      `json:"action"`
	Color  ColorButton `json:"color"`
}

type Action struct {
	Type  string `json:"type"`
	Label string `json:"label"`
}

type ColorButton string

const (
	Negative ColorButton = "negative"
	Positive ColorButton = "positive"
	Default  ColorButton = "default"
	Primary  ColorButton = "primary"
)

func CreateButtons(buttons ...[]Button) [][]Button {
	length := len(buttons)
	if length > 8 {
		panic("Amount rows are more 8.")
	}

	return buttons
}

func CreateButtonsInRow(textAndColorButtons map[string]ColorButton) []Button {
	length := len(textAndColorButtons)
	if length > 4 {
		panic("Amount buttons in the row are more 4.")
	}

	var buttons []Button

	for text, color := range textAndColorButtons {
		button := new(Button)
		button.Action.Type = "text"
		button.Action.Label = text
		button.Color = color

		buttons = append(buttons, *button)
	}

	return buttons
}

func (client *Client) SendKeyboard(id int64, message string, isOneTime bool, buttons [][]Button) (int64, *Error) {
	config := NewMessage(NewDstFromUserID(id), message)
	values := addParameters(config)

	keyboard := new(Keyboard)
	keyboard.OneTime = isOneTime
	keyboard.Buttons = buttons
	values.Add("keyboard", convertToJson(keyboard))

	request := NewRequest("messages.send", config.AccessToken, values)
	res, err := client.Do(request)
	if err != nil {
		return 0, err
	}

	answer, error := strconv.ParseInt(res.Response.String(), 10, 64)
	if error != nil {
		return 0, NewError(ErrBadResponseCode, error.Error())
	}

	return answer, nil
}

func convertToJson(i interface{}) string {
	jsonObject, _ := json.Marshal(i)
	stringJsonObject := string(jsonObject)
	return stringJsonObject
}

func addParameters(config MessageConfig) url.Values {
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

	if config.Attachment != "" {
		values.Add("attachment", config.Attachment)
	}

	return values
}
