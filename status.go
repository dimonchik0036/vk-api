package vkapi

import (
	"encoding/json"
	"net/url"
	"strconv"
)

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

func (client *Client) SetMyStatus(text string) (err error) {
	return client.SetStatus(0, text)
}

func (client *Client) Status(user_id int64) (text string, err error) {
	req := Request{}
	req.Method = "status.get"
	req.Values = url.Values{}
	if user_id != 0 {
		req.Values.Add("user_id", strconv.FormatInt(user_id, 10))
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

	err = json.Unmarshal((*res).Response.Bytes(), &Text)
	if err != nil {
		return "", err
	}

	return Text.Text, nil
}

func (client *Client) MyStatus() (text string, err error) {
	return client.Status(0)
}
