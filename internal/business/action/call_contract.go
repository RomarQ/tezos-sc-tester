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

type CallContractAction struct {
	json struct {
		Recipient  string          `json:"recipient"`
		Sender     string          `json:"sender"`
		Entrypoint string          `json:"entrypoint"`
		Amount     string          `json:"amount"`
		Parameter  json.RawMessage `json:"parameter"`
	}
	Recipient  string
	Sender     string
	Entrypoint string
	Amount     business.Mutez
	Parameter  ast.Node
}

// Unmarshal action
func (action *CallContractAction) Unmarshal(bytes json.RawMessage) error {
	err := json.Unmarshal(bytes, &action.json)
	if err != nil {
		return err
	}

	// Validate action
	if err = action.validate(); err != nil {
		return err
	}

	// "recipient" field
	action.Recipient = action.json.Recipient
	// "sender" field
	action.Sender = action.json.Sender
	// "entrypoint" field
	action.Entrypoint = action.json.Entrypoint

	// "amount" field
	action.Amount, err = business.MutezOfString(action.json.Amount)
	if err != nil {
		return err
	}

	// "parameter" field
	action.Parameter, err = michelson.ParseJSON(action.json.Parameter)
	if err != nil {
		logger.Debug("%+v", action.json.Parameter)
		return fmt.Errorf(`invalid parameter.`)
	}

	return nil
}

// Perform the action
func (action CallContractAction) Run(mockup business.Mockup) ActionResult {

	err := mockup.Transfer(business.CallContractArgument{
		Recipient:  action.Recipient,
		Source:     action.Sender,
		Entrypoint: action.Entrypoint,
		Amount:     action.Amount,
		Parameter:  micheline.Print(action.Parameter, ""),
	})
	if err != nil {
		return action.buildFailureResult(err.Error())
	}

	return action.buildSuccessResult(map[string]interface{}{})
}

func (action CallContractAction) validate() error {
	missingFields := make([]string, 0)
	if action.json.Recipient == "" {
		missingFields = append(missingFields, "recipient")
	} else if err := utils.ValidateString(STRING_IDENTIFIER_REGEX, action.json.Recipient); err != nil {
		return err
	}
	if action.json.Sender == "" {
		missingFields = append(missingFields, "sender")
	} else if err := utils.ValidateString(STRING_IDENTIFIER_REGEX, action.json.Sender); err != nil {
		return err
	}
	if action.json.Entrypoint == "" {
		missingFields = append(missingFields, "entrypoint")
	} else if err := utils.ValidateString(ENTRYPOINT_REGEX, action.json.Entrypoint); err != nil {
		return err
	}
	if action.json.Parameter == nil {
		missingFields = append(missingFields, "parameter")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("Action of kind (%s) misses the following fields [%s].", CallContract, strings.Join(missingFields, ", "))
	}

	return nil
}

func (action CallContractAction) buildSuccessResult(result map[string]interface{}) ActionResult {
	return ActionResult{
		Status: Success,
		Kind:   CallContract,
		Result: result,
		Action: action.json,
	}
}

func (action CallContractAction) buildFailureResult(details string) ActionResult {
	return ActionResult{
		Status: Failure,
		Kind:   CallContract,
		Action: action.json,
		Result: map[string]interface{}{
			"details": details,
		},
	}
}
