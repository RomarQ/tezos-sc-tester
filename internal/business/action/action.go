package action

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tidwall/gjson"

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
		Kind    ActionKind  `json:"kind"`
		Payload interface{} `json:"payload"`
	}
	IAction interface {
		Run(mockup business.Mockup) (interface{}, bool)
		Unmarshal() error
		Marshal() json.RawMessage
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
func GetActions(body io.ReadCloser) ([]IAction, error) {
	rawActions := make([]json.RawMessage, 0)

	err := json.NewDecoder(body).Decode(&rawActions)
	if err != nil {
		return nil, err
	}

	actions := make([]IAction, 0)
	for _, rawAction := range rawActions {
		var action IAction

		kind := gjson.GetBytes(rawAction, "kind")
		switch kind.String() {
		default:
			return nil, fmt.Errorf("Unexpected action kind (%s).", kind)
		case string(CallContract):
			action = &CallContractAction{
				raw: rawAction,
			}
		case string(OriginateContract):
			action = &OriginateContractAction{
				raw: rawAction,
			}
		case string(CreateImplicitAccount):
			action = &CreateImplicitAccountAction{
				raw: rawAction,
			}
		}

		if err = action.Unmarshal(); err != nil {
			return nil, Error.DetailedHttpError(http.StatusBadRequest, err.Error(), rawAction)
		}
		actions = append(actions, action)
	}

	return actions, err
}

// ApplyActions executes each test action
func ApplyActions(mockup business.Mockup, actions []IAction) []ActionResult {
	responses := make([]ActionResult, 0)

	for _, action := range actions {
		result, ok := action.Run(mockup)
		if ok {
			responses = append(responses, buildSuccessResult(result, action))
		} else {
			responses = append(responses, buildFailureResult(result, action))

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

func buildSuccessResult(result interface{}, action IAction) ActionResult {
	return ActionResult{
		Status: Success,
		Action: action.Marshal(),
		Result: result,
	}
}

func buildFailureResult(details interface{}, action IAction) ActionResult {
	return ActionResult{
		Status: Failure,
		Action: action.Marshal(),
		Result: map[string]interface{}{
			"details": details,
		},
	}
}
