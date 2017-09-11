package vkapi

import (
	"net/url"
	"strconv"
	"strings"
)

const (
	PostTypePost     = "post"
	PostTypeCopy     = "copy"
	PostTypeReply    = "reply"
	PostTypePostpone = "postpone"
	PostTypeSuggest  = "suggest"
)

type Wall struct {
	ID           int64         `json:"id"`
	OwnerID      int64         `json:"owner_id"`
	FromID       int64         `json:"from_id"`
	Date         int64         `json:"date"`
	Text         string        `json:"text"`
	ReplyOwnerID int64         `json:"reply_owner_id"`
	ReplyPostID  int64         `json:"reply_post_id"`
	FriendsOnly  int           `json:"friends_only"`
	Comments     *Comments     `json:"comments"`
	Likes        *Likes        `json:"likes"`
	Reposts      *Reposts      `json:"reposts"`
	Views        *Views        `json:"views"`
	PostType     string        `json:"post_type"`
	PostSource   *PostSource   `json:"post_source"`
	Attachments  *[]Attachment `json:"attachments"`
	Geo          *Geo          `json:"geo"`
	SingerID     int64         `json:"singer_id"`
	CopyHistory  *[]Wall       `json:"copy_history"`
	CanPin       int           `json:"can_pin"`
	CanDelete    int           `json:"can_delete"`
	CanEdit      int           `json:"can_edit"`
	IsPinned     int           `json:"is_pinned"`
	MarkedAsAds  int           `json:"marked_as_ads"`
}

func (w *Wall) URL() string {
	return "https://vk.com/wall" + strconv.FormatInt(w.OwnerID, 10) + "_" + strconv.FormatInt(w.ID, 10)
}

func (client *Client) GetWall(dst Destination, count int64, offset int64, filter string, extended bool, fields ...string) (int64, []Wall, []Users, []Group, *Error) {
	values := url.Values{}
	switch {
	case dst.UserID != 0:
		values.Set("owner_id", ConcatInt64ToString(dst.UserID))
	case dst.GroupID != 0:
		values.Set("owner_id", ConcatInt64ToString(-dst.GroupID))
	case dst.Domain != "":
		values.Set("domain", dst.Domain)
	}

	values.Set("count", ConcatInt64ToString(count))
	values.Set("offset", ConcatInt64ToString(offset))
	values.Set("filter", filter)

	if extended {
		values.Set("extended", "1")
		values.Set("fields", strings.Join(fields, ","))
	}

	res, err := client.Do(NewRequest("wall.get", "", values))
	if err != nil {
		return 0, []Wall{}, []Users{}, []Group{}, err
	}

	Answer := struct {
		Count    int64   `json:"count"`
		Items    []Wall  `json:"items"`
		Profiles []Users `json:"profiles"`
		Groups   []Group `json:"groups"`
	}{}

	if err := res.To(&Answer); err != nil {
		return 0, []Wall{}, []Users{}, []Group{}, NewError(ErrBadCode, err.Error())
	}

	return Answer.Count, Answer.Items, Answer.Profiles, Answer.Groups, nil
}

type Comments struct {
	Count   int `json:"count"`
	CanPost int `json:"can_post"`
}

type Likes struct {
	Count      int `json:"count"`
	UserLikes  int `json:"user_likes"`
	CanLike    int `json:"can_like"`
	CanPublish int `json:"can_publish"`
}

type Reposts struct {
	Count        int `json:"count"`
	UserReposted int `json:"user_reposted"`
}

type Views struct {
	Count int `json:"count"`
}

type PostSource struct {
	Type     string `json:"type"`
	Platform string `json:"platform"`
	Data     string `json:"data"`
	URL      string `json:"url"`
}

type Geo struct {
	Type        string `json:"type"`
	Coordinates string `json:"coordinates"`
	Place       struct {
		ID         int64   `json:"id"`
		Title      string  `json:"title"`
		Latitude   float64 `json:"latitude"`
		Longitude  float64 `json:"longitude"`
		Created    int64   `json:"created"`
		Icon       string  `json:"icon"`
		Country    string  `json:"country"`
		City       string  `json:"city"`
		Type       int     `json:"type"`
		GroupID    int64   `json:"group_id"`
		GroupPhoto int64   `json:"group_photo"`
		Checkins   int64   `json:"checkins"`
		Updated    int64   `json:"updated"`
		Address    int64   `json:"address"`
	} `json:"place"`
}
