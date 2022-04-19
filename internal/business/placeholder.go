package business

import (
	"bytes"
	"fmt"
	"regexp"
)

var (
	PLACEHOLDER__ADDRESS_OF_ACCOUNT = "TEST__ADDRESS_OF_ACCOUNT__"
)

func ExpandAccountPlaceholders(addresses map[string]string, b []byte) []byte {
	address_of_account_regex := regexp.MustCompile(fmt.Sprintf("%s([a-zA-Z0-9_]+)", PLACEHOLDER__ADDRESS_OF_ACCOUNT))

	placeholders := address_of_account_regex.FindAll(b, -1)
	for _, placeholder := range placeholders {
		accountID := bytes.Replace(placeholder, []byte(PLACEHOLDER__ADDRESS_OF_ACCOUNT), []byte{}, 1)

		b = bytes.ReplaceAll(b, placeholder, []byte(addresses[string(accountID)]))
	}

	return b
}
