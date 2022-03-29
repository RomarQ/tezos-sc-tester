package action

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/tidwall/gjson"

	"github.com/romarq/visualtez-testing/internal/business"
)

type (
	TestStatus string
	TestResult struct {
		Status      TestStatus  `json:"status"`
		Kind        ActionKind  `json:"kind"`
		Description string      `json:"description,omitempty"`
		Action      interface{} `json:"action"`
	}
	Action struct {
		Kind    ActionKind  `json:"kind"`
		Payload interface{} `json:"payload"`
	}
	IAction interface {
		Run(mockup business.Mockup) error
		Unmarshal(bytes json.RawMessage) error
	}
)

const (
	Failure TestStatus = "failure"
	Success            = "success"
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

func ApplyActions(mockup business.Mockup, actions []IAction) []TestResult {
	getSuccessResponse := func(kind ActionKind, action interface{}) TestResult {
		return TestResult{
			Status: Success,
			Kind:   kind,
			Action: action,
		}
	}

	getFailureResponse := func(kind ActionKind, description string, action interface{}) TestResult {
		return TestResult{
			Status:      Failure,
			Kind:        kind,
			Description: description,
			Action:      action,
		}
	}

	responses := make([]TestResult, 0)
	for _, action := range actions {
		if err := action.Run(mockup); err != nil {
			responses = append(
				responses,
				getFailureResponse(CreateImplicitAccount, err.Error(), action),
			)
		} else {
			responses = append(
				responses,
				getSuccessResponse(CreateImplicitAccount, action),
			)
		}
	}

	return responses
}
