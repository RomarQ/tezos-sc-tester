package action

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/romarq/visualtez-testing/internal/business"
)

type ModifyBlockLevelAction struct {
	raw  json.RawMessage
	json struct {
		Kind    string `json:"kind"`
		Payload struct {
			Level *int32 `json:"level"`
		} `json:"payload"`
	}
	Level int32
}

// Unmarshal action
func (action *ModifyBlockLevelAction) Unmarshal() error {
	err := json.Unmarshal(action.raw, &action.json)
	if err != nil {
		return err
	}

	// Validate action
	if err = action.validate(); err != nil {
		return err
	}

	// "level" field
	action.Level = *action.json.Payload.Level

	return nil
}

// Marshal returns the JSON of the action (cached)
func (action ModifyBlockLevelAction) Marshal() json.RawMessage {
	return action.raw
}

// Perform the action
func (action ModifyBlockLevelAction) Run(mockup business.Mockup) (interface{}, bool) {
	err := mockup.UpdateBlockLevel(action.Level)
	if err != nil {
		return err, false
	}

	return map[string]int32{
		"level": action.Level,
	}, true
}

func (action ModifyBlockLevelAction) validate() error {
	missingFields := make([]string, 0)
	if action.json.Payload.Level == nil {
		missingFields = append(missingFields, "level")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("Action of kind (%s) misses the following fields [%s].", ModifyBlockLevel, strings.Join(missingFields, ", "))
	}

	return nil
}
