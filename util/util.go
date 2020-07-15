package util

import (
	"strconv"
)

func ToInt(n string) int {
	i, _ := strconv.Atoi(n)
	return i
}

func ToString(n int) string {
	return strconv.Itoa(n)
}
