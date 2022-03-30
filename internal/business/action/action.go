package action

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/tidwall/gjson"

	"github.com/romarq/visualtez-testing/internal/business"
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

// Unmarshal actions
func GetActions(body io.ReadCloser) ([]IAction, error) {
	rawActions := make([]json.RawMessage, 0)

	err := json.NewDecoder(body).Decode(&rawActions)
	if err != nil {
		return nil, err
	}

	actions := make([]IAction, 0)
	for _, rawAction := range rawActions {
		kind := gjson.GetBytes(rawAction, `kind`)
		payload := gjson.GetBytes(rawAction, `payload`)
		switch kind.String() {
		default:
			return nil, fmt.Errorf("Unexpected action kind (%s).", kind)
		case string(OriginateContract):
			action := &OriginateContractAction{}
			if err = action.Unmarshal(json.RawMessage(payload.Raw)); err != nil {
				return nil, err
			}
			actions = append(actions, action)
		case string(CreateImplicitAccount):
			action := &CreateImplicitAccountAction{}
			if err = action.Unmarshal(json.RawMessage(payload.Raw)); err != nil {
				return nil, err
			}
			actions = append(actions, action)
		}
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
