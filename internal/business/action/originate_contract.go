package action

import (
	"encoding/json"
	"fmt"

	"github.com/romarq/visualtez-testing/internal/business"
	"github.com/romarq/visualtez-testing/internal/logger"
)

type OriginateContractAction struct {
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
	Code    string  `json:"code"`
	Storage string  `json:"storage"`
}

// Unmarshal action
func (action *OriginateContractAction) Unmarshal(bytes json.RawMessage) error {
	if err := json.Unmarshal(bytes, &action); err != nil {
		return err
	}

	// Validate action
	return action.validate()
}

// Perform the action
func (action OriginateContractAction) Run(mockup business.Mockup) ActionResult {
	address, err := mockup.Originate("bootstrap2", action.Name, action.Balance, action.Code, string(action.Storage))
	if err != nil {
		logger.Debug("[Task #%s] - %s", mockup.TaskID, err)
		return action.buildFailureResult("Could not originate contract.")
	}
	return action.buildSuccessResult(map[string]interface{}{
		"address": address,
	})
}

func (action OriginateContractAction) validate() error {
	missingFields := make([]string, 0)
	if action.Name == "" {
		missingFields = append(missingFields, "name")
	}
	if action.Code == "" {
		missingFields = append(missingFields, "code")
	}
	if action.Storage == "" {
		missingFields = append(missingFields, "storage")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("Action of kind (%s) misses the following fields %s.", OriginateContract, missingFields)
	}

	return nil
}

func (action OriginateContractAction) buildSuccessResult(result map[string]interface{}) ActionResult {
	return ActionResult{
		Status: Success,
		Kind:   OriginateContract,
		Result: result,
		Action: action,
	}
}

func (action OriginateContractAction) buildFailureResult(details string) ActionResult {
	return ActionResult{
		Status: Failure,
		Kind:   OriginateContract,
		Action: action,
		Result: map[string]interface{}{
			"details": details,
		},
	}
}
