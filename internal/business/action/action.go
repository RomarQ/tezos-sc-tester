package action

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/romarq/visualtez-testing/internal/business"
	Error "github.com/romarq/visualtez-testing/internal/error"
)

type (
	ActionStatus string
	ActionResult struct {
		Status ActionStatus `json:"status"`
		Action interface{}  `json:"action"`
		Result interface{}  `json:"result,omitempty"`
	}
	Action struct {
		Kind    ActionKind      `json:"kind"`
		Payload json.RawMessage `json:"payload"`
	}
	IAction interface {
		Run(mockup business.Mockup) (interface{}, bool)
		Unmarshal(action Action) error
		Action() interface{}
	}
)

const (
	Failure ActionStatus = "failure"
	Success ActionStatus = "success"
)

const (
	STRING_IDENTIFIER_REGEX = "^[a-zA-Z0-9_]+$"
	ENTRYPOINT_REGEX        = "^[a-zA-Z0-9_]{1,31}$"
)

// GetActions unmarshal test actions
func GetActions(rawActions []Action) ([]IAction, error) {
	actions := make([]IAction, 0)
	for _, rawAction := range rawActions {
		var action IAction

		switch rawAction.Kind {
		default:
			return nil, fmt.Errorf("Unexpected action kind (%s).", rawAction.Kind)
		case AssertAccountBalance:
			action = &AssertAccountBalanceAction{}
		case AssertContractStorage:
			action = &AssertContractStorageAction{}
		case CallContract:
			action = &CallContractAction{}
		case OriginateContract:
			action = &OriginateContractAction{}
		case CreateImplicitAccount:
			action = &CreateImplicitAccountAction{}
		case ModifyChainID:
			action = &ModifyChainIdAction{}
		case PackData:
			action = &PackDataAction{}
		}

		if err := action.Unmarshal(rawAction); err != nil {
			return nil, Error.DetailedHttpError(http.StatusBadRequest, err.Error(), rawAction)
		}
		actions = append(actions, action)
	}

	return actions, nil
}

// ApplyActions executes each test action
func ApplyActions(mockup business.Mockup, actions []IAction) []ActionResult {
	responses := make([]ActionResult, 0)

	for _, action := range actions {
		result, ok := action.Run(mockup)
		if ok {
			responses = append(responses, buildResult(Success, result, action))
		} else {
			responses = append(responses, buildResult(Failure, result, action))
		}
	}

	return responses
}

func expandPlaceholders(mockup business.Mockup, str string) string {
	// Expand addresses
	b := business.ExpandAccountPlaceholders(mockup.Addresses, []byte(str))
	// Expand balances
	b = business.ExpandBalancePlaceholders(mockup, b)

	return string(b)
}

// replaceBigMaps converts all 'big_map' types to 'map'
// This is necessary for testing the storage updates
func replaceBigMaps(str string) string {
	return strings.ReplaceAll(str, "big_map", "map")
}

func buildResult(status ActionStatus, result interface{}, action IAction) ActionResult {
	switch v := result.(type) {
	case string:
		result = map[string]interface{}{
			"details": v,
		}
	case error:
		result = map[string]interface{}{
			"details": v.Error(),
		}
	}

	return ActionResult{
		Status: status,
		Action: action.Action(),
		Result: result,
	}
}
