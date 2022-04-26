package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateKey(t *testing.T) {
	t.Run("Generate Key",
		func(t *testing.T) {
			keyPair, err := GenerateKey()
			assert.Nil(t, err, "Must not fail")
			assert.Equal(
				t,
				"edsk",
				keyPair.String()[0:4],
			)
		})

	t.Run("Validate Chain ID", func(t *testing.T) {
		assert.True(t, ValidateChainID("NetXynUjJNZm7wi"), "A valid chain_id")
		assert.False(t, ValidateChainID("NetSomething"), "An invalid chain_id")
	})
}
