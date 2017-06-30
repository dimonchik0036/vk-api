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

type LongPoll struct {
	Host      string `json:"server"`
	Path      string `json:"path"`
	Key       string `json:"key"`
	Ts        int64  `json:"ts"`
	LPVersion int    `json:"-"`
	NeedPts   int    `json:"-"`
}

type LPUpdate []interface{}

type LPUpdates struct {
	Failed  int64      `json:"failed"`
	Ts      int64      `json:"ts"`
	Updates []LPUpdate `json:"updates"`
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

func (client *Client) GetLPUpdates(config LPConfig) (LPUpdates, error) {
	if client.apiClient == nil {
		return LPUpdates{}, errors.New("A api client was not initialized")
	}

	if client.LongPoll == nil {
		return LPUpdates{}, errors.New("A long poll was not initialized")
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
		return LPUpdates{}, err
	}

	res, err := client.apiClient.httpClient.Do(req)
	if err != nil {
		return LPUpdates{}, err
	}

	if res.StatusCode != http.StatusOK {
		return LPUpdates{}, errors.New(res.Status)
	}

	var answer LPUpdates
	if err = json.NewDecoder(res.Body).Decode(&answer); err != nil {
		return LPUpdates{}, err
	}

	return answer, nil
}

func (client *Client) GetLPUpdatesChan(bufSize int, config LPConfig) (LPChan, error) {
	ch := make(chan LPUpdate, bufSize)

	go func() {
		for {
			update, err := client.GetLPUpdates(config)
			if err != nil {
				log.Println(err)
				log.Println("Failed to get updates, retrying in 3 seconds...")
				time.Sleep(time.Second * 3)

				continue
			}

			switch update.Failed {
			case 0:
				if len(update.Updates) == 0 {
					continue
				}

				for i := len(update.Updates) - 1; i >= 0; i-- {
					ch <- update.Updates[i]
				}

				client.LongPoll.Ts = update.Ts
			case 1:
				client.LongPoll.Ts = update.Ts
				log.Println("Ts updated")
			case 2, 3:
				if err := client.InitLongPoll(client.LongPoll.NeedPts, client.LongPoll.LPVersion); err != nil {
					log.Println("Long poll update error:", err)
					return
				}

				log.Println("Long poll config updated")
			}
		}
	}()

	return ch, nil
}
