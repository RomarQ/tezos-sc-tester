package action

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/romarq/visualtez-testing/internal/business"
	"github.com/romarq/visualtez-testing/internal/business/michelson"
	"github.com/romarq/visualtez-testing/internal/business/michelson/ast"
	MichelsonJSON "github.com/romarq/visualtez-testing/internal/business/michelson/json"
	"github.com/romarq/visualtez-testing/internal/business/michelson/micheline"
	"github.com/romarq/visualtez-testing/internal/logger"
	"github.com/romarq/visualtez-testing/internal/utils"
)

type CallContractAction struct {
	raw  json.RawMessage
	json struct {
		Kind    string `json:"kind"`
		Payload struct {
			Recipient     string          `json:"recipient"`
			Sender        string          `json:"sender"`
			Level         int32           `json:"level"`
			Timestamp     string          `json:"timestamp"`
			Entrypoint    string          `json:"entrypoint"`
			Amount        string          `json:"amount"`
			Parameter     json.RawMessage `json:"parameter"`
			ExpectFailure bool            `json:"expect_failure"`
		} `json:"payload"`
	}
	Recipient     string
	Sender        string
	Level         int32
	Timestamp     *time.Time
	Entrypoint    string
	Amount        business.Mutez
	Parameter     ast.Node
	ExpectFailure bool
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
	// "level" field
	action.Level = action.json.Payload.Level
	if action.Level == 0 {
		action.Level = 1
	}
	// "timestamp" field
	if action.json.Payload.Timestamp != "" {
		timestamp, err := utils.ParseRFC3339Timestamp(action.json.Payload.Timestamp)
		if err != nil {
			return fmt.Errorf(`field "timestamp" must use RFC3339 format. %s`, err)
		}
		action.Timestamp = &timestamp
	}
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
		return fmt.Errorf(`invalid parameter. %s`, err)
	}

	// "expect_failure" field
	action.ExpectFailure = action.json.Payload.ExpectFailure

	return nil
}

// Marshal returns the JSON of the action (cached)
func (a CallContractAction) Marshal() json.RawMessage {
	return a.raw
}

// Perform the action
func (action CallContractAction) Run(mockup business.Mockup) (interface{}, bool) {
	// Update the level of the head block
	// The transfer operation will create a new block
	err := mockup.UpdateHeadBlockLevel(action.Level - 1)
	if err != nil {
		return err, false
	}

	if action.Timestamp != nil {
		// Update the timestamp of the head block
		// Subtract one second because the next block will increment
		// the timestamp by one second
		timestamp := action.Timestamp.Add(-time.Second)
		err := mockup.UpdateHeadBlockTimestamp(utils.FormatRFC3339Timestamp(timestamp))
		if err != nil {
			return err, false
		}
	}

	parameterMicheline := replaceBigMaps(micheline.Print(action.Parameter, ""))
	parameterMicheline = expandPlaceholders(mockup, parameterMicheline)
	err = mockup.Transfer(business.CallContractArgument{
		Recipient:  action.Recipient,
		Source:     action.Sender,
		Entrypoint: action.Entrypoint,
		Amount:     action.Amount,
		Parameter:  parameterMicheline,
	})
	if err != nil {
		return err, action.ExpectFailure
	}

	storage, err := mockup.GetContractStorage(action.Recipient)
	if err != nil {
		err = fmt.Errorf("could not fetch storage for contract (%s).", action.Recipient)
		logger.Debug("[%s] %s", CallContract, err.Error())
		return err, false
	}

	actualStorageJSON, err := MichelsonJSON.Print(storage, "", "  ")
	if err != nil {
		err = fmt.Errorf("failed to print actual contract storage to JSON")
		logger.Debug("[%s] %s", AssertContractStorage, err.Error())
		return err, false
	}

	return map[string]interface{}{
		"storage": actualStorageJSON,
	}, !action.ExpectFailure
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

	if action.Level > 99999999 {
		return fmt.Errorf("The block level cannot be higher than 99999999.")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("Action of kind (%s) misses the following fields [%s].", CallContract, strings.Join(missingFields, ", "))
	}

	return nil
}
