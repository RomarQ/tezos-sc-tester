package utils_test

import (
	"testing"

	"github.com/romarq/tezos-sc-tester/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {
	t.Run("Generate Key",
		func(t *testing.T) {
			keyPair, err := utils.GenerateKey()
			assert.Nil(t, err, "Must not fail")
			assert.Equal(
				t,
				"edsk",
				keyPair.String()[0:4],
			)
		})

	t.Run("Validate Chain ID", func(t *testing.T) {
		assert.True(t, utils.ValidateChainID("NetXynUjJNZm7wi"), "A valid chain_id")
		assert.False(t, utils.ValidateChainID("NetSomething"), "An invalid chain_id")
	})
}
