package business

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComposeArguments(t *testing.T) {
	t.Run("Test composeArguments",
		func(t *testing.T) {
			args := composeArguments(
				TezosClientArgument{
					Kind:       Mode,
					Parameters: []string{"mockup"},
				},
				TezosClientArgument{
					Kind:       Protocol,
					Parameters: []string{"abc"},
				},
				TezosClientArgument{
					Kind:       BaseDirectory,
					Parameters: []string{"src"},
				},
				TezosClientArgument{
					Kind:       ProtocolConstants,
					Parameters: []string{"ProtocolConstants"},
				},
				TezosClientArgument{
					Kind:       BootstrapAccounts,
					Parameters: []string{"BootstrapAccounts"},
				},
				TezosClientArgument{
					Kind:       BurnCap,
					Parameters: []string{"BurnCap"},
				},
				TezosClientArgument{
					Kind:       Fee,
					Parameters: []string{"Fee"},
				},
			)
			assert.Equal(
				t,
				[]string{
					"-M", "mockup",
					"-p", "abc",
					"-d", "src",
					"--protocol-constants", "ProtocolConstants",
					"--bootstrap-accounts", "BootstrapAccounts",
					"--burn-cap", "BurnCap",
					"--fee", "Fee",
				},
				args,
			)
		})
}
