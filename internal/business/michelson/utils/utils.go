package utils

import (
	"regexp"

	"github.com/romarq/visualtez-testing/internal/utils"
)

var (
	regex_instruction = regexp.MustCompile("^[0-9A-Z_]+$")
	reserved_words    = []string{
		"storage",
		"parameter",
		"code",
		"view",
	}
)

func IsInstruction(text string) bool {
	return regex_instruction.MatchString(text)
}
func IsReservedWord(word string) bool {
	return utils.Contains(reserved_words, word)
}
