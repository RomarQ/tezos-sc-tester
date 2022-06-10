package action

import (
	"encoding/json"
	"fmt"

	"github.com/romarq/tezos-sc-tester/internal/business"
)

type ModifyBlockLevelAction struct {
	json struct {
		Kind    ActionKind `json:"kind"`
		Payload struct {
			Level int32 `json:"level"`
		} `json:"payload"`
	}
	Level int32
}

// Unmarshal action
func (action *ModifyBlockLevelAction) Unmarshal(ac Action) error {
	action.json.Kind = ac.Kind
	err := json.Unmarshal(ac.Payload, &action.json.Payload)
	if err != nil {
		return err
	}

	// Validate action
	if err = action.validate(); err != nil {
		return err
	}

	// "level" field
	action.Level = action.json.Payload.Level
	if action.Level == 0 {
		action.Level = 1
	}

	return nil
}

// Marshal returns the JSON of the action (cached)
func (action ModifyBlockLevelAction) Action() interface{} {
	return action.json
}

// Perform the action
func (action ModifyBlockLevelAction) Run(mockup business.Mockup) (interface{}, bool) {
	// Update the level of the head block
	// The transfer operation will create a new block
	err := mockup.UpdateHeadBlockLevel(action.Level - 1)
	if err != nil {
		return err, false
	}

	return map[string]string{}, true
}

func (action ModifyBlockLevelAction) validate() error {
	if action.json.Payload.Level < 1 {
		return fmt.Errorf("The block level must be higher than 0.")
	}
	if action.json.Payload.Level > 99999999 {
		return fmt.Errorf("The block level cannot be higher than 99999999.")
	}

	return nil
}
