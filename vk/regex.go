package vk

import (
	"regexp"
)

var (
	ProfileRegex        = regexp.MustCompile(`(http(s?)://)?(m\.)?vk.com/([a-z0-9._]+)$`)
	ScreenNameOrIDRegex = regexp.MustCompile(`([a-z0-9._]+)$`)
)
