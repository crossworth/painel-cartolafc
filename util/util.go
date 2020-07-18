package util

import (
	"strconv"
)

func ToInt(n string) int {
	i, _ := strconv.Atoi(n)
	return i
}

func ToInt64(n string) int64 {
	i, _ := strconv.ParseInt(n, 10, 64)
	return i
}

func ToIntWithDefault(n string, def int) int {
	i, err := strconv.Atoi(n)
	if err != nil {
		return def
	}

	return i
}

func ToIntWithDefault64(n string, def int64) int64 {
	i, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		return def
	}

	return i
}

func ToString(n int) string {
	return strconv.Itoa(n)
}
