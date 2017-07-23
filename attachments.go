package vkapi

import "fmt"

type Attachment struct {
	Type      string    `json:"type"`
	AccessKey string    `json:"access_key"`
	Photo     *Photo    `json:"photo"`
	Audio     *Audio    `json:"audio"`
	Video     *Video    `json:"video"`
	Document  *Document `json:"doc"`
}

type Photo struct {
	ID        int64    `json:"id"`
	AlbumID   int64    `json:"album_id"`
	OwnerID   int64    `json:"owner_id"`
	UserID    int64    `json:"user_id"`
	Text      string   `json:"text"`
	Date      int64    `json:"date"`
	Sizes     *[]Sizes `json:"sizes"`
	Photo75   string   `json:"photo_75"`
	Photo130  string   `json:"photo_130"`
	Photo604  string   `json:"photo_604"`
	Photo807  string   `json:"photo_807"`
	Photo1280 string   `json:"photo_1280"`
	Photo2560 string   `json:"photo_2560"`
	Width     int      `json:"width"`
	Height    int      `json:"height"`
}

type Sizes struct {
	Src    string `json:"src"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Type   string `json:"type"`
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

type Document struct {
	ID        int64  `json:"id"`
	OwnerID   int64  `json:"owner_id"`
	Title     string `json:"title"`
	Size      int64  `json:"size"`
	Ext       string `json:"ext"`
	Url       string `json:"url"`
	Date      int64  `json:"date"`
	Type      int    `json:"type"`
	AccessKey string `json:"access_key"`
	Preview   struct {
		Photo *struct {
			Sizes *[]Sizes `json:"sizes"`
		} `json:"photo"`
		Graffiti *Graffiti `json:"graffiti"`
		AudioMsg *AudioMsg `json:"audio_msg"`
	} `json:"preview"`
}

func (doc *Document) IsTxt() bool {
	return doc.Type == 1
}

func (doc *Document) IsArch() bool {
	return doc.Type == 2
}

func (doc *Document) IsGif() bool {
	return doc.Type == 3
}

func (doc *Document) IsImages() bool {
	return doc.Type == 4
}

func (doc *Document) IsAudio() bool {
	return doc.Type == 5
}

func (doc *Document) IsVideo() bool {
	return doc.Type == 6
}

func (doc *Document) IsEBooks() bool {
	return doc.Type == 7
}

func (doc *Document) IsUnknown() bool {
	return doc.Type == 8
}

type Graffiti struct {
	Src    string `json:"src"`
	Width  int64  `json:"width"`
	Height int64  `json:"height"`
}

type AudioMsg struct {
	Duration int64   `json:"duration"`
	Waveform []int64 `json:"waveform"`
	LinkOgg  string  `json:"link_ogg"`
	LinkMp3  string  `json:"link_mp3"`
}

type Video struct {
	ID          int64  `json:"id"`
	OwnerID     int64  `json:"owner_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int64  `json:"duration"`
	Photo130    string `json:"photo_130"`
	Photo320    string `json:"photo_320"`
	Photo640    string `json:"photo_640"`
	Photo800    string `json:"photo_800"`
	Date        int64  `json:"date"`
	AddingDate  int64  `json:"adding_date"`
	Views       int64  `json:"views"`
	Comments    int64  `json:"comments"`
	Player      string `json:"player"`
	AccessKey   string `json:"access_key"`
	Processing  int    `json:"processing"`
	Live        int    `json:"live"`
	Upcoming    int    `json:"upcoming"`
}

func (v *Video) GetMaxPreview() string {
	if v.Photo800 != "" {
		return v.Photo800
	}
	if v.Photo640 != "" {
		return v.Photo640
	}
	if v.Photo320 != "" {
		return v.Photo320
	}
	if v.Photo130 != "" {
		return v.Photo130
	}

	return ""
}

func (client *Client) AddAttachment(fieldname string, file interface{}) string {
	switch fieldname {
	case "photo":
		server, err := client.GetMessagesUploadServer()
		if err != nil {
			return ""
		}

		res, err := client.UploadFile(server.UploadURL, fieldname, file)
		if err != nil {
			return ""
		}

		photo, err := client.SaveMessagesPhoto(res)
		if err != nil {
			return ""
		}

		return fmt.Sprintf("photo%d_%d", photo.OwnerID, photo.ID)
	default:
		return ""
	}
}

func (client *Client) SendPhoto(dst Destination, file interface{}) (int64, *Error) {
	config := MessageConfig{
		Destination: dst,
		Attachment:  client.AddAttachment("photo", file),
	}

	return client.SendMessage(config)
}
