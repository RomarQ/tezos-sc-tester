package business

import (
	"fmt"
	"math/big"
)

type (
	Tez struct {
		v *big.Float
	}
	Mutez struct {
		v *big.Float
	}
)

// TezOf construct a Tez value from big.Float
func TezOfFloat(value *big.Float) Tez {
	return Tez{
		v: value,
	}
}

// TezOfString
func TezOfString(value string) (t Tez, err error) {
	v, ok := new(big.Float).SetString(value)
	if !ok {
		err = fmt.Errorf(`invalid tez value: %s.`, value)
	}
	t.v = v
	return
}

// MutezOf construct a Mutez value from big.Float
func MutezOfFloat(value *big.Float) Mutez {
	return Mutez{
		v: value,
	}
}

// MutezOfString
func MutezOfString(value string) (t Mutez, err error) {
	v, ok := new(big.Float).SetString(value)
	if !ok {
		err = fmt.Errorf(`invalid mutez value: %s.`, value)
	}
	t.v = v
	return
}

// AddMutez adds two Mutez values and returns a the result
func AddMutez(m1 Mutez, m2 Mutez) Mutez {
	return MutezOfFloat(new(big.Float).Add(m1.v, m2.v))
}

// ToTez convert mutez (uꜩ) to tez (ꜩ)
// 1 uꜩ => 0.000001 ꜩ
func (m Mutez) ToTez() Tez {
	return TezOfFloat(new(big.Float).Mul(m.v, big.NewFloat(0.000001)))
}

// String stringify a value of type Mutez
func (m Mutez) String() string {
	return m.v.String()
}

// ToMutez convert tez (ꜩ) to mutez (uꜩ)
// 1 ꜩ => 1000000 uꜩ
func (t Tez) ToMutez() Mutez {
	return MutezOfFloat(new(big.Float).Mul(t.v, big.NewFloat(1000000)))
}

// String stringify a value of type Tez
func (t Tez) String() string {
	return t.v.Text('f', 6)
}
