package action

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	t.Run("Test CreateImplicitAccountAction Unmarshal",
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
				int64(10),
				action.Balance,
				"Assert balance",
			)
		})
}
