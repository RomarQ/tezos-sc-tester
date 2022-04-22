package action

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/romarq/visualtez-testing/internal/business"
	"github.com/romarq/visualtez-testing/internal/utils"
)

type ModifyChainIdAction struct {
	raw  json.RawMessage
	json struct {
		Kind    string `json:"kind"`
		Payload struct {
			ChainID string `json:"chain_id"`
		} `json:"payload"`
	}
	ChainID string
}

// Unmarshal action
func (action *ModifyChainIdAction) Unmarshal() error {
	err := json.Unmarshal(action.raw, &action.json)
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
func (action ModifyChainIdAction) Marshal() json.RawMessage {
	return action.raw
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
	} else if err := utils.ValidateString(STRING_IDENTIFIER_REGEX, action.json.Payload.ChainID); err != nil {
		return err
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("Action of kind (%s) misses the following fields [%s].", ModifyChainID, strings.Join(missingFields, ", "))
	}

	return nil
}
