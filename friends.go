package vkapi

import (
	"encoding/json"
	"net/url"
	"strings"
)

const (
	FriendFieldHints  = "hints"
	FriendFieldRandom = "random"
	FriendFieldMobile = "mobile"
	FriendFieldName   = "name"
)

// GetFriends will be return array userID or array Users.
func (client *Client) GetFriends(userID int64, order string, count int64, offset int64, nameCase string, fields ...string) (friends []Users, err *Error) {
	values := url.Values{}
	values.Add("user_id", ConcatInt64ToString(userID))
	values.Add("order", order)
	values.Add("count", ConcatInt64ToString(count))
	values.Add("offset", ConcatInt64ToString(offset))
	values.Add("name_case", nameCase)
	values.Add("fields", strings.Join(fields, ","))

	res, err := client.Do(NewRequest("friends.get", "", values))
	if err != nil {
		return []Users{}, err
	}

	Answer := struct {
		Items Raw `json:"items"`
	}{}

	if err := res.To(&Answer); err != nil {
		return []Users{}, NewError(ErrBadCode, err.Error())
	}

	var ids []int64
	if err := json.Unmarshal(Answer.Items.Bytes(), &ids); err == nil {
		for _, id := range ids {
			friends = append(friends, Users{ID: id})
		}

		return friends, nil
	}

	if err := json.Unmarshal(Answer.Items.Bytes(), &friends); err != nil {
		return []Users{}, NewError(ErrBadCode, err.Error())
	}

	return friends, nil
}
