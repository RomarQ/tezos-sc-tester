package business

import (
	"fmt"
	"regexp"

	TZGO "blockwatch.cc/tzgo/tezos"
)

// Generate an implicit account
func GenerateKey() (TZGO.PrivateKey, error) {
	return TZGO.GenerateKey(TZGO.KeyTypeEd25519)
}

// ValidateString validates string against regex expression
func ValidateString(regex string, name string) error {
	if match, err := regexp.MatchString(regex, name); !match || err != nil {
		return fmt.Errorf("String (%s) does not match pattern '%s'.", name, regex)
	}

	return nil
}
