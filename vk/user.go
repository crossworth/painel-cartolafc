package vk

import (
	"encoding/json"
	"errors"
	"net/url"

	vkapi "github.com/himidori/golang-vk-api"
)

var (
	ErrUserIDNotFound = errors.New("id do perfil do usuário não encontrado")
)

func (v *VKClient) ScreenNameToUserID(screenNameOrID string) (int, error) {
	values := url.Values{}
	values.Add("user_ids", screenNameOrID)

	resp, err := v.client.MakeRequest("users.get", values)
	if err != nil {
		return 0, err
	}

	var userList []*vkapi.User
	err = json.Unmarshal(resp.Response, &userList)
	if err != nil {
		return 0, err
	}

	if len(userList) != 1 {
		return 0, ErrUserIDNotFound
	}

	return userList[0].UID, nil
}
