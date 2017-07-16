package vkapi

type Attachment struct {
	Type      string `json:"type"`
	AccessKey string `json:"access_key"`
	Photo     *Photo `json:"photo"`
	Audio     *Audio `json:"audio"`
}

type Photo struct {
	ID      int64  `json:"id"`
	AlbumID int64  `json:"album_id"`
	OwnerID int64  `json:"owner_id"`
	UserID  int64  `json:"user_id"`
	Text    string `json:"text"`
	Date    int64  `json:"date"`
	Sizes   *[]struct {
		Src    string `json:"src"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
		Type   string `json:"type"`
	} `json:"sizes"`
	Photo75   string `json:"photo_75"`
	Photo130  string `json:"photo_130"`
	Photo604  string `json:"photo_604"`
	Photo807  string `json:"photo_807"`
	Photo1280 string `json:"photo_1280"`
	Photo2560 string `json:"photo_2560"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

func (p *Photo) GetMaxSizePhoto() string {
	if p.Photo2560 != "" {
		return p.Photo2560
	}
	if p.Photo1280 != "" {
		return p.Photo1280
	}
	if p.Photo807 != "" {
		return p.Photo807
	}
	if p.Photo604 != "" {
		return p.Photo604
	}
	if p.Photo130 != "" {
		return p.Photo130
	}
	if p.Photo75 != "" {
		return p.Photo75
	}

	return ""
}

type Audio struct {
	ID       int64  `json:"id"`
	OwnerID  int64  `json:"owner_id"`
	Artist   string `json:"artist"`
	Title    string `json:"title"`
	Duration int64  `json:"duration"`
	Url      string `json:"url"`
	LyricsID int64  `json:"lyrics_id"`
	AlbumID  int64  `json:"album_id"`
	GenreID  int64  `json:"genre_id"`
	Date     int64  `json:"date"`
	NoSearch int    `json:"no_search"`
}
