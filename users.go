package vkapi

import (
	"fmt"
	"net/url"
	"strings"
)

// VKUser describes the structure of the vk user.
type VKUser struct {
	Me      Users
	Friends []Users
}

// UsersInfo returns array Users with the selected fields
// if the request was successful.
func (client *Client) UsersInfo(dst Destination, fieldArgs ...string) (users []Users, err *Error) {
	values := url.Values{}
	if dst.ScreenName != "" {
		values.Add("user_ids", dst.ScreenName)
	} else {
		values = dst.Values()
	}

	if len(fieldArgs) > 0 {
		args := strings.Join(fieldArgs, ",")
		values.Set("fields", args)
	}

	res, err := client.Do(NewRequest("users.get", "", values))
	if err != nil {
		return nil, err
	}

	if err := res.To(&users); err != nil {
		return nil, NewError(ErrBadCode, err.Error())
	}

	return
}

// InitMyProfile fills in the selected Client.VKUser data.
func (client *Client) InitMyProfile(fieldArgs ...string) *Error {
	users, err := client.UsersInfo(Destination{}, fieldArgs...)
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return NewError(ErrBadCode, "An unexpected error occurred.")
	}

	client.VKUser.Me = users[0]
	return nil
}

// Users describes the structure of the users.
type Users struct {
	// Full description at https://vk.com/dev/objects/user
	ID          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Deactivated string `json:"deactivated"`
	Hidden      int    `json:"hidden"`

	// Optional fields
	About                  string          `json:"about,omitempty"`
	Activities             string          `json:"activities,omitempty"`
	Bdate                  string          `json:"bdate,omitempty"`
	Blacklisted            int             `json:"blacklisted,omitempty"`
	BlacklistedByMe        int             `json:"blacklisted_by_me,omitempty"`
	Books                  string          `json:"books,omitempty"`
	CanPost                int             `json:"can_post,omitempty"`
	CanSeeAllPosts         int             `json:"can_see_all_posts,omitempty"`
	CanSeeAudio            int             `json:"can_see_audio,omitempty"`
	CanSendFriendRequest   int             `json:"can_send_friend_request,omitempty"`
	CanWritePrivateMessage int             `json:"can_write_private_message,omitempty"`
	Career                 *[]Career       `json:"career,omitempty"`
	City                   *City           `json:"city,omitempty"`
	CommonCount            int64           `json:"common_count,omitempty"`
	Skype                  string          `json:"skype,omitempty"`
	Facebook               string          `json:"facebook,omitempty"`
	Twitter                string          `json:"twitter,omitempty"`
	Livejournal            string          `json:"livejournal,omitempty"`
	Instagram              string          `json:"instagram,omitempty"`
	Contacts               *Contacts       `json:"contacts,omitempty"`
	Counters               *Counters       `json:"counters,omitempty"`
	Country                *Country        `json:"country,omitempty"`
	CropPhoto              *CropPhoto      `json:"crop_photo,omitempty"`
	Domain                 string          `json:"domain,omitempty"`
	University             int64           `json:"university,omitempty"`
	UniversityName         string          `json:"university_name,omitempty"`
	Faculty                int64           `json:"faculty,omitempty"`
	FacultyName            string          `json:"faculty_name,omitempty"`
	Graduation             int64           `json:"graduation,omitempty"`
	FirstNameNom           string          `json:"first_name_nom,omitempty"`
	FirstNameGen           string          `json:"first_name_gen,omitempty"`
	FirstNameDat           string          `json:"first_name_dat,omitempty"`
	FirstNameAcc           string          `json:"first_name_acc,omitempty"`
	FirstNameIns           string          `json:"first_name_ins,omitempty"`
	FirstNameAbl           string          `json:"first_name_abl,omitempty"`
	FollowersCount         int64           `json:"followers_count,omitempty"`
	FriendStatus           int             `json:"friend_status,omitempty"`
	Games                  string          `json:"games,omitempty"`
	HasMobile              int             `json:"has_mobile,omitempty"`
	HasPhoto               int             `json:"hasPhoto,omitempty"`
	HomeTown               string          `json:"home_town,omitempty"`
	Interests              string          `json:"interests,omitempty"`
	IsFavorite             int             `json:"is_favorite,omitempty"`
	IsFriend               int             `json:"is_friend,omitempty"`
	IsHiddenFromFeed       int             `json:"is_hidden_from_feed,omitempty"`
	LastNameNom            string          `json:"last_name_nom,omitempty"`
	LastNameGen            string          `json:"last_name_gen,omitempty"`
	LastNameDat            string          `json:"last_name_dat,omitempty"`
	LastNameAcc            string          `json:"last_name_acc,omitempty"`
	LastNameIns            string          `json:"last_name_ins,omitempty"`
	LastNameAbl            string          `json:"last_name_abl,omitempty"`
	LastSeen               *LastSeen       `json:"last_seen,omitempty"`
	MaidenName             string          `json:"maiden_name,omitempty"`
	Military               *[]Military     `json:"military,omitempty"`
	Movies                 string          `json:"movies,omitempty"`
	Music                  string          `json:"music,omitempty"`
	Nickname               string          `json:"nickname,omitempty"`
	Occupation             *Occupation     `json:"occupation,omitempty"`
	Online                 int             `json:"online,omitempty"`
	Personal               *Personal       `json:"personal,omitempty"`
	Photo50                string          `json:"photo_50,omitempty"`
	Photo100               string          `json:"photo_100,omitempty"`
	Photo200Orig           string          `json:"photo_200_orig,omitempty"`
	Photo200               string          `json:"photo_200,omitempty"`
	Photo400Orig           string          `json:"photo_400_orig,omitempty"`
	PhotoId                string          `json:"photo_id,omitempty"`
	PhotoMax               string          `json:"photo_max,omitempty"`
	PhotoMaxOrig           string          `json:"photo_max_orig,omitempty"`
	Quotes                 string          `json:"quotes,omitempty"`
	Relatives              *[]Relatives    `json:"relatives,omitempty"`
	Relation               int             `json:"relation,omitempty"`
	Schools                *[]Schools      `json:"schools,omitempty"`
	ScreenName             string          `json:"screen_name,omitempty"`
	Sex                    int             `json:"sex,omitempty"`
	Site                   string          `json:"site,omitempty"`
	Status                 string          `json:"status,omitempty"`
	StatusAudio            string          `json:"status_audio,omitempty"`
	Timezone               int             `json:"timezone,omitempty"`
	Tv                     string          `json:"tv,omitempty"`
	Universities           *[]Universities `json:"universities,omitempty"`
	Verified               int             `json:"verified,omitempty"`
	WallComments           int             `json:"wall_comments,omitempty"`
}

