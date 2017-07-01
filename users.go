package vkapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"strings"
)

func (client *Client) UsersInfo(userIds []string, fieldArgs []string) (users []Users, err *Error) {
	var req Request
	req.Method = "users.get"
	req.Values = url.Values{}

	if len(userIds) > 0 {
		ids := strings.Join(userIds, ",")
		req.Values.Set("user_ids", ids)
	}

	if len(userIds) > 0 {
		args := strings.Join(fieldArgs, ",")
		req.Values.Set("fields", args)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	log.Println("Answer:", res.Response.String())

	if err := json.Unmarshal(res.Response.Bytes(), &users); err != nil {
		return nil, NewError(ErrBadCode, err.Error())
	}

	return
}

func (client *Client) InitMyProfile(fieldArgs []string) error {
	users, err := client.UsersInfo([]string{}, fieldArgs)
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return errors.New("An unexpected error occurred.")
	}

	client.User = users[0]
	return nil
}

type Users struct {
	// Full description at https://vk.com/dev/objects/user
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Deactivated string `json:"deactivated"`
	Hidden      int    `json:"hidden"`

	// Optional fields
	About                  string          `json:"about"`
	Activities             string          `json:"activities"`
	Bdate                  string          `json:"bdate"`
	Blacklisted            int             `json:"blacklisted"`
	BlacklistedByMe        int             `json:"blacklisted_by_me"`
	Books                  string          `json:"books"`
	CanPost                int             `json:"can_post"`
	CanSeeAllPosts         int             `json:"can_see_all_posts"`
	CanSeeAudio            int             `json:"can_see_audio"`
	CanSendFriendRequest   int             `json:"can_send_friend_request"`
	CanWritePrivateMessage int             `json:"can_write_private_message"`
	Career                 *[]Career       `json:"career"`
	City                   *City           `json:"city"`
	CommonCount            int64           `json:"common_count"`
	Skype                  string          `json:"skype"`
	Facebook               string          `json:"facebook"`
	Twitter                string          `json:"twitter"`
	Livejournal            string          `json:"livejournal"`
	Instagram              string          `json:"instagram"`
	Contacts               *Contacts       `json:"contacts"`
	Counters               *Counters       `json:"counters"`
	Country                *Country        `json:"country"`
	CropPhoto              *CropPhoto      `json:"crop_photo"`
	Domain                 string          `json:"domain"`
	University             int64           `json:"university"`
	UniversityName         string          `json:"university_name"`
	Faculty                int64           `json:"faculty"`
	FacultyName            string          `json:"faculty_name"`
	Graduation             int64           `json:"graduation"`
	FirstNameNom           string          `json:"first_name_nom"`
	FirstNameGen           string          `json:"first_name_gen"`
	FirstNameDat           string          `json:"first_name_dat"`
	FirstNameAcc           string          `json:"first_name_acc"`
	FirstNameIns           string          `json:"first_name_ins"`
	FirstNameAbl           string          `json:"first_name_abl"`
	FollowersCount         int64           `json:"followers_count"`
	FriendStatus           int             `json:"friend_status"`
	Games                  string          `json:"games"`
	HasMobile              int             `json:"has_mobile"`
	HasPhoto               int             `json:"hasPhoto"`
	HomeTown               string          `json:"home_town"`
	Interests              string          `json:"interests"`
	IsFavorite             int             `json:"is_favorite"`
	IsFriend               int             `json:"is_friend"`
	IsHiddenFromFeed       int             `json:"is_hidden_from_feed"`
	LastNameNom            string          `json:"last_name_nom"`
	LastNameGen            string          `json:"last_name_gen"`
	LastNameDat            string          `json:"last_name_dat"`
	LastNameAcc            string          `json:"last_name_acc"`
	LastNameIns            string          `json:"last_name_ins"`
	LastNameAbl            string          `json:"last_name_abl"`
	LastSeen               *LastSeen       `json:"last_seen"`
	MaidenName             string          `json:"maiden_name"`
	Military               *[]Military     `json:"military"`
	Movies                 string          `json:"movies"`
	Music                  string          `json:"music"`
	Nickname               string          `json:"nickname"`
	Occupation             *Occupation     `json:"occupation"`
	Online                 int             `json:"online"`
	Personal               *Personal       `json:"personal"`
	Photo50                string          `json:"photo_50"`
	Photo100               string          `json:"photo_100"`
	Photo200Orig           string          `json:"photo_200_orig"`
	Photo200               string          `json:"photo_200"`
	Photo400Orig           string          `json:"photo_400_orig"`
	PhotoId                string          `json:"photo_id"`
	PhotoMax               string          `json:"photo_max"`
	PhotoMaxOrig           string          `json:"photo_max_orig"`
	Quotes                 string          `json:"quotes"`
	Relatives              *[]Relatives    `json:"relatives"`
	Relation               int             `json:"relation"`
	Schools                *[]Schools      `json:"schools"`
	ScreenName             string          `json:"screen_name"`
	Sex                    int             `json:"sex"`
	Site                   string          `json:"site"`
	Status                 string          `json:"status"`
	StatusAudio            string          `json:"status_audio"`
	Timezone               int             `json:"timezone"`
	Tv                     string          `json:"tv"`
	Universities           *[]Universities `json:"universities"`
	Verified               int             `json:"verified"`
	WallComments           int             `json:"wall_comments"`
}

