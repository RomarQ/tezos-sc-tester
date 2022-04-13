package michelson

import (
	"regexp"
)

var (
	regex_instruction = regexp.MustCompile("^[0-9A-Z_]+$")
	regex_identifier  = regexp.MustCompile("^[a-zA-Z0-9_]+$")
	regex_digit       = regexp.MustCompile("^-?[0-9]+$")
	regex_hex         = regexp.MustCompile("^[0-9a-fA-F]+$")
	reserved_words    = []string{
		"storage",
		"parameter",
		"code",
		"view",
	}
)

func isHex(text string) bool {
	return regex_hex.MatchString(text)
}

func isDigit(text string) bool {
	return regex_digit.MatchString(text)
}

func isIdentifier(text string) bool {
	return regex_identifier.MatchString(text)
}

func isInstruction(text string) bool {
	return regex_instruction.MatchString(text)

}
