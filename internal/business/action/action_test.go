package action

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/romarq/visualtez-testing/internal/business"
	"github.com/stretchr/testify/assert"
)

func TestGetActions(t *testing.T) {
	t.Run("Test GetActions (No errors)",
		func(t *testing.T) {
			rawActions := []Action{
				{
					Kind: CreateImplicitAccount,
					Payload: json.RawMessage(`
						{
							"name": "alice",
							"balance": "10"
						}
					`),
				},
			}
			actions, err := GetActions(rawActions)
			assert.Nil(t, err, "Must not fail")
			assert.Len(
				t,
				actions,
				1,
				"Validate parsed actions",
			)
		})
	t.Run("Test GetActions (With errors)",
		func(t *testing.T) {
			rawActions := []Action{
				{
					Kind: CreateImplicitAccount,
					Payload: json.RawMessage(`
						{
							"name": "alice",
							"balance": "10"
						}
					`),
				},
				{
					Kind: "THIS_ACTION_DOES_NOT_EXIST",
				},
			}
			actions, err := GetActions(rawActions)
			assert.Equal(t, "Unexpected action kind (THIS_ACTION_DOES_NOT_EXIST).", err.Error(), "Must fail")
			assert.ElementsMatch(
				t,
				[]IAction{},
				actions,
				"Expects an empty slice",
			)
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
						Action: action_createImplicitAccount_alice.Action(),
						Result: map[string]interface{}{},
					},
					{
						Status: Failure,
						Action: action_createImplicitAccount_bob.Action(),
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

func (action CreateImplicitAccountActionMock) Run(mockup business.Mockup) (interface{}, bool) {
	if action.Name == "bob" {
		return "ERROR", false
	}
	return map[string]interface{}{}, true
}
