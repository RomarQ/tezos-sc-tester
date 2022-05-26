package action

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/romarq/tezos-sc-tester/internal/business"
	"github.com/romarq/tezos-sc-tester/internal/utils"
)

type ModifyChainIdAction struct {
	json struct {
		Kind    ActionKind `json:"kind"`
		Payload struct {
			ChainID string `json:"chain_id"`
		} `json:"payload"`
	}
	ChainID string
}

// Unmarshal action
func (action *ModifyChainIdAction) Unmarshal(ac Action) error {
	action.json.Kind = ac.Kind
	err := json.Unmarshal(ac.Payload, &action.json.Payload)
	if err != nil {
		return err
	}

	// Validate action
	if err = action.validate(); err != nil {
		return err
	}

	// "chain_id" field
	action.ChainID = action.json.Payload.ChainID

	return nil
}

// Marshal returns the JSON of the action (cached)
func (action ModifyChainIdAction) Action() interface{} {
	return action.json
}

// Perform the action
func (action ModifyChainIdAction) Run(mockup business.Mockup) (interface{}, bool) {
	err := mockup.UpdateChainID(action.ChainID)
	if err != nil {
		return err, false
	}

	return map[string]string{
		"chain_id": action.ChainID,
	}, true
}

func (action ModifyChainIdAction) validate() error {
	missingFields := make([]string, 0)
	if action.json.Payload.ChainID == "" {
		missingFields = append(missingFields, "chain_id")
	} else if !utils.ValidateChainID(action.json.Payload.ChainID) {
		return fmt.Errorf(`"chain_id" is invalid: %s`, action.json.Payload.ChainID)
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("Action of kind (%s) misses the following fields [%s].", ModifyChainID, strings.Join(missingFields, ", "))
	}

	return nil
}
