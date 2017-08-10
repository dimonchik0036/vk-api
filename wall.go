package vkapi

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
		ID         int64  `json:"id"`
		Title      string `json:"title"`
		Latitude   int    `json:"latitude"`
		Longitude  int    `json:"longitude"`
		Created    int64  `json:"created"`
		Icon       string `json:"icon"`
		Country    string `json:"country"`
		City       string `json:"city"`
		Type       int    `json:"type"`
		GroupID    int64  `json:"group_id"`
		GroupPhoto int64  `json:"group_photo"`
		Checkins   int64  `json:"checkins"`
		Updated    int64  `json:"updated"`
		Address    int64  `json:"address"`
	} `json:"place"`
}
