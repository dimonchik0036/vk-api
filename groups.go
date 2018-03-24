package vkapi

import (
	"strings"
)

type VKGroup struct {
	Members []int64
}

type Group struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	ScreenName  string `json:"screen_name"`
	IsClosed    int    `json:"is_closed"`
	Deactivated string `json:"deactivated"`
	IsAdmin     int    `json:"is_admin"`
	AdminLevel  int    `json:"admin_level"`
	IsMember    int    `json:"is_member"`
	InvitedBy   int    `json:"invited_by"`
	Type        string `json:"type"`
	Photo50     string `json:"photo_50"`
	Photo100    string `json:"photo_100"`
	Photo200    string `json:"photo_200"`
	/* TODO Add options */
}

func (client *Client) GetGroupById(dst Destination, fields ...string) ([]Group, *Error) {
	values := dst.Values()

	values.Set("fields", strings.Join(fields, ","))

	res, err := client.Do(NewRequest("groups.getById", "", values))
	if err != nil {
		return []Group{}, err
	}

	var Groups []Group

	if err := res.To(&Groups); err != nil {
		return []Group{}, NewError(ErrBadCode, err.Error())
	}

	return Groups, nil
}
