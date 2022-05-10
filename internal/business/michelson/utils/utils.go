package utils

import (
	"regexp"
)

var (
	reserved_words = []string{
		"storage",
		"parameter",
		"code",
		"view",
	}
)

func IsInstruction(text string) bool {
	return regexp.MustCompile("^[0-9A-Z_]+$").MatchString(text)
}
func IsReservedWord(word string) bool {
	for _, item := range reserved_words {
		if item == word {
			return true
		}
	}
	return false
}
