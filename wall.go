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
	ToID         int64         `json:"to_id"`
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
	var peerID int64
	if w.OwnerID == 0 {
		peerID = w.ToID
	}

	return "https://vk.com/wall" + strconv.FormatInt(peerID, 10) + "_" + strconv.FormatInt(w.ID, 10)
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

type PostConfig struct {
	OwnerID       int64     `json:"owner_id"`
	FriendsOnly   int       `json:"friends_only"`
	FromGroup     int       `json:"from_group"`
	Attachments   string    `json:"attachments"`
	Message       string    `json:"message"`
	Services      string    `json:"services"`
	PublishDate   Timestamp `json:"publish_date"`
	Signed        int       `json:"signed"`
	Geo           bool      `json:"-"`
	Lat           float64   `json:"lat"`
	Long          float64   `json:"long"`
	PlaceID       int64     `json:"place_id"`
	PostID        int64     `json:"post_id"`
	Guid          string    `json:"guid"`
	MarkAsAds     int       `json:"mark_as_ads"`
	CloseComments int       `json:"close_comments"`
}

// SetGeo sets the location.
func (p *PostConfig) SetGeo(lat float64, long float64) {
	p.Geo = true
	p.Lat = lat
	p.Long = long
}

// PostWall tries to post a message to the wall with a configuration
// from PostConfig and returns the post id if it succeeds.
func (client *Client) PostWall(config PostConfig) (int64, *Error) {
	values := url.Values{}
	if config.OwnerID != 0 {
		values.Set("owner_id", ConcatInt64ToString(config.OwnerID))
	}
	if config.Message != "" {
		values.Add("message", config.Message)
	}
	if config.Geo {
		values.Add("lat", strconv.FormatFloat(config.Lat, 'f', -1, 64))
		values.Add("long", strconv.FormatFloat(config.Long, 'f', -1, 64))
	}
	if config.Attachments != "" {
		values.Add("attachments", config.Attachments)
	}
	if config.FriendsOnly != 0 {
		values.Set("friends_only", string(config.FriendsOnly))
	}
	if config.FromGroup != 0 {
		values.Set("from_group", string(config.FromGroup))
	}
	if config.Services != "" {
		values.Add("services", config.Services)
	}
	if config.Signed != 0 {
		values.Set("signed", string(config.Signed))
	}
	if config.PublishDate != 0 {
		values.Set("publish_date", ConcatInt64ToString(int64(config.PublishDate)))
	}
	if config.PlaceID != 0 {
		values.Set("place_id", ConcatInt64ToString(config.PlaceID))
	}
	if config.PostID != 0 {
		values.Set("post_id", ConcatInt64ToString(config.PostID))
	}
	if config.Guid != "" {
		values.Add("guid", config.Guid)
	}
	if config.MarkAsAds != 0 {
		values.Set("mark_as_ads", string(config.MarkAsAds))
	}
	if config.CloseComments != 0 {
		values.Set("close_comments", string(config.CloseComments))
	}

	res, err := client.Do(NewRequest("wall.post", "", values))
	if err != nil {
		return 0, err
	}

	Answer := struct {
		PostID int64 `json:"post_id"`
	}{}

	if err := res.To(&Answer); err != nil {
		return 0, NewError(ErrBadCode, err.Error())
	}

	return Answer.PostID, nil
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
