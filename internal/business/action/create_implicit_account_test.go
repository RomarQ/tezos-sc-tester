package action

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal_CreateImplicitAccountAction(t *testing.T) {
	t.Run("Test CreateImplicitAccountAction Unmarshal (Valid)",
		func(t *testing.T) {
			rawAction := Action{
				Kind: CreateImplicitAccount,
				Payload: json.RawMessage(`
					{
						"name":    "bob",
						"balance": "10"
					}
				`),
			}
			action := CreateImplicitAccountAction{}
			err := action.Unmarshal(rawAction)
			assert.Nil(t, err, "Must not fail")
			assert.Equal(
				t,
				"bob",
				action.Name,
				"Assert name",
			)
			assert.Equal(
				t,
				"10",
				action.Balance.String(),
				"Assert balance",
			)
		})

	t.Run("Test CreateImplicitAccountAction Unmarshal (Invalid)",
		func(t *testing.T) {
			rawAction := Action{
				Kind: CreateImplicitAccount,
				Payload: json.RawMessage(`
					{
						"name":    "bob A",
						"balance": "10"
					}
				`),
			}
			action := CreateImplicitAccountAction{}
			err := action.Unmarshal(rawAction)
			assert.NotNil(t, err, "Must fail (name is invalid)")
			assert.Equal(t, err.Error(), "String (bob A) does not match pattern '^[a-zA-Z0-9_]+$'.", "Assert error message")
		})
}