func (user *Users) MainInfo(sep string) string {
	return fmt.Sprintf("ID: %d%sFirst name: %s%sLast name: %s", user.ID, sep, user.FirstName, sep, user.LastName)
}

func (client *Client) GetMainInfo(dst Destination, sep string) ([]string, *Error) {
	users, err := client.UsersInfo(dst)
	if err != nil {
		return []string{}, err
	}

	var string []string
	for _, u := range users {
		string = append(string, u.MainInfo(sep))
	}

	return string, nil
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
	Photo *Photo `json:"photo"`
	Crop  *Rect  `json:"crop"`
	Rect  *Rect  `json:"rect"`
}

type Rect struct {
	X  float64
	X2 float64
	Y  float64
	Y2 float64
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

const (
	UserFieldAbout                  = "about"
	UserFieldActivities             = "activities"
	UserFieldBirthdayDate           = "bdate"
	UserFieldBlacklisted            = "blacklisted"
	UserFieldBlacklistedByMe        = "blacklisted_by_me"
	UserFieldBooks                  = "books"
	UserFieldCanPost                = "can_post"
	UserFieldCanSeeAllPosts         = "can_see_all_posts"
	UserFieldCanSeeAudio            = "can_see_audio"
	UserFieldCanSendFriendRequest   = "can_send_friend_request"
	UserFieldCanWritePrivateMessage = "can_write_private_message"
	UserFieldCarrer                 = "career"
	UserFieldCity                   = "city"
	UserFieldCommonCount            = "common_count"
	UserFieldConnections            = "connections"
	UserFieldContacts               = "contacts"
	UserFieldCounters               = "counters"
	UserFieldCountry                = "country"
	UserFieldCropPhoto              = "crop_photo"
	UserFieldDomain                 = "domain"
	UserFieldEducation              = "education"
	UserFieldFirstNameNom           = "first_name_nom"
	UserFieldFirstNameGen           = "first_name_gen"
	UserFieldFirstNameDat           = "first_name_dat"
	UserFieldFirstNameAcc           = "first_name_acc"
	UserFieldFirstNameIns           = "first_name_ins"
	UserFieldFirstNameAbl           = "first_name_abl"
	UserFieldFollowersCount         = "followers_count"
	UserFieldFriendStatus           = "friend_status"
	UserFieldGames                  = "games"
	UserFieldHasMobile              = "has_mobile"
	UserFieldHasPhoto               = "hasPhoto"
	UserFieldHomeTown               = "home_town"
	UserFieldInterests              = "interests"
	UserFieldIsFavorite             = "is_favorite"
	UserFieldIsFriend               = "is_friend"
	UserFieldIsHiddenFromFeed       = "is_hidden_from_feed"
	UserFieldLastNameNom            = "last_name_nom"
	UserFieldLastNameGen            = "last_name_gen"
	UserFieldLastNameDat            = "last_name_dat"
	UserFieldLastNameAcc            = "last_name_acc"
	UserFieldLastNameIns            = "last_name_ins"
	UserFieldLastNameAbl            = "last_name_abl"
	UserFieldLastSeen               = "last_seen"
	UserFieldMaidenName             = "maiden_name"
	UserFieldMilitary               = "military"
	UserFieldMovies                 = "movies"
	UserFieldMusic                  = "music"
	UserFieldNickname               = "nickname"
	UserFieldOccupation             = "occupation"
	UserFieldOnline                 = "online"
	UserFieldPersonal               = "personal"
	UserFieldPhoto50                = "photo_50"
	UserFieldPhoto100               = "photo_100"
	UserFieldPhoto200Orig           = "photo_200_orig"
	UserFieldPhoto200               = "photo_200"
	UserFieldPhoto400Orig           = "photo_400_orig"
	UserFieldPhotoId                = "photo_id"
	UserFieldPhotoMax               = "photo_max"
	UserFieldPhotoMaxOrig           = "photo_max_orig"
	UserFieldQuotes                 = "quotes"
	UserFieldRelatives              = "relatives"
	UserFieldRelation               = "relation"
	UserFieldSchool                 = "schools"
	UserFieldScreenName             = "screen_name"
	UserFieldSex                    = "sex"
	UserFieldSite                   = "site"
	UserFieldStatus                 = "status"
	UserFieldStatusAudio            = "status_audio"
	UserFieldTimezone               = "timezone"
	UserFieldTv                     = "tv"
	UserFieldUniversities           = "universities"
	UserFieldVerified               = "verified"
	UserFieldWallComments           = "wall_comments"
)

var UserFieldAll = []string{UserFieldAbout,
	UserFieldActivities,
	UserFieldBirthdayDate,
	UserFieldBlacklisted,
	UserFieldBlacklistedByMe,
	UserFieldBooks,
	UserFieldCanPost,
	UserFieldCanSeeAllPosts,
	UserFieldCanSeeAudio,
	UserFieldCanSendFriendRequest,
	UserFieldCanWritePrivateMessage,
	UserFieldCarrer,
	UserFieldCity,
	UserFieldCommonCount,
	UserFieldConnections,
	UserFieldContacts,
	UserFieldCounters,
	UserFieldCountry,
	UserFieldCropPhoto,
	UserFieldDomain,
	UserFieldEducation,
	UserFieldFirstNameNom,
	UserFieldFirstNameGen,
	UserFieldFirstNameDat,
	UserFieldFirstNameAcc,
	UserFieldFirstNameIns,
	UserFieldFirstNameAbl,
	UserFieldFollowersCount,
	UserFieldFriendStatus,
	UserFieldGames,
	UserFieldHasMobile,
	UserFieldHasPhoto,
	UserFieldHomeTown,
	UserFieldInterests,
	UserFieldIsFavorite,
	UserFieldIsFriend,
	UserFieldIsHiddenFromFeed,
	UserFieldLastNameNom,
	UserFieldLastNameGen,
	UserFieldLastNameDat,
	UserFieldLastNameAcc,
	UserFieldLastNameIns,
	UserFieldLastNameAbl,
	UserFieldLastSeen,
	UserFieldMaidenName,
	UserFieldMilitary,
	UserFieldMovies,
	UserFieldMusic,
	UserFieldNickname,
	UserFieldOccupation,
	UserFieldOnline,
	UserFieldPersonal,
	UserFieldPhoto50,
	UserFieldPhoto100,
	UserFieldPhoto200Orig,
	UserFieldPhoto200,
	UserFieldPhoto400Orig,
	UserFieldPhotoId,
	UserFieldPhotoMax,
	UserFieldPhotoMaxOrig,
	UserFieldQuotes,
	UserFieldRelatives,
	UserFieldRelation,
	UserFieldSchool,
	UserFieldScreenName,
	UserFieldSex,
	UserFieldSite,
	UserFieldStatus,
	UserFieldStatusAudio,
	UserFieldTimezone,
	UserFieldTv,
	UserFieldUniversities,
	UserFieldVerified,
	UserFieldWallComments}
