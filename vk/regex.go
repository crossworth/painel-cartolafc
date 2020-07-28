package vk

import (
	"regexp"
)

var (
	ProfileRegex              = regexp.MustCompile(`(http(s?)://)?(m\.)?vk.com/([\-?a-z0-9._]+)$`)
	ScreenNameOrIDMemberRegex = regexp.MustCompile(`([a-z0-9._]+)$`)
	ScreenNameOrIDGroupRegex  = regexp.MustCompile(`([\-a-z0-9._]+)$`)
)
