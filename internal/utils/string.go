package utils

import (
	"regexp"
	"strings"
)

func SanitizeId(id string) (string, bool) {
	regex, _ := regexp.Compile("^[0-9]+$")

	if regex.MatchString(id) {
		return id, true
	}

	return "", false
}

func SanitizeUsername(username string) (string, bool) {
	regex, _ := regexp.Compile("^[A-Za-z0-9][A-Za-z0-9_]{2,24}$")

	if regex.MatchString(username) {
		return strings.ToLower(username), true
	}

	return "", false
}
