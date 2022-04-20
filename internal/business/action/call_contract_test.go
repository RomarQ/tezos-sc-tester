package action

import (
	"encoding/json"
	"testing"

	"github.com/romarq/visualtez-testing/internal/business/michelson/ast"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshal_CallContractAction(t *testing.T) {
	t.Run("Test CallContractAction Unmarshal (Valid)",
		func(t *testing.T) {
			action := CallContractAction{
				raw: json.RawMessage(`
					{
						"kind": "call_contract",
						"payload": {
							"recipient":	"contract_1",
							"sender":		"sender_name",
							"entrypoint":	"do_something",
							"amount":		"10",
							"parameter":	{ "prim": "Unit" }
						}
					}
				`),
			}
			err := action.Unmarshal()
			assert.Nil(t, err, "Must not fail")
			assert.Equal(
				t,
				"contract_1",
				action.Recipient,
				"Assert name",
			)
			assert.Equal(
				t,
				"sender_name",
				action.Sender,
				"Assert sender",
			)
			assert.Equal(
				t,
				"10",
				action.Amount.String(),
				"Assert amount",
			)
			assert.Equal(
				t,
				ast.Prim{
					Prim: "Unit",
				},
				action.Parameter,
				"Assert parameter",
			)
		})
	t.Run("Test CallContractAction Unmarshal (Invalid name)",
		func(t *testing.T) {
			action := CallContractAction{
				raw: json.RawMessage(`
					{
						"kind": "call_contract",
						"payload": {
							"recipient":	"contract 1",
							"sender":		"sender_name",
							"amount":		"10",
							"parameter":	{ "prim": "Unit" }
						}
					}
				`),
			}
			err := action.Unmarshal()
			assert.NotNil(t, err, "Must fail (name is invalid)")
			assert.Equal(t, err.Error(), "String (contract 1) does not match pattern '^[a-zA-Z0-9_]+$'.", "Assert error message")
		})
	t.Run("Test CallContractAction Unmarshal (Invalid sender)",
		func(t *testing.T) {
			action := CallContractAction{
				raw: json.RawMessage(`
					{
						"kind": "call_contract",
						"payload": {
							"recipient":	"contract_1",
							"sender":		"sender name",
							"amount":		"10",
							"parameter":	{ "prim": "Unit" }
						}
					}
				`),
			}
			err := action.Unmarshal()
			assert.NotNil(t, err, "Must fail (sender is invalid)")
			assert.Equal(t, err.Error(), "String (sender name) does not match pattern '^[a-zA-Z0-9_]+$'.", "Assert error message")
		})
	t.Run("Test CallContractAction Unmarshal (Invalid entrypoint length)",
		func(t *testing.T) {
			action := CallContractAction{
				raw: json.RawMessage(`
					{
						"kind": "call_contract",
						"payload": {
							"recipient":	"contract_1",
							"sender":		"sender_name",
							"entrypoint":	"abcdefghijlmnopqrstuvxz123456789",
							"amount":		"10",
							"parameter":	{ "prim": "Unit" }
						}
					}
				`),
			}
			err := action.Unmarshal()
			assert.NotNil(t, err, "Must fail (Invalid entrypoint length)")
			assert.Equal(t, err.Error(), "String (abcdefghijlmnopqrstuvxz123456789) does not match pattern '^[a-zA-Z0-9_]{1,31}$'.", "Assert error message")
		})
	t.Run("Test CallContractAction Unmarshal (Invalid chars in entrypoint)",
		func(t *testing.T) {
			action := CallContractAction{
				raw: json.RawMessage(`
					{
						"kind": "call_contract",
						"payload": {
							"recipient":	"contract_1",
							"sender":		"sender_name",
							"entrypoint":	"a.a",
							"amount":		"10",
							"parameter":	{ "prim": "Unit" }
						}
					}
				`),
			}
			err := action.Unmarshal()
			assert.NotNil(t, err, "Must fail (Invalid chars in entrypoint)")
			assert.Equal(t, err.Error(), "String (a.a) does not match pattern '^[a-zA-Z0-9_]{1,31}$'.", "Assert error message")
		})
	t.Run("Test CallContractAction Unmarshal (Missing fields)",
		func(t *testing.T) {
			action := CallContractAction{raw: json.RawMessage(`{}`)}
			err := action.Unmarshal()
			assert.NotNil(t, err, "Must fail (Missing fields)")
			assert.Equal(t, err.Error(), "Action of kind (call_contract) misses the following fields [recipient, sender, entrypoint, parameter].", "Assert error message")
		})
}
