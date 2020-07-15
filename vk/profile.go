package vk

import (
	"errors"
	"regexp"
)

var (
	profileRegex = regexp.MustCompile(`(http(s?)://)?(m\.)?vk.com/([a-z0-9.]+)$`)
)

var (
	ErrInvalidProfileURL = errors.New("link de perfil inválido")
)

func ProfileScreenNameOrIDFromURL(profileURL string) (string, error) {
	if !profileRegex.Match([]byte(profileURL)) {
		return "", ErrInvalidProfileURL
	}

	matches := profileRegex.FindStringSubmatch(profileURL)
	return matches[4], nil
}
