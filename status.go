package vkapi

import (
	"net/url"
	"strconv"
)

// SetStatus set status in your group.
func (client *Client) SetStatus(groupID int64, text string) (err error) {
	values := url.Values{}
	values.Add("text", text)

	if groupID != 0 {
		values.Add("group_id", strconv.FormatInt(groupID, 10))
	}
	_, err = client.Do(NewRequest("status.set", "", values))
	if err != nil {
		return err
	}

	return
}

// SetMyStatus set status on your page.
func (client *Client) SetMyStatus(text string) (err error) {
	return client.SetStatus(0, text)
}

// Status returns the status from user or group page.
func (client *Client) Status(dst Destination) (text string, err error) {
	res, error := client.Do(NewRequest("status.get", "", dst.Values()))
	if err != nil {
		return "", error
	}

	Text := struct {
		Text string `json:"text"`
	}{}

	if err := res.To(&Text); err != nil {
		return "", err
	}

	return Text.Text, nil
}

// MyStatus returns the status from the Client page.
func (client *Client) MyStatus() (text string, err error) {
	return client.Status(Destination{})
}
