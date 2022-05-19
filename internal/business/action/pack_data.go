package action

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/romarq/visualtez-testing/internal/business"
	"github.com/romarq/visualtez-testing/internal/business/michelson"
	"github.com/romarq/visualtez-testing/internal/business/michelson/ast"
	"github.com/romarq/visualtez-testing/internal/business/michelson/micheline"
	"github.com/romarq/visualtez-testing/internal/logger"
)

type PackDataAction struct {
	json struct {
		Kind    ActionKind `json:"kind"`
		Payload struct {
			Data json.RawMessage `json:"data"`
			Type json.RawMessage `json:"type"`
		} `json:"payload"`
	}
	Data ast.Node
	Type ast.Node
}

// Unmarshal action
func (action *PackDataAction) Unmarshal(ac Action) error {
	action.json.Kind = ac.Kind
	err := json.Unmarshal(ac.Payload, &action.json.Payload)
	if err != nil {
		return err
	}

	// Validate action
	if err = action.validate(); err != nil {
		return err
	}

	// "data" field
	action.Data, err = michelson.ParseJSON(action.json.Payload.Data)
	if err != nil {
		logger.Debug("%+v", action.json.Payload.Data)
		return fmt.Errorf(`invalid michelson value.`)
	}

	// "type" field
	action.Type, err = michelson.ParseJSON(action.json.Payload.Type)
	if err != nil {
		logger.Debug("%+v", action.json.Payload.Type)
		return fmt.Errorf(`invalid michelson type.`)
	}

	return nil
}

// Marshal returns the JSON of the action (cached)
func (action PackDataAction) Action() interface{} {
	return action.json
}

// Run performs action (Serializes a michelson value)
func (action PackDataAction) Run(mockup business.Mockup) (interface{}, bool) {
	dataMicheline := expandPlaceholders(mockup, micheline.Print(action.Data, ""))
	typeMicheline := expandPlaceholders(mockup, micheline.Print(action.Type, ""))

	byteString, err := mockup.SerializeData(dataMicheline, typeMicheline)
	if err != nil {
		logger.Debug("[Task #%s] - %s", mockup.TaskID, err)
		return fmt.Sprintf("could not serialize michelson data. %s", err), false
	}

	return map[string]string{
		"bytes": byteString,
	}, true
}

// validate validates the action fields before interpreting them
func (action PackDataAction) validate() error {
	missingFields := make([]string, 0)
	if action.json.Payload.Data == nil {
		missingFields = append(missingFields, "code")
	}
	if action.json.Payload.Type == nil {
		missingFields = append(missingFields, "storage")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("Action of kind (%s) misses the following fields [%s].", PackData, strings.Join(missingFields, ", "))
	}

	return nil
}
