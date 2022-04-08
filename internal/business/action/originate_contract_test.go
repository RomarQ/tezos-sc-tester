package action

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal_OriginateContractAction(t *testing.T) {
	t.Run("Test OriginateContractAction Unmarshal (Valid)",
		func(t *testing.T) {
			rawJson, err := json.MarshalIndent(
				map[string]interface{}{
					"name":    "contract_1",
					"balance": 10,
					"code":    "{storage unit ; parameter (unit %entrypoint) ; code { CDR ; NIL operation ; PAIR } }",
					"storage": "Unit",
				},
				"",
				"",
			)
			assert.Nil(t, err)
			action := OriginateContractAction{}
			err = action.Unmarshal(rawJson)
			assert.Nil(t, err, "Must not fail")
			assert.Equal(
				t,
				"contract_1",
				action.Name,
				"Assert name",
			)
			assert.Equal(
				t,
				float64(10),
				action.Balance,
				"Assert balance",
			)
			assert.Equal(
				t,
				"{storage unit ; parameter (unit %entrypoint) ; code { CDR ; NIL operation ; PAIR } }",
				action.Code,
				"Assert code",
			)
			assert.Equal(
				t,
				"Unit",
				action.Storage,
				"Assert storage",
			)
		})
	t.Run("Test OriginateContractAction Unmarshal (Invalid name)",
		func(t *testing.T) {
			rawJson, err := json.MarshalIndent(
				map[string]interface{}{
					"name":    "contract 1",
					"balance": 10,
				},
				"",
				"",
			)
			assert.Nil(t, err)
			action := OriginateContractAction{}
			err = action.Unmarshal(rawJson)
			assert.NotNil(t, err, "Must fail (name is invalid)")
			assert.Equal(t, err.Error(), "String (contract 1) does not match pattern '^[a-zA-Z0-9._-]+$'.", "Assert error message")
		})
	t.Run("Test OriginateContractAction Unmarshal (Missing fields)",
		func(t *testing.T) {
			rawJson := json.RawMessage(`{}`)
			action := OriginateContractAction{}
			err := action.Unmarshal(rawJson)
			assert.NotNil(t, err, "Must fail (Missing fields)")
			assert.Equal(t, err.Error(), "Action of kind (originate_contract) misses the following fields [name, code, storage].", "Assert error message")
		})
}
