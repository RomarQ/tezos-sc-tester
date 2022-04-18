package business

import "math/big"

type (
	TTez   = big.Float
	TMutez = big.Float
)

// MutezOfTez convert tez (ꜩ) to mutez (uꜩ)
// 1 ꜩ => 1000000 uꜩ
func MutezOfTez(tez *TTez) *TMutez {
	return new(big.Float).Mul(tez, big.NewFloat(1000000))
}

// MutezOfTez convert mutez (uꜩ) to tez (ꜩ)
// 1 uꜩ => 0.000001 ꜩ
func TezOfMutez(mutez *TMutez) *TTez {
	return new(big.Float).Mul(mutez, big.NewFloat(0.000001))
}
