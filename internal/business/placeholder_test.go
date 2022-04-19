package business

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandPlaceholders(t *testing.T) {
	t.Run("Expand Account Placeholders", func(t *testing.T) {
		addresses := map[string]string{
			"a1": "tz1",
			"a2": "tz2",
		}
		bytes := []byte("TEST__ADDRESS_OF_ACCOUNT__a1----TEST__ADDRESS_OF_ACCOUNT__a2")
		assert.Equal(t, string(ExpandAccountPlaceholders(addresses, bytes)), "tz1----tz2")
	})
}
