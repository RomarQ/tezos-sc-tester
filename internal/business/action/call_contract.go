package action

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/romarq/visualtez-testing/internal/business"
	"github.com/romarq/visualtez-testing/internal/logger"
	"github.com/romarq/visualtez-testing/internal/utils"
)

type CallContractAction struct {
	Recipient  string `json:"recipient"`
	Sender     string `json:"sender"`
	Entrypoint string `json:"entrypoint"`
	Amount     string `json:"amount"`
	Parameter  string `json:"parameter"`
}

// Unmarshal action
func (action *CallContractAction) Unmarshal(bytes json.RawMessage) error {
	if err := json.Unmarshal(bytes, &action); err != nil {
		return err
	}

	// Validate action
	return action.validate()
}

// Perform the action
func (action CallContractAction) Run(mockup business.Mockup) ActionResult {
	amount, ok := new(business.TMutez).SetString(action.Amount)
	if !ok {
		errMsg := fmt.Sprintf("invalid mutez value (%s).", action.Amount)
		logger.Debug("[Task #%s] - %s", mockup.TaskID, errMsg)
		return action.buildFailureResult(errMsg)
	}

	err := mockup.Transfer(business.CallContractArgument{
		Recipient:  action.Recipient,
		Source:     action.Sender,
		Entrypoint: action.Entrypoint,
		Amount:     amount,
		Parameter:  action.Parameter,
	})
	if err != nil {
		return action.buildFailureResult(err.Error())
	}

	return action.buildSuccessResult(map[string]interface{}{})
}

func (action CallContractAction) validate() error {
	missingFields := make([]string, 0)
	if action.Recipient == "" {
		missingFields = append(missingFields, "recipient")
	} else if err := utils.ValidateString(STRING_IDENTIFIER_REGEX, action.Recipient); err != nil {
		return err
	}
	if action.Sender == "" {
		missingFields = append(missingFields, "sender")
	} else if err := utils.ValidateString(STRING_IDENTIFIER_REGEX, action.Sender); err != nil {
		return err
	}
	if action.Entrypoint == "" {
		missingFields = append(missingFields, "entrypoint")
	} else if err := utils.ValidateString(ENTRYPOINT_REGEX, action.Entrypoint); err != nil {
		return err
	}
	if action.Parameter == "" {
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
		Action: action,
	}
}

func (action CallContractAction) buildFailureResult(details string) ActionResult {
	return ActionResult{
		Status: Failure,
		Kind:   CallContract,
		Action: action,
		Result: map[string]interface{}{
			"details": details,
		},
	}
}
