package action

import (
	"encoding/json"
	"fmt"

	TZGO "blockwatch.cc/tzgo/tezos"
)

type CreateImplicitAccountActionPayload struct {
	Name    string `json:"name"`
	Balance int64  `json:"balance"`
}

type CreateImplicitAccountAction struct {
	Kind    ActionKind                         `json:"kind"`
	Payload CreateImplicitAccountActionPayload `json:"payload"`
}

func (action *CreateImplicitAccountAction) validate() error {
	if action.Payload.Name == "" {
		return fmt.Errorf("Actions of kind (%s) must contain a name field.", CreateImplicitAccount)
	}
	return nil
}

// Unmarshal action
func (action *CreateImplicitAccountAction) Unmarshal(bytes json.RawMessage) error {
	if err := json.Unmarshal(bytes, &action); err != nil {
		return err
	}

	// Validate action
	return action.validate()
}

// Generate an implicit account
func (action *CreateImplicitAccountAction) GenerateKey() (TZGO.PrivateKey, error) {
	keyPair, err := TZGO.GenerateKey(TZGO.KeyTypeEd25519)

	return keyPair, err
}
