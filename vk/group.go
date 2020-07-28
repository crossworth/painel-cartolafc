package vk

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"strings"
)

type GroupInfo struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ScreenName   string `json:"screen_name"`
	IsClosed     int    `json:"is_closed"`
	Type         string `json:"type"`
	IsAdmin      int    `json:"is_admin"`
	IsMember     int    `json:"is_member"`
	IsAdvertiser int    `json:"is_advertiser"`
	Photo50      string `json:"photo_50"`
	Photo100     string `json:"photo_100"`
	Photo200     string `json:"photo_200"`
}

func (v *VKClient) IsGroup(id string) bool {
	return ScreenNameOrIDGroupRegex.MatchString(id)
}

func (v *VKClient) GroupScreenNameToProfileID(context context.Context, screenNameOrID string) (int, string, error) {
	if !ScreenNameOrIDGroupRegex.MatchString(screenNameOrID) {
		return 0, "", ErrInvalidScreenNameOrId
	}

	screenNameOrID = strings.TrimPrefix(screenNameOrID, "-")

	values := url.Values{}
	values.Add("group_ids", "")
	values.Add("group_id", screenNameOrID)

	log.Println("haa")

	resp, err := v.client.MakeRequest("groups.getById", values)
	if err != nil {
		return 0, "", err
	}

	var groupInfo []*GroupInfo
	err = json.Unmarshal(resp.Response, &groupInfo)
	if err != nil {
		return 0, "", err
	}

	if len(groupInfo) != 1 {
		return 0, "", ErrProfileIDNotFound
	}

	return -groupInfo[0].ID, groupInfo[0].ScreenName, nil
}
