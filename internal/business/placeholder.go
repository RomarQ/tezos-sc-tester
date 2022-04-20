package business

import (
	"bytes"
	"fmt"
	"regexp"
)

var (
	PLACEHOLDER__ADDRESS_OF_ACCOUNT = "TEST__ADDRESS_OF_ACCOUNT__"
	PLACEHOLDER__BALANCE_OF_ACCOUNT = "TEST__BALANCE_OF_ACCOUNT__"
)

// ExpandAccountPlaceholders expands the real account address from a placeholder that identifies the account
func ExpandAccountPlaceholders(addresses map[string]string, b []byte) []byte {
	regex := regexp.MustCompile(fmt.Sprintf("%s([a-zA-Z0-9_]+)", PLACEHOLDER__ADDRESS_OF_ACCOUNT))

	placeholders := regex.FindAll(b, -1)
	for _, placeholder := range placeholders {
		accountID := bytes.Replace(placeholder, []byte(PLACEHOLDER__ADDRESS_OF_ACCOUNT), []byte{}, 1)

		b = bytes.ReplaceAll(b, placeholder, []byte(addresses[string(accountID)]))
	}

	return b
}

// ExpandBalancePlaceholders expands the account balance from a placeholder that identifies the account
func ExpandBalancePlaceholders(mockup Mockup, b []byte) []byte {
	regex := regexp.MustCompile(fmt.Sprintf("%s([a-zA-Z0-9_]+)", PLACEHOLDER__BALANCE_OF_ACCOUNT))

	placeholders := regex.FindAll(b, -1)
	for _, placeholder := range placeholders {
		accountID := bytes.Replace(placeholder, []byte(PLACEHOLDER__BALANCE_OF_ACCOUNT), []byte{}, 1)
		balance := mockup.GetBalance(string(accountID))
		b = bytes.ReplaceAll(b, placeholder, []byte(balance.String()))
	}

	return b
}
