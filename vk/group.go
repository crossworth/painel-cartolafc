package vk

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
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

type isMemberResponse struct {
	Member    bool `json:"member"`
	CanInvite bool `json:"can_invite"`
}

func (v *VKClient) IsUserIDOnGroup(context context.Context, userID string, communityID int) (bool, error) {

	values := url.Values{}
	values.Add("group_id", strconv.Itoa(communityID))
	values.Add("user_id", userID)

	resp, err := v.client.MakeRequest("groups.isMember", values)
	if err != nil {
		return false, err
	}

	var member isMemberResponse
	err = json.Unmarshal(resp.Response, &member)
	if err != nil {
		var memberNumber int
		err2 := json.Unmarshal(resp.Response, &memberNumber)
		if err2 == nil {
			return memberNumber == 1, nil
		}

		return false, err
	}

	return member.Member, nil
}
