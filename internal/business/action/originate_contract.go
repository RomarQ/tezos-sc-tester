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
	"github.com/romarq/visualtez-testing/internal/utils"
)

type OriginateContractAction struct {
	raw  json.RawMessage
	json struct {
		Kind    string `json:"kind"`
		Payload struct {
			Name    string          `json:"name"`
			Balance string          `json:"balance"`
			Code    json.RawMessage `json:"code"`
			Storage json.RawMessage `json:"storage"`
		} `json:"payload"`
	}
	Name    string
	Balance business.Mutez
	Code    ast.Node
	Storage ast.Node
}

// Unmarshal action
func (action *OriginateContractAction) Unmarshal() error {
	err := json.Unmarshal(action.raw, &action.json)
	if err != nil {
		return err
	}

	// Validate action
	if err = action.validate(); err != nil {
		return err
	}

	// "name" field
	action.Name = action.json.Payload.Name

	// "balance" field
	action.Balance, err = business.MutezOfString(action.json.Payload.Balance)
	if err != nil {
		return err
	}

	// "code" field
	action.Code, err = michelson.ParseJSON(action.json.Payload.Code)
	if err != nil {
		logger.Debug("%+v", action.json.Payload.Code)
		return fmt.Errorf(`invalid code.`)
	}

	// "storage" field
	action.Storage, err = michelson.ParseJSON(action.json.Payload.Storage)
	if err != nil {
		logger.Debug("%+v", action.json.Payload.Storage)
		return fmt.Errorf(`invalid storage.`)
	}

	return nil
}

// Marshal returns the JSON of the action (cached)
func (a OriginateContractAction) Marshal() json.RawMessage {
	return a.raw
}

// Run performs action (Originates a contract)
func (action OriginateContractAction) Run(mockup business.Mockup) (interface{}, bool) {
	if mockup.ContainsAddress(action.Name) {
		return fmt.Sprintf("Name (%s) is already in use.", action.Name), false
	}

	codeMicheline := expandPlaceholders(mockup, micheline.Print(action.Code, ""))
	storageMicheline := expandPlaceholders(mockup, micheline.Print(action.Storage, ""))
	address, err := mockup.Originate(mockup.Config.Tezos.Originator, action.Name, action.Balance, codeMicheline, storageMicheline)
	if err != nil {
		logger.Debug("[Task #%s] - %s", mockup.TaskID, err)
		return fmt.Sprintf("could not originate contract. %s", err), false
	}

	// Cache contract info
	mockup.CacheAccountAddress(action.Name, address)
	err = mockup.CacheContract(action.Name, action.Code)
	if err != nil {
		logger.Debug("[Task #%s] - %s", mockup.TaskID, err)
		return err, false
	}

	return map[string]interface{}{
		"address": address,
	}, true
}

// validate validates the action fields before interpreting them
func (action OriginateContractAction) validate() error {
	missingFields := make([]string, 0)
	if action.json.Payload.Name == "" {
		missingFields = append(missingFields, "name")
	} else if err := utils.ValidateString(STRING_IDENTIFIER_REGEX, action.json.Payload.Name); err != nil {
		return err
	}
	if action.json.Payload.Code == nil {
		missingFields = append(missingFields, "code")
	}
	if action.json.Payload.Storage == nil {
		missingFields = append(missingFields, "storage")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("Action of kind (%s) misses the following fields [%s].", OriginateContract, strings.Join(missingFields, ", "))
	}

	return nil
}
