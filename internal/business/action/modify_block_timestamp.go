package action

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/romarq/tezos-sc-tester/internal/business"
	"github.com/romarq/tezos-sc-tester/internal/utils"
)

type ModifyBlockTimestampAction struct {
	json struct {
		Kind    ActionKind `json:"kind"`
		Payload struct {
			Timestamp string `json:"timestamp"`
		} `json:"payload"`
	}
	Timestamp time.Time
}

// Unmarshal action
func (action *ModifyBlockTimestampAction) Unmarshal(ac Action) error {
	action.json.Kind = ac.Kind
	err := json.Unmarshal(ac.Payload, &action.json.Payload)
	if err != nil {
		return err
	}

	// Validate action
	if err = action.validate(); err != nil {
		return err
	}

	// "timestamp" field
	timestamp, err := utils.ParseRFC3339Timestamp(action.json.Payload.Timestamp)
	if err != nil {
		return fmt.Errorf("field 'timestamp' must use RFC3339 format. %s", err)
	}
	action.Timestamp = timestamp

	return nil
}

// Marshal returns the JSON of the action (cached)
func (action ModifyBlockTimestampAction) Action() interface{} {
	return action.json
}

// Perform the action
func (action ModifyBlockTimestampAction) Run(mockup business.Mockup) (interface{}, bool) {
	// Update the timestamp of the head block
	// Subtract one second because the next block will increment
	// the timestamp by one second
	timestamp := action.Timestamp.Add(-time.Second)
	err := mockup.UpdateHeadBlockTimestamp(utils.FormatRFC3339Timestamp(timestamp))
	if err != nil {
		return err, false
	}

	return map[string]string{}, true
}

func (action ModifyBlockTimestampAction) validate() error {
	missingFields := make([]string, 0)
	if action.json.Payload.Timestamp == "" {
		missingFields = append(missingFields, "timestamp")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("Action of kind (%s) misses the following fields [%s].", ModifyChainID, strings.Join(missingFields, ", "))
	}

	return nil
}
