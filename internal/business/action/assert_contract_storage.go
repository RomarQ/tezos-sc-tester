package action

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/romarq/visualtez-testing/internal/business"
	"github.com/romarq/visualtez-testing/internal/business/michelson"
	"github.com/romarq/visualtez-testing/internal/business/michelson/ast"
	MichelsonJSON "github.com/romarq/visualtez-testing/internal/business/michelson/json"
	"github.com/romarq/visualtez-testing/internal/business/michelson/micheline"
	"github.com/romarq/visualtez-testing/internal/logger"
	"github.com/romarq/visualtez-testing/internal/utils"
)

type AssertContractStorageAction struct {
	raw  json.RawMessage
	json struct {
		Kind    string `json:"kind"`
		Payload struct {
			ContractName string          `json:"contract_name"`
			Storage      json.RawMessage `json:"storage"`
		} `json:"payload"`
	}
	ContractName string
	Storage      ast.Node
}

// Unmarshal action
func (action *AssertContractStorageAction) Unmarshal() error {
	err := json.Unmarshal(action.raw, &action.json)
	if err != nil {
		return err
	}

	// Validate action
	if err = action.validate(); err != nil {
		return err
	}

	// "contract_name" field
	action.ContractName = action.json.Payload.ContractName

	// "storage" field
	action.Storage, err = michelson.ParseJSON(action.json.Payload.Storage)
	if err != nil {
		logger.Debug("%+v", action.json.Payload.Storage)
		return fmt.Errorf(`invalid michelson.`)
	}

	return nil
}

// Marshal returns the JSON of the action (cached)
func (action AssertContractStorageAction) Marshal() json.RawMessage {
	return action.raw
}

// Perform the action
func (action AssertContractStorageAction) Run(mockup business.Mockup) (interface{}, bool) {
	storage, err := mockup.GetContractStorage(action.ContractName)
	if err != nil {
		err = fmt.Errorf("could not fetch storage for contract (%s). %s", action.ContractName, err)
		logger.Debug("[%s] %s", AssertContractStorage, err)
		return err, false
	}

	actualStorageJSON, err := MichelsonJSON.Print(storage, "", "  ")
	if err != nil {
		err = fmt.Errorf("failed to print actual contract storage to JSON. %s", err)
		logger.Debug("[%s] %s", AssertContractStorage, err)
		return err, false
	}

	// Get the storage type (the expected data needs to be normalize against the type)
	storageTypeMicheline := micheline.Print(mockup.GetCachedContract(action.ContractName).StorageType, "")

	expectedStorageMicheline := expandPlaceholders(mockup, micheline.Print(action.Storage, ""))
	expectedStorageAST, err := mockup.NormalizeData(expectedStorageMicheline, storageTypeMicheline, business.Readable)
	if err != nil {
		err = fmt.Errorf("failed to parse 'micheline'. %s", err)
		logger.Debug("[%s] %s", AssertContractStorage, err)
		return err, false
	}
	expectedStorageJSON, err := MichelsonJSON.Print(expectedStorageAST, "", "  ")
	if err != nil {
		err = fmt.Errorf("failed to print expected contract storage to JSON. %s", err)
		logger.Debug("[%s] %s", AssertContractStorage, err)
		return err, false
	}

	if expectedStorageAST != storage {
		return map[string]json.RawMessage{
			"expected": expectedStorageJSON,
			"actual":   actualStorageJSON,
		}, false
	}

	return map[string]json.RawMessage{
		"storage": actualStorageJSON,
	}, true
}

func (action AssertContractStorageAction) validate() error {
	missingFields := make([]string, 0)
	if action.json.Payload.ContractName == "" {
		missingFields = append(missingFields, "contract_name")
	} else if err := utils.ValidateString(STRING_IDENTIFIER_REGEX, action.json.Payload.ContractName); err != nil {
		return err
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("Action of kind (%s) misses the following fields [%s].", AssertContractStorage, strings.Join(missingFields, ", "))
	}

	return nil
}
