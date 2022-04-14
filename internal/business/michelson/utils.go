package michelson

import (
	"regexp"
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

func isInstruction(text string) bool {
	return regex_instruction.MatchString(text)
}
