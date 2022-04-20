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
	raw  json.RawMessage
	json struct {
		Kind    string `json:"kind"`
		Payload struct {
			Recipient  string          `json:"recipient"`
			Sender     string          `json:"sender"`
			Entrypoint string          `json:"entrypoint"`
			Amount     string          `json:"amount"`
			Parameter  json.RawMessage `json:"parameter"`
		} `json:"payload"`
	}
	Recipient  string
	Sender     string
	Entrypoint string
	Amount     business.Mutez
	Parameter  ast.Node
}

// Unmarshal action
func (action *CallContractAction) Unmarshal() error {
	err := json.Unmarshal(action.raw, &action.json)
	if err != nil {
		return err
	}

	// Validate action
	if err = action.validate(); err != nil {
		return err
	}

	// "recipient" field
	action.Recipient = action.json.Payload.Recipient
	// "sender" field
	action.Sender = action.json.Payload.Sender
	// "entrypoint" field
	action.Entrypoint = action.json.Payload.Entrypoint

	// "amount" field
	action.Amount, err = business.MutezOfString(action.json.Payload.Amount)
	if err != nil {
		return err
	}

	// "parameter" field
	action.Parameter, err = michelson.ParseJSON(action.json.Payload.Parameter)
	if err != nil {
		logger.Debug("%+v", action.json.Payload.Parameter)
		return fmt.Errorf(`invalid parameter.`)
	}

	return nil
}

// Marshal returns the JSON of the action (cached)
func (a CallContractAction) Marshal() json.RawMessage {
	return a.raw
}

// Perform the action
func (action CallContractAction) Run(mockup business.Mockup) (interface{}, bool) {
	err := mockup.Transfer(business.CallContractArgument{
		Recipient:  action.Recipient,
		Source:     action.Sender,
		Entrypoint: action.Entrypoint,
		Amount:     action.Amount,
		Parameter:  expandPlaceholders(mockup, micheline.Print(action.Parameter, "")),
	})
	if err != nil {
		return err.Error(), false
	}

	return map[string]interface{}{}, true
}

func (action CallContractAction) validate() error {
	missingFields := make([]string, 0)
	if action.json.Payload.Recipient == "" {
		missingFields = append(missingFields, "recipient")
	} else if err := utils.ValidateString(STRING_IDENTIFIER_REGEX, action.json.Payload.Recipient); err != nil {
		return err
	}
	if action.json.Payload.Sender == "" {
		missingFields = append(missingFields, "sender")
	} else if err := utils.ValidateString(STRING_IDENTIFIER_REGEX, action.json.Payload.Sender); err != nil {
		return err
	}
	if action.json.Payload.Entrypoint == "" {
		missingFields = append(missingFields, "entrypoint")
	} else if err := utils.ValidateString(ENTRYPOINT_REGEX, action.json.Payload.Entrypoint); err != nil {
		return err
	}
	if action.json.Payload.Parameter == nil {
		missingFields = append(missingFields, "parameter")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("Action of kind (%s) misses the following fields [%s].", CallContract, strings.Join(missingFields, ", "))
	}

	return nil
}
