package action

import (
	"io"
	"math/big"
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
						"balance": "10"
					}
				}
			]
			`))
			actions, err := GetActions(reqBody)
			assert.Nil(t, err, "Must not fail")
			assert.Len(
				t,
				actions,
				1,
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
						"balance": "10"
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
			action_createImplicitAccount_alice := CreateImplicitAccountAction{
				Name:    "alice",
				Balance: business.MutezOfFloat(big.NewFloat(10)),
			}
			action_createImplicitAccount_bob := CreateImplicitAccountAction{
				Name:    "bob",
				Balance: business.MutezOfFloat(big.NewFloat(10)),
			}
			actions := []IAction{
				&CreateImplicitAccountActionMock{action_createImplicitAccount_alice},
				&CreateImplicitAccountActionMock{action_createImplicitAccount_bob},
			}
			results := ApplyActions(business.Mockup{}, actions)
			assert.Equal(
				t,
				[]ActionResult{
					{
						Status: Success,
						Kind:   CreateImplicitAccount,
						Action: action_createImplicitAccount_alice.json,
						Result: map[string]interface{}{},
					},
					{
						Status: Failure,
						Kind:   CreateImplicitAccount,
						Action: action_createImplicitAccount_bob.json,
						Result: map[string]interface{}{
							"details": "ERROR",
						},
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

func (action CreateImplicitAccountActionMock) Run(mockup business.Mockup) ActionResult {
	if action.Name == "bob" {
		return action.buildFailureResult("ERROR")
	}
	return action.buildSuccessResult(map[string]interface{}{})
}
