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
		Status ActionStatus           `json:"status"`
		Kind   ActionKind             `json:"kind"`
		Action interface{}            `json:"action"`
		Result map[string]interface{} `json:"result,omitempty"`
	}
	Action struct {
		Kind    ActionKind  `json:"kind"`
		Payload interface{} `json:"payload"`
	}
	IAction interface {
		Run(mockup business.Mockup) ActionResult
		Unmarshal(bytes json.RawMessage) error
	}
)

const (
	Failure ActionStatus = "failure"
	Success              = "success"
)

const (
	STRING_IDENTIFIER_REGEX = "^[a-zA-Z0-9._-]+$"
	ENTRYPOINT_REGEX        = "^[a-zA-Z0-9_]{1,31}$"
)

// Unmarshal actions
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
			action = &CallContractAction{}
		case string(OriginateContract):
			action = &OriginateContractAction{}
		case string(CreateImplicitAccount):
			action = &CreateImplicitAccountAction{}
		}

		payload := gjson.GetBytes(rawAction, "payload")
		if err = action.Unmarshal(json.RawMessage(payload.Raw)); err != nil {
			return nil, Error.DetailedHttpError(http.StatusBadRequest, err.Error(), rawAction)
		}
		actions = append(actions, action)
	}

	return actions, err
}

func ApplyActions(mockup business.Mockup, actions []IAction) []ActionResult {
	responses := make([]ActionResult, 0)

	for _, action := range actions {
		responses = append(responses, action.Run(mockup))
	}

	return responses
}