type Career struct {
	GroupId   int64  `json:"group_id"`
	Company   string `json:"company"`
	CountryId int64  `json:"country_id"`
	CityId    int64  `json:"city_id"`
	CityName  string `json:"city_name"`
	From      int64  `json:"from"`
	Until     int64  `json:"until"`
	Position  string `json:"position"`
}

type City struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type Contacts struct {
	MobilePhone string `json:"mobile_phone"`
	HomePhone   string `json:"home_phone"`
}

type Counters struct {
	Albums        int `json:"albums"`
	Videos        int `json:"videos"`
	Audios        int `json:"audios"`
	Photos        int `json:"photos"`
	Notes         int `json:"notes"`
	Friends       int `json:"friends"`
	Groups        int `json:"groups"`
	OnlineFriends int `json:"online_friends"`
	MutualFriends int `json:"mutual_friends"`
	UserVideos    int `json:"user_videos"`
	Followers     int `json:"followers"`
	Pages         int `json:"pages"`
}

type Country struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type CropPhoto struct {
	/*Photo struct{} `json:"photo"`*/
	Crop *Rect `json:"crop"`
	Rect *Rect `json:"rect"`
}

type Rect struct {
	x  float64
	x2 float64
	y  float64
	y2 float64
}

type LastSeen struct {
	Time     int64 `json:"time"`
	Platform int   `json:"platform"`
}

type Military struct {
	Unit      string `json:"unit"`
	UnitId    int    `json:"unit_id"`
	CountryId int    `json:"country_id"`
	From      int    `json:"from"`
	Until     int    `json:"until"`
}

type Occupation struct {
	Type string `json:"type"`
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Personal struct {
	Political  int      `json:"political"`
	Langs      []string `json:"langs"`
	Religion   string   `json:"religion"`
	InspiredBy string   `json:"inspired_by"`
	PeopleMain int      `json:"people_main"`
	LifeMain   int      `json:"life_main"`
	Smoking    int      `json:"smoking"`
	Alcohol    int      `json:"alcohol"`
}

type Relatives struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Schools struct {
	Id            string `json:"id"`
	Country       int    `json:"country"`
	City          int    `json:"city"`
	Name          string `json:"name"`
	YearFrom      int    `json:"year_from"`
	YearTo        int    `json:"year_to"`
	YearGraduated int    `json:"year_graduated"`
	Class         string `json:"class"`
	Speciality    string `json:"speciality"`
	Type          int    `json:"type"`
	TypeStr       string `json:"type_str"`
}

type Universities struct {
	Id              int    `json:"id"`
	Country         int    `json:"country"`
	City            int    `json:"city"`
	Name            string `json:"name"`
	Faculty         int    `json:"faculty"`
	FacultyName     string `json:"faculty_name"`
	Chair           int    `json:"chair"`
	ChairName       string `json:"chair_name"`
	Graduation      int    `json:"graduation"`
	EducationForm   string `json:"education_form"`
	EducationStatus string `json:"education_status"`
}

var AllFields = []string{"about",
	"activities",
	"bdate",
	"blacklisted",
	"blacklisted_by_me",
	"books",
	"can_post",
	"can_see_all_posts",
	"can_see_audio",
	"can_send_friend_request",
	"can_write_private_message",
	"career",
	"city",
	"common_count",
	"connections",
	"contacts",
	"counters",
	"country",
	"crop_photo",
	"domain",
	"education",
	"first_name_nom",
	"first_name_gen",
	"first_name_dat",
	"first_name_acc",
	"first_name_ins",
	"first_name_abl",
	"followers_count",
	"friend_status",
	"games",
	"has_mobile",
	"hasPhoto",
	"home_town",
	"interests",
	"is_favorite",
	"is_friend",
	"is_hidden_from_feed",
	"last_name_nom",
	"last_name_gen",
	"last_name_dat",
	"last_name_acc",
	"last_name_ins",
	"last_name_abl",
	"last_seen",
	"maiden_name",
	"military",
	"movies",
	"music",
	"nickname",
	"occupation",
	"online",
	"personal",
	"photo_50",
	"photo_100",
	"photo_200_orig",
	"photo_200",
	"photo_400_orig",
	"photo_id",
	"photo_max",
	"photo_max_orig",
	"quotes",
	"relatives",
	"relation",
	"schools",
	"screen_name",
	"sex",
	"site",
	"status",
	"status_audio",
	"timezone",
	"tv",
	"universities",
	"verified",
	"wall_comments"}
