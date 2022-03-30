package action

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/romarq/visualtez-testing/internal/business"
	"github.com/stretchr/testify/assert"
)

func TestGetActions(t *testing.T) {
	t.Run("Test GetActions (No errors)",
		func(t *testing.T) {
			reqBody := io.NopCloser(strings.NewReader(`
			[
				{
					"kind": "create_implicit_account",
					"payload": {
						"name": "alice",
						"balance": 10
					}
				}
			]
			`))
			actions, err := GetActions(reqBody)
			assert.Nil(t, err, "Must not fail")
			assert.ElementsMatch(
				t,
				[]IAction{
					&CreateImplicitAccountAction{
						Name:    "alice",
						Balance: float64(10),
					},
				},
				actions,
				"Validate parsed actions",
			)

			reqBody.Close()
		})
	t.Run("Test GetActions (With errors)",
		func(t *testing.T) {
			reqBody := io.NopCloser(strings.NewReader(`
			[
				{
					"kind": "create_implicit_account",
					"payload": {
						"name": "alice",
						"balance": 10
					}
				},
				{
					"kind": "THIS_ACTION_DOES_NOT_EXIST"
				}
			]
			`))
			actions, err := GetActions(reqBody)
			assert.Equal(t, "Unexpected action kind (THIS_ACTION_DOES_NOT_EXIST).", err.Error(), "Must fail")
			assert.ElementsMatch(
				t,
				[]IAction{},
				actions,
				"Expects an empty slice",
			)

			reqBody.Close()
		})
}

func TestApplyActions(t *testing.T) {
	t.Run("Test ApplyActions",
		func(t *testing.T) {
			action_createImplicitAccount_alice := &CreateImplicitAccountActionMock{
				CreateImplicitAccountAction: CreateImplicitAccountAction{
					Name:    "alice",
					Balance: float64(10),
				},
			}
			action_createImplicitAccount_bob := &CreateImplicitAccountActionMock{
				CreateImplicitAccountAction: CreateImplicitAccountAction{
					Name:    "bob",
					Balance: float64(10),
				},
			}
			actions := []IAction{
				action_createImplicitAccount_alice,
				action_createImplicitAccount_bob,
			}
			results := ApplyActions(business.Mockup{}, actions)
			assert.ElementsMatch(
				t,
				[]TestResult{
					{
						Status: Success,
						Kind:   CreateImplicitAccount,
						Action: action_createImplicitAccount_alice,
					},
					{
						Status:      Failure,
						Kind:        CreateImplicitAccount,
						Description: "FAIL",
						Action:      action_createImplicitAccount_bob,
					},
				},
				results,
				"Validate actions results",
			)
		})
}

// Mocks

type CreateImplicitAccountActionMock struct {
	CreateImplicitAccountAction
}

func (action *CreateImplicitAccountActionMock) Run(mockup business.Mockup) error {
	if action.Name == "bob" {
		return fmt.Errorf("FAIL")
	}
	return nil
}
