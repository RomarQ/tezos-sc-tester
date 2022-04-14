package micheline

import (
	"testing"

	"github.com/romarq/visualtez-testing/internal/business/michelson/micheline/token"
	"github.com/stretchr/testify/assert"
)

func TestTokenize(t *testing.T) {

	type (
		output struct {
			position int
			kind     token.Kind
			text     string
		}
		test struct {
			Input  string
			Output []output
		}
	)

	runTests := func(t *testing.T, list []test) {
		for _, test := range list {
			scanner := Scanner{}
			scanner.Init(test.Input)

			for _, output := range test.Output {
				position, kind, text := scanner.Scan()
				assert.Equal(t, kind.String(), output.kind.String(), "Assert token kind")
				if output.position != 0 {
					assert.Equal(t, position, output.position, "Assert token position")
				}
				if output.text != "" {
					assert.Equal(t, text, output.text, "Assert token text")
				}
			}
		}
	}

	t.Run("Tokenize Micheline Values", func(t *testing.T) {
		runTests(t, []test{
			{
				Input: `{ "1" ; "2" }`,
				Output: []output{
					{
						kind: token.Open_brace,
						text: "{",
					},
					{
						position: 2,
						kind:     token.String,
						text:     "1",
					},
					{
						position: 6,
						kind:     token.Semi,
						text:     ";",
					},
					{
						position: 8,
						kind:     token.String,
						text:     "2",
					},
					{
						position: 12,
						kind:     token.Close_brace,
						text:     "}",
					},
				},
			},
			{
				Input: `{"abc";"2"}`,
				Output: []output{
					{
						kind: token.Open_brace,
						text: "{",
					},
					{
						position: 1,
						kind:     token.String,
						text:     "abc",
					},
					{
						position: 6,
						kind:     token.Semi,
						text:     ";",
					},
					{
						position: 7,
						kind:     token.String,
						text:     "2",
					},
					{
						position: 10,
						kind:     token.Close_brace,
						text:     "}",
					},
				},
			},
			{
				Input: `{ 0x00 ; 0x01 }`,
				Output: []output{
					{
						kind: token.Open_brace,
						text: "{",
					},
					{
						kind: token.Bytes,
						text: "0x00",
					},
					{
						kind: token.Semi,
						text: ";",
					},
					{
						kind: token.Bytes,
						text: "0x01",
					},
					{
						kind: token.Close_brace,
						text: "}",
					},
				},
			},
			{
				Input: `{0x0abc;0xdcba;0x10}`,
				Output: []output{
					{
						kind: token.Open_brace,
						text: "{",
					},
					{
						kind: token.Bytes,
						text: "0x0abc",
					},
					{
						kind: token.Semi,
						text: ";",
					},
					{
						kind: token.Bytes,
						text: "0xdcba",
					}, {
						kind: token.Semi,
						text: ";",
					},
					{
						kind: token.Bytes,
						text: "0x10",
					},
					{
						kind: token.Close_brace,
						text: "}",
					},
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
				Output: []output{
					// {
					{
						kind: token.Open_brace,
						text: "{",
					},
					// storage (unit %abc);
					{
						kind: token.Identifier,
						text: "storage",
					},
					{
						kind: token.Open_paren,
						text: "(",
					},
					{
						kind: token.Identifier,
						text: "unit",
					},
					{
						kind: token.Annot,
						text: "%abc",
					}, {
						kind: token.Close_paren,
						text: ")",
					},
					{
						kind: token.Semi,
						text: ";",
					},
					// parameter (unit %default);
					{
						kind: token.Identifier,
						text: "parameter",
					},
					{
						kind: token.Open_paren,
						text: "(",
					},
					{
						kind: token.Identifier,
						text: "unit",
					},
					{
						kind: token.Annot,
						text: "%default",
					}, {
						kind: token.Close_paren,
						text: ")",
					},
					{
						kind: token.Semi,
						text: ";",
					},
					// code {}; # Test
					{
						kind: token.Identifier,
						text: "code",
					},
					{
						kind: token.Open_brace,
						text: "{",
					},
					{
						kind: token.Close_brace,
						text: "}",
					},
					{
						kind: token.Semi,
						text: ";",
					},
					{
						kind: token.Comment,
						text: "# Test",
					},
					// }
					{
						kind: token.Close_brace,
						text: "}",
					},
				},
			},
		})
	})
}
