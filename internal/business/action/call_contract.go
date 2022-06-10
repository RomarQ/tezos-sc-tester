package action

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/romarq/tezos-sc-tester/internal/business"
	"github.com/romarq/tezos-sc-tester/internal/business/michelson"
	"github.com/romarq/tezos-sc-tester/internal/business/michelson/ast"
	MichelsonJSON "github.com/romarq/tezos-sc-tester/internal/business/michelson/json"
	"github.com/romarq/tezos-sc-tester/internal/business/michelson/micheline"
	"github.com/romarq/tezos-sc-tester/internal/logger"
	"github.com/romarq/tezos-sc-tester/internal/utils"
)

type CallContractAction struct {
	json struct {
		Kind    ActionKind `json:"kind"`
		Payload struct {
			Recipient      string          `json:"recipient"`
			Sender         string          `json:"sender"`
			Entrypoint     string          `json:"entrypoint"`
			Amount         string          `json:"amount"`
			Parameter      json.RawMessage `json:"parameter"`
			ExpectFailwith json.RawMessage `json:"expect_failwith,omitempty"`
		} `json:"payload"`
	}
	Recipient      string
	Sender         string
	Entrypoint     string
	Amount         business.Mutez
	Parameter      ast.Node
	ExpectFailwith ast.Node
}

// Unmarshal action
func (action *CallContractAction) Unmarshal(ac Action) error {
	action.json.Kind = ac.Kind
	err := json.Unmarshal(ac.Payload, &action.json.Payload)
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
		return fmt.Errorf("invalid 'parameter'. %s", err)
	}

	// "expect_failwith" field
	if action.json.Payload.ExpectFailwith != nil {
		action.ExpectFailwith, err = michelson.ParseJSON(action.json.Payload.ExpectFailwith)
		if err != nil {
			logger.Debug("%+v", action.json.Payload.ExpectFailwith)
			return fmt.Errorf("invalid 'expect_failwith'. %s", err)
		}
	}

	return nil
}

// Marshal returns the JSON of the action (cached)
func (action CallContractAction) Action() interface{} {
	return action.json
}

// Perform the action
func (action CallContractAction) Run(mockup business.Mockup) (interface{}, bool) {
	parameterMicheline := replaceBigMaps(micheline.Print(action.Parameter, ""))
	parameterMicheline = expandPlaceholders(mockup, parameterMicheline)
	err := mockup.Transfer(business.CallContractArgument{
		Recipient:  action.Recipient,
		Source:     action.Sender,
		Entrypoint: action.Entrypoint,
		Amount:     action.Amount,
		Parameter:  parameterMicheline,
	})
	if err != nil {
		if action.ExpectFailwith == nil {
			// The transfer was not expected to fail
			return err, false
		}
		// The transfer was expected to fail.
		// Extract the error emitted with (FAILWITH), the error is a micheline value
		michelineError, err := utils.ExtractFailWithError(err.Error())
		if err != nil {
			return err, false
		}

		// Validate the error against the user input
		if michelineError.String() != action.ExpectFailwith.String() {
			michelsonJson, err := MichelsonJSON.Print(michelineError, "", "  ")
			if err != nil {
				errMsg := fmt.Sprintf("failed to print (FAILWITH) result to michelson JSON. %s", err.Error())
				logger.Debug("[%s] %s", AssertContractStorage, errMsg)
			}
			return map[string]json.RawMessage{
				"expected": action.json.Payload.ExpectFailwith,
				"actual":   michelsonJson,
			}, false
		}
	}

	storage, err := mockup.GetContractStorage(action.Recipient)
	if err != nil {
		logger.Debug("[%s] %s", CallContract, err.Error())
		return fmt.Errorf("could not fetch storage for contract (%s).", action.Recipient), false
	}

	actualStorageJSON, err := MichelsonJSON.Print(storage, "", "  ")
	if err != nil {
		logger.Debug("[%s] %s", AssertContractStorage, err.Error())
		return fmt.Errorf("failed to print actual contract storage to JSON"), false
	}

	return map[string]interface{}{
		"storage": actualStorageJSON,
	}, true
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
