package business

import (
	"testing"

	"github.com/romarq/visualtez-testing/internal/config"
	"github.com/romarq/visualtez-testing/internal/logger"
	"github.com/stretchr/testify/assert"
)

func runMockupAndTeardown(t *testing.T, testFunc func(mockup Mockup)) {
	mockup := Mockup{
		TaskID: "task",
		Config: config.Config{
			Log: config.LogConfig{
				Location: "../../.tmp_test/api.log",
				Level:    "debug",
			},
			Tezos: config.TezosConfig{
				DefaultProtocol: "ProtoALphaALphaALphaALphaALphaALphaALphaALphaDdp3zK",
				BaseDirectory:   "../../tezos-bin",
				RevealFee:       1000,
				Originator:      "bootstrap2",
			},
		},
	}
	logger.SetupLogger(mockup.Config.Log.Location, mockup.Config.Log.Level)

	err := mockup.Bootstrap()
	assert.NoError(t, err)

	testFunc(mockup)

	err = mockup.Teardown()
	assert.NoError(t, err)
}

func TestComposeArguments(t *testing.T) {
	t.Run("Validate addresses after mockup bootstrap", func(t *testing.T) {
		runMockupAndTeardown(t, func(mockup Mockup) {
			assert.Equal(t, mockup.Addresses, map[string]string{
				"bootstrap1": "tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx",
				"bootstrap2": "tz1gjaF81ZRRvdzjobyfVNsAeSC6PScjfQwN",
				"bootstrap3": "tz1faswCTDciRzE4oJ9jn2Vm2dvjeyA9fUzU",
				"bootstrap4": "tz1b7tUupMgCNw2cCLpKTkSD1NZzB5TkP2sv",
				"bootstrap5": "tz1ddb9NMYHZi5UzPdzTZMYQQZoMub195zgv",
			})
		})
	})
	t.Run("Compose arguments",
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
