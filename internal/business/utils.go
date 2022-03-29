package business

import (
	TZGO "blockwatch.cc/tzgo/tezos"
)

// Generate an implicit account
func GenerateKey() (TZGO.PrivateKey, error) {
	return TZGO.GenerateKey(TZGO.KeyTypeEd25519)
}
