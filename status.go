package vkapi

import (
	"net/url"
	"strconv"
)

// SetStatus set status in your group.
func (client *Client) SetStatus(groupId int64, text string) (err error) {
	req := Request{}
	req.Method = "status.set"
	req.Values = url.Values{}
	req.Values.Add("text", text)
	if groupId != 0 {
		req.Values.Add("group_id", strconv.FormatInt(groupId, 10))
	}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return
}

// SetMyStatus set status on your page.
func (client *Client) SetMyStatus(text string) (err error) {
	return client.SetStatus(0, text)
}

// Status returns the status from the user page.
func (client *Client) Status(userId int64) (text string, err error) {
	req := Request{}
	req.Method = "status.get"
	req.Values = url.Values{}
	if userId != 0 {
		req.Values.Add("user_id", strconv.FormatInt(userId, 10))
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if (*res).ServerError() != nil {
		return "", (*res).Error
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
	return client.Status(0)
}
