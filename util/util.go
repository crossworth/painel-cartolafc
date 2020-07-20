package util

import (
	"strconv"
)

func ToInt(n string) int {
	i, _ := strconv.Atoi(n)
	return i
}

func ToIntWithDefault(n string, def int) int {
	i, err := strconv.Atoi(n)
	if err != nil {
		return def
	}

	return i
}

func ToString(n int) string {
	return strconv.Itoa(n)
}
