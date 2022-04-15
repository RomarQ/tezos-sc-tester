package utils

import (
	"regexp"

	"github.com/romarq/visualtez-testing/pkg/utils"
)

var (
	regex_instruction = regexp.MustCompile("^[0-9A-Z_]+$")
	RESERVED_WORDS    = []string{
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
	return utils.Contains(RESERVED_WORDS, word)
}
