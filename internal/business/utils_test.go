package business

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateKey(t *testing.T) {
	t.Run("Test GenerateKey",
		func(t *testing.T) {
			keyPair, err := GenerateKey()
			assert.Nil(t, err, "Must not fail")
			assert.Equal(
				t,
				"edsk",
				keyPair.String()[0:4],
			)
		})
}
