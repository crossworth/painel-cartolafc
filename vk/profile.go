package vk

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"

	vkapi "github.com/himidori/golang-vk-api"
)

var (
	ErrInvalidProfileURL = errors.New("link de perfil inválido")

	ErrProfileIDNotFound     = errors.New("id do perfil do usuário não encontrado")
	ErrInvalidScreenNameOrId = errors.New("o nome de perfil ou id é inválido")
)

// ProfileScreenNameOrIDFromURL retorna o screen name ou id de um link de perfil
func ProfileScreenNameOrIDFromURL(profileURL string) (string, error) {
	if !ProfileRegex.Match([]byte(profileURL)) {
		return "", ErrInvalidProfileURL
	}

	matches := ProfileRegex.FindStringSubmatch(profileURL)
	return matches[4], nil
}

func (v *VKClient) ScreenNameToProfileID(context context.Context, screenNameOrID string) (int, string, error) {
	if !ScreenNameOrIDMemberRegex.MatchString(screenNameOrID) {
		return 0, "", ErrInvalidScreenNameOrId
	}

	values := url.Values{}
	values.Add("user_ids", screenNameOrID)
	values.Add("fields", "screen_name")

	resp, err := v.client.MakeRequest("users.get", values)
	if err != nil {
		return 0, "", err
	}

	var userList []*vkapi.User
	err = json.Unmarshal(resp.Response, &userList)
	if err != nil {
		return 0, "", err
	}

	if len(userList) != 1 {
		return 0, "", ErrProfileIDNotFound
	}

	return userList[0].UID, userList[0].ScreenName, nil
}
