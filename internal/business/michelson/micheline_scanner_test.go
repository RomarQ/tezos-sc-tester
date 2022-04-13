package michelson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type test struct {
	Input  string
	Output []token
}

func TestTokenize(t *testing.T) {
	t.Run("Tokenize Micheline Values", func(t *testing.T) {
		runTests(t, []test{
			{
				Input: `{ "1" ; "2" }`,
				Output: []token{
					buildToken(Open_brace, "{"),
					buildToken(String, "1"),
					buildToken(Semi, ";"),
					buildToken(String, "2"),
					buildToken(Close_brace, "}"),
				},
			},
			{
				Input: `{"abc";"2"}`,
				Output: []token{
					buildToken(Open_brace, "{"),
					buildToken(String, "abc"),
					buildToken(Semi, ";"),
					buildToken(String, "2"),
					buildToken(Close_brace, "}"),
				},
			},
			{
				Input: `{ 0x00 ; 0x01 }`,
				Output: []token{
					buildToken(Open_brace, "{"),
					buildToken(Bytes, "00"),
					buildToken(Semi, ";"),
					buildToken(Bytes, "01"),
					buildToken(Close_brace, "}"),
				},
			},
			{
				Input: `{0x0abc;0xdcba;0x10}`,
				Output: []token{
					buildToken(Open_brace, "{"),
					buildToken(Bytes, "0abc"),
					buildToken(Semi, ";"),
					buildToken(Bytes, "dcba"),
					buildToken(Semi, ";"),
					buildToken(Bytes, "10"),
					buildToken(Close_brace, "}"),
				},
			},
		})
	})
	t.Run("Tokenize Micheline (Contract)", func(t *testing.T) {
		runTests(t, []test{
			{
				Input: `
					{
						storage (unit %abc);
						parameter (unit %default);
						code {}; # Test
					}
				`,
				Output: []token{
					buildToken(Open_brace, "{"),

					buildToken(Identifier, "storage"),
					buildToken(Open_paren, "("),
					buildToken(Identifier, "unit"),
					buildToken(Annot, "%abc"),
					buildToken(Close_paren, ")"),
					buildToken(Semi, ";"),

					buildToken(Identifier, "parameter"),
					buildToken(Open_paren, "("),
					buildToken(Identifier, "unit"),
					buildToken(Annot, "%default"),
					buildToken(Close_paren, ")"),
					buildToken(Semi, ";"),

					buildToken(Identifier, "code"),
					buildToken(Open_brace, "{"),
					buildToken(Close_brace, "}"),
					buildToken(Semi, ";"),

					buildToken(Comment, "# Test"),

					buildToken(Close_brace, "}"),
				},
			},
		})
	})
}

func runTests(t *testing.T, list []test) {
	for _, test := range list {
		scanner := InitMichelineScanner(test.Input)
		tokens, err := scanner.Tokenize()
		assert.Nil(t, err, "Must not fail")
		assert.Equal(t, tokens, test.Output, "Verify tokens")
	}
}
