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

type AssertContractStorageAction struct {
	json struct {
		Kind    ActionKind `json:"kind"`
		Payload struct {
			ContractName string          `json:"contract_name"`
			Storage      json.RawMessage `json:"storage"`
		} `json:"payload"`
	}
	ContractName string
	Storage      ast.Node
}

// Unmarshal action
func (action *AssertContractStorageAction) Unmarshal(ac Action) error {
	action.json.Kind = ac.Kind
	err := json.Unmarshal(ac.Payload, &action.json.Payload)
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
		return fmt.Errorf(`invalid michelson. %s`, err)
	}

	return nil
}

// Marshal returns the JSON of the action (cached)
func (action AssertContractStorageAction) Action() interface{} {
	return action.json
}

// Perform the action
func (action AssertContractStorageAction) Run(mockup business.Mockup) (result interface{}, success bool) {
	defer func() {
		if err := recover(); err != nil {
			result = err
			success = false
		}
	}()

	// Get current storage (already normalized)
	storage, err := mockup.GetContractStorage(action.ContractName)
	if err != nil {
		errMsg := fmt.Errorf("could not fetch storage for contract (%s)", action.ContractName)
		logger.Debug("[%s] %s. %s", AssertContractStorage, errMsg, err)
		return errMsg, false
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

	if expectedStorageAST.String() != storage.String() {
		return map[string]json.RawMessage{
			"expected": expectedStorageJSON,
			"actual":   actualStorageJSON,
		}, false
	}

	return map[string]json.RawMessage{
		"storage": actualStorageJSON,
	}, true
}

// validate validates the action fields before interpreting them
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
