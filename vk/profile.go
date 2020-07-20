package vk

import (
	"errors"
)

var (
	ErrInvalidProfileURL = errors.New("link de perfil inválido")
)

// ProfileScreenNameOrIDFromURL retorna o screen name ou id de um link de perfil
func ProfileScreenNameOrIDFromURL(profileURL string) (string, error) {
	if !ProfileRegex.Match([]byte(profileURL)) {
		return "", ErrInvalidProfileURL
	}

	matches := ProfileRegex.FindStringSubmatch(profileURL)
	return matches[4], nil
}
