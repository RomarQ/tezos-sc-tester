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
	json struct {
		Name    string          `json:"name"`
		Balance string          `json:"balance"`
		Code    json.RawMessage `json:"code"`
		Storage json.RawMessage `json:"storage"`
	}
	Name    string
	Balance business.Mutez
	Code    ast.Node
	Storage ast.Node
}

// Unmarshal action
func (action *OriginateContractAction) Unmarshal(bytes json.RawMessage) error {
	err := json.Unmarshal(bytes, &action.json)
	if err != nil {
		return err
	}

	// Validate action
	if err = action.validate(); err != nil {
		return err
	}

	// "name" field
	action.Name = action.json.Name

	// "balance" field
	action.Balance, err = business.MutezOfString(action.json.Balance)
	if err != nil {
		return err
	}

	// "code" field
	action.Code, err = michelson.ParseJSON(action.json.Code)
	if err != nil {
		logger.Debug("%+v", action.json.Code)
		return fmt.Errorf(`invalid code.`)
	}

	// "storage" field
	action.Storage, err = michelson.ParseJSON(action.json.Storage)
	if err != nil {
		logger.Debug("%+v", action.json.Storage)
		return fmt.Errorf(`invalid storage.`)
	}

	return nil
}

// Perform action (Originates a contract)
func (action OriginateContractAction) Run(mockup business.Mockup) ActionResult {
	if mockup.ContainsAddress(action.Name) {
		return action.buildFailureResult(fmt.Sprintf("Name (%s) is already in use.", action.Name))
	}

	codeMicheline := micheline.Print(action.Code, "")
	storageMicheline := micheline.Print(action.Storage, "")
	address, err := mockup.Originate(mockup.Config.Tezos.Originator, action.Name, action.Balance, codeMicheline, storageMicheline)
	if err != nil {
		logger.Debug("[Task #%s] - %s", mockup.TaskID, err)
		return action.buildFailureResult(fmt.Sprintf("could not originate contract. %s", err))
	}

	// Save new address
	mockup.SetAddress(action.Name, address)

	return action.buildSuccessResult(map[string]interface{}{
		"address": address,
	})
}

func (action OriginateContractAction) validate() error {
	missingFields := make([]string, 0)
	if action.json.Name == "" {
		missingFields = append(missingFields, "name")
	} else if err := utils.ValidateString(STRING_IDENTIFIER_REGEX, action.json.Name); err != nil {
		return err
	}
	if action.json.Code == nil {
		missingFields = append(missingFields, "code")
	}
	if action.json.Storage == nil {
		missingFields = append(missingFields, "storage")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("Action of kind (%s) misses the following fields [%s].", OriginateContract, strings.Join(missingFields, ", "))
	}

	return nil
}

func (action OriginateContractAction) buildSuccessResult(result map[string]interface{}) ActionResult {
	return ActionResult{
		Status: Success,
		Kind:   OriginateContract,
		Result: result,
		Action: action.json,
	}
}

func (action OriginateContractAction) buildFailureResult(details string) ActionResult {
	return ActionResult{
		Status: Failure,
		Kind:   OriginateContract,
		Action: action.json,
		Result: map[string]interface{}{
			"details": details,
		},
	}
}
