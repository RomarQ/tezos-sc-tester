package action

import (
	"encoding/json"
	"testing"

	"github.com/romarq/visualtez-testing/internal/business/michelson"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshal_OriginateContractAction(t *testing.T) {
	t.Run("Test OriginateContractAction Unmarshal (Valid)",
		func(t *testing.T) {
			code := json.RawMessage(`
				[
					{ "prim": "storage", "args": [ { "prim": "unit" } ] },
					{ "prim": "parameter", "args": [ { "prim": "unit", "annots": ["%entrypoint"] } ] },
					{ "prim": "code", "args": [ { "prim": "CDR" }, { "prim": "NIL", "args": [ { "prim": "operation" } ] }, { "prim": "PAIR" } ] }
				]
			`)
			storage := json.RawMessage(`
				{ "prim": "Unit" }
			`)
			rawJson, err := json.MarshalIndent(
				map[string]interface{}{
					"kind": OriginateContract,
					"payload": map[string]interface{}{
						"name":    "contract_1",
						"balance": "10",
						"code":    code,
						"storage": storage,
					},
				},
				"",
				"",
			)
			assert.Nil(t, err)
			action := OriginateContractAction{raw: rawJson}
			err = action.Unmarshal()
			assert.Nil(t, err, "Must not fail")
			assert.Equal(
				t,
				"contract_1",
				action.Name,
				"Assert name",
			)
			assert.Equal(
				t,
				"10",
				action.Balance.String(),
				"Assert balance",
			)
			ast, _ := michelson.ParseJSON(code)
			assert.Equal(
				t,
				ast.String(),
				action.Code.String(),
				"Assert code",
			)
			ast, _ = michelson.ParseJSON(storage)
			assert.Equal(
				t,
				ast.String(),
				action.Storage.String(),
				"Assert storage",
			)
		})
	t.Run("Test OriginateContractAction Unmarshal (Invalid name)",
		func(t *testing.T) {
			rawJson, err := json.MarshalIndent(
				map[string]interface{}{
					"kind": OriginateContract,
					"payload": map[string]interface{}{
						"name":    "contract 1",
						"balance": "10",
					},
				},
				"",
				"",
			)
			assert.Nil(t, err)
			action := OriginateContractAction{raw: rawJson}
			err = action.Unmarshal()
			assert.NotNil(t, err, "Must fail (name is invalid)")
			assert.Equal(t, err.Error(), "String (contract 1) does not match pattern '^[a-zA-Z0-9_]+$'.", "Assert error message")
		})
	t.Run("Test OriginateContractAction Unmarshal (Missing fields)",
		func(t *testing.T) {
			rawJson := json.RawMessage(`{}`)
			action := OriginateContractAction{raw: rawJson}
			err := action.Unmarshal()
			assert.NotNil(t, err, "Must fail (Missing fields)")
			assert.Equal(t, err.Error(), "Action of kind (originate_contract) misses the following fields [name, code, storage].", "Assert error message")
		})
}
