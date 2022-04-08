package action

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal_CreateImplicitAccountAction(t *testing.T) {
	t.Run("Test CreateImplicitAccountAction Unmarshal (Valid)",
		func(t *testing.T) {
			action := CreateImplicitAccountAction{}
			err := action.Unmarshal(
				json.RawMessage(`
					{
						"name":    "bob",
						"balance": 10
					}
				`),
			)
			assert.Nil(t, err, "Must not fail")
			assert.Equal(
				t,
				"bob",
				action.Name,
				"Assert name",
			)
			assert.Equal(
				t,
				float64(10),
				action.Balance,
				"Assert balance",
			)
		})

	t.Run("Test CreateImplicitAccountAction Unmarshal (Invalid)",
		func(t *testing.T) {
			action := CreateImplicitAccountAction{}
			err := action.Unmarshal(
				json.RawMessage(`
					{
						"name":    "bob A",
						"balance": 10
					}
				`),
			)
			assert.NotNil(t, err, "Must fail (name is invalid)")
			assert.Equal(t, err.Error(), "String (bob A) does not match pattern '^[a-zA-Z0-9._-]+$'.", "Assert error message")
		})
}
