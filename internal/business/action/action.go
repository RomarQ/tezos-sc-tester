package action

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/tidwall/gjson"

	Mockup "github.com/romarq/visualtez-testing/internal/business"
)

type TestStatus string

const (
	Failure TestStatus = "failure"
	Success            = "success"
)

type TestResult struct {
	Status      TestStatus  `json:"status"`
	Description string      `json:"description,omitempty"`
	Action      interface{} `json:"action"`
}

// Unmarshal actions
func GetActions(body io.ReadCloser) ([]interface{}, error) {
	rawActions := make([]json.RawMessage, 0)

	err := json.NewDecoder(body).Decode(&rawActions)
	if err != nil {
		return nil, err
	}

	actions := make([]interface{}, 0)
	for _, rawAction := range rawActions {
		kind := gjson.GetBytes(rawAction, `kind`)
		switch kind.String() {
		default:
			return nil, fmt.Errorf("Unexpected action kind (%s).", kind)
		case string(CreateImplicitAccount):
			action := CreateImplicitAccountAction{}
			err = action.Unmarshal(rawAction)
			if err != nil {
				return nil, err
			}
			actions = append(actions, action)
		}
	}

	return actions, err
}

func ApplyActions(m Mockup.Mockup, taskID string, actions []interface{}) []TestResult {
	getSuccessResponse := func(action interface{}) TestResult {
		return TestResult{
			Status: Success,
			Action: action,
		}
	}

	getFailureResponse := func(description string, action interface{}) TestResult {
		return TestResult{
			Status:      Failure,
			Description: description,
			Action:      action,
		}
	}

	responses := make([]TestResult, 0)
	for _, action := range actions {
		switch content := action.(type) {
		case CreateImplicitAccountAction:
			keyPair, err := content.GenerateKey()
			if err != nil {
				responses = append(
					responses,
					getFailureResponse("Could not generate wallet.", content),
				)
				continue
			}

			// Import private key
			privateKey := keyPair.String()
			err = m.ImportSecret(privateKey, content.Payload.Name)
			if err != nil {
				responses = append(
					responses,
					getFailureResponse("Could not import wallet.", content),
				)
				continue
			}

			// Fund wallet
			err = m.Transfer(content.Payload.Balance, "bootstrap1", keyPair.Address().String())
			if err != nil {
				responses = append(
					responses,
					getFailureResponse("Could not fund wallet.", content),
				)
				continue
			}

			// Reveal wallet
			err = m.RevealWallet(content.Payload.Name)
			if err != nil {
				responses = append(
					responses,
					getFailureResponse("Could not reveal wallet.", content),
				)
				continue
			}

			responses = append(
				responses,
				getSuccessResponse(content),
			)
		}
	}

	return responses
}
