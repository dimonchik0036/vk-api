package vkapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	FlagMessageUnread = 1 << iota
	FlagMessageOutBox
	FlagMessageReplied
	FlagMessageImportant
	FlagMessageChat
	FlagMessageFriends
	FlagMessageSpam
	FlagMessageDeleted
	FlagMessageFixed
	FlagMessageMedia
	FlagMessageHidden = 65536

	LPModeAttachments   = 2
	LPModeExtendedEvent = 8
	LPModePts           = 32
	LPModeExtra         = 64
	LPModeRandomID      = 128

	LPCodeNewMessage = 4
)

type LongPoll struct {
	Host      string `json:"server"`
	Path      string `json:"path"`
	Key       string `json:"key"`
	Ts        int64  `json:"ts"`
	LPVersion int    `json:"-"`
	NeedPts   int    `json:"-"`
}

type LPUpdate struct {
	Code    int64
	Update  []interface{}
	Message *LPMessage
}

func (update *LPUpdate) UnmarshalUpdate(mode int) error {
	update.Code = int64(update.Update[0].(float64))

	switch update.Code {
	case LPCodeNewMessage:
		message := new(LPMessage)

		message.ID = int64(update.Update[1].(float64))
		message.Flags = int64(update.Update[2].(float64))
		message.FromID = int64(update.Update[3].(float64))
		message.Timestamp = int64(update.Update[4].(float64))
		message.Text = update.Update[5].(string)

		if mode&LPModeAttachments == LPModeAttachments {
			message.Attachments = make(map[string]string)
			for key, value := range update.Update[6].(map[string]interface{}) {
				message.Attachments[key] = value.(string)
			}
		}

		if mode&LPModeRandomID&LPModeRandomID == (LPModeAttachments | LPModeRandomID) {
			message.RandomId = int64(update.Update[7].(float64))
		} else {
			if mode&LPModeRandomID == LPModeRandomID {
				message.RandomId = int64(update.Update[6].(float64))
			}
		}

		update.Message = message
	}

	return nil
}

type LPMessage struct {
	ID          int64
	Flags       int64
	FromID      int64
	Timestamp   int64
	Text        string
	Attachments map[string]string
	RandomId    int64
}

type LPAnswer struct {
	Failed  int64           `json:"failed"`
	Ts      int64           `json:"ts"`
	Updates [][]interface{} `json:"updates"`
}

type LPChan <-chan LPUpdate

func (client *Client) InitLongPoll(needPts int, lpVersion int) *Error {
	var req Request
	req.Method = "messages.getLongPollServer"

	v := url.Values{}
	v.Add("need_pts", strconv.FormatInt(int64(needPts), 10))
	v.Add("lp_version", strconv.FormatInt(int64(lpVersion), 10))
	req.Values = v

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	client.LongPoll = new(LongPoll)
	if err := json.Unmarshal(res.Response.Bytes(), &client.LongPoll); err != nil {
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

type LPConfig struct {
	Wait int
	Mode int
}

func (client *Client) GetLPAnswer(config LPConfig) (LPAnswer, error) {
	if client.apiClient == nil {
		return LPAnswer{}, errors.New("A api client was not initialized")
	}

	if client.LongPoll == nil {
		return LPAnswer{}, errors.New("A long poll was not initialized")
	}

	values := url.Values{}
	values.Add("act", "a_check")
	values.Add("key", client.LongPoll.Key)
	values.Add("ts", strconv.FormatInt(client.LongPoll.Ts, 10))
	values.Add("wait", strconv.FormatInt(int64(config.Wait), 10))
	values.Add("mode", strconv.FormatInt(int64(config.Mode), 10))
	values.Add("version", strconv.FormatInt(int64(client.LongPoll.LPVersion), 10))

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
		return LPAnswer{}, err
	}

	if res.StatusCode != http.StatusOK {
		return LPAnswer{}, errors.New(res.Status)
	}

	var answer LPAnswer
	if err = json.NewDecoder(res.Body).Decode(&answer); err != nil {
		return LPAnswer{}, err
	}

	return answer, nil
}

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
			LPUpdate.UnmarshalUpdate(config.Mode)
			LPUpdates = append(LPUpdates, LPUpdate)
		}

		client.LongPoll.Ts = answer.Ts
		return LPUpdates, nil
	case 1:
		client.LongPoll.Ts = answer.Ts
		if client.apiClient.Log {
			client.apiClient.Logger.Println("Ts updated")
		}

	case 2, 3:
		if err := client.InitLongPoll(client.LongPoll.NeedPts, client.LongPoll.LPVersion); err != nil {
			if client.apiClient.Log {
				client.apiClient.Logger.Println("Long poll update error:", err)
			}
			return []LPUpdate{}, err
		}

		if client.apiClient.Log {
			client.apiClient.Logger.Println("Long poll config updated")
		}
	}

	return []LPUpdate{}, nil
}

func (client *Client) GetLPUpdatesChan(bufSize int, config LPConfig) (LPChan, *bool, error) {
	ch := make(chan LPUpdate, bufSize)
	run := true

	go func() {
		for run {
			updates, err := client.GetLPUpdates(config)
			if err != nil {
				log.Println("Failed to get updates, retrying in 3 seconds...")
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
