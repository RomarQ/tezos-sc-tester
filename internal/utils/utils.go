package utils

import (
	"encoding/json"
	"fmt"
	"regexp"

	"blockwatch.cc/tzgo/tezos"
)

// Contains verifies if a list contains a given element
func Contains[T comparable](list []T, x T) bool {
	for _, item := range list {
		if item == x {
			return true
		}
	}
	return false
}

// GenerateKey generates a tezos wallet with Ed25519 curve
func GenerateKey() (tezos.PrivateKey, error) {
	return tezos.GenerateKey(tezos.KeyTypeEd25519)
}

// ValidateChainID validate chain_id hash
func ValidateChainID(chainID string) bool {
	var h tezos.ChainIdHash
	if err := h.UnmarshalText([]byte(chainID)); err != nil {
		return false
	}
	return true
}

// ValidateString validates string against regex expression
func ValidateString(regex string, name string) error {
	if match, err := regexp.MatchString(regex, name); !match || err != nil {
		return fmt.Errorf("String (%s) does not match pattern '%s'.", name, regex)
	}
	return nil
}

// PrettifyJSON
func PrettifyJSON(o interface{}) string {
	prettyJSON, _ := json.MarshalIndent(o, "", "  ")
	return string(prettyJSON)
}
