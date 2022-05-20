package utils

import (
	"fmt"
	"regexp"
	"time"

	"blockwatch.cc/tzgo/tezos"
	"github.com/romarq/visualtez-testing/internal/business/michelson"
	"github.com/romarq/visualtez-testing/internal/business/michelson/ast"
)

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

// ParseRFC3339Timestamp parse RFC3339 timestamp
func ParseRFC3339Timestamp(timestamp string) (time.Time, error) {
	return time.Parse(time.RFC3339, timestamp)
}

// FormatRFC3339Timestamp format timestamp to RFC3339
func FormatRFC3339Timestamp(timestamp time.Time) string {
	return timestamp.Format(time.RFC3339)
}

// ExtractFailWithError extracts the Micheline value emitted
// by (FAILWITH) instruction
func ExtractFailWithError(output string) (ast.Node, error) {
	pattern := regexp.MustCompile("script reached FAILWITH instruction\nwith (.*)\n")

	match := pattern.FindStringSubmatch(output)
	if len(match) < 2 {
		return nil, fmt.Errorf("could not extract micheline from FAILWITH output.")
	}

	return michelson.ParseMicheline(match[1])
}
