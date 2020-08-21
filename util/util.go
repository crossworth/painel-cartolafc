package util

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/crossworth/cartola-web-admin/logger"
)

// Deprecated: Remove this
func ToInt(n string) int {
	i, _ := strconv.Atoi(n)
	return i
}

// Deprecated: Remove this
func ToIntWithDefault(n string, def int) int {
	i, err := strconv.Atoi(n)
	if err != nil {
		return def
	}

	return i
}

// Deprecated: Remove this
func ToIntWithDefaultMin(n string, def int) int {
	i, err := strconv.Atoi(n)
	if err != nil {
		return def
	}

	if i < def {
		return def
	}

	return i
}

// Deprecated: Remove this
func ToString(n int) string {
	return strconv.Itoa(n)
}

// Deprecated: Remove this
func StringWithDefault(n string, def string) string {
	if n != "" {
		return n
	}

	return def
}

// Deprecated: Remove this
func BoolWithDefault(n string, def bool) bool {
	if strings.ToLower(n) == "true" || strings.ToLower(n) == "1" || strings.ToLower(n) == "on" {
		return true
	}

	if strings.ToLower(n) == "false" || strings.ToLower(n) == "0" || strings.ToLower(n) == "off" {
		return false
	}

	return def
}

func StringOrDefault(input string, defaultStr string) string {
	if input == "" {
		return defaultStr
	}

	return input
}

func GetStringFromEnvOrDefault(key string, defaultStr string) string {
	return StringOrDefault(os.Getenv(key), defaultStr)
}

func GetStringFromEnvOrFatalError(key string) string {
	envContent := os.Getenv(key)

	if envContent == "" {
		logger.Log.Fatal().Msg(fmt.Sprintf("variável de ambiente %q não definida", key))
	}

	return envContent
}

func IntJoin(ints []int, sep string) string {
	var result []string

	for _, i := range ints {
		result = append(result, strconv.Itoa(i))
	}

	return strings.Join(result, sep)
}

func StringToIntSlice(input string) []int {
	parts := strings.Split(input, ",")
	var result []int

	for _, i := range parts {
		n, _ := strconv.Atoi(i)
		result = append(result, n)
	}

	return result
}

func GetIntFromEnvOrFatalError(key string) int {
	envContent := os.Getenv(key)

	if envContent == "" {
		logger.Log.Fatal().Msg(fmt.Sprintf("variável de ambiente %q não definida", key))
	}

	intVal, err := strconv.Atoi(envContent)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg(fmt.Sprintf("variável de ambiente %q não é um número inteiro, %q", key, envContent))
	}

	return intVal
}
