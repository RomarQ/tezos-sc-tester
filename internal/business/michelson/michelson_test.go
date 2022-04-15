package michelson

import (
	"encoding/json"
	"testing"

	"github.com/romarq/visualtez-testing/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestJSONOfMicheline(t *testing.T) {

	type (
		test struct {
			Input      string
			Output     json.RawMessage
			ShouldFail bool
		}
	)

	runTests := func(t *testing.T, list []test) {
		for _, test := range list {
			result, err := JSONOfMicheline(test.Input)
			if test.ShouldFail {
				assert.Error(t, err, "Must fail")
			} else {
				assert.NoError(t, err, "Must not fail")
			}
			assert.Equal(t, utils.PrettifyJSON(result), utils.PrettifyJSON(test.Output), "Assert result")
		}
	}

	t.Run("Values", func(t *testing.T) {
		runTests(t, []test{
			{
				Input:      "1",
				Output:     json.RawMessage(`{ "int": "1" }`),
				ShouldFail: false,
			},
			{
				Input:      `"Hello World"`,
				Output:     json.RawMessage(`{ "string": "Hello World" }`),
				ShouldFail: false,
			},
			{
				Input:      `0x01`,
				Output:     json.RawMessage(`{ "bytes": "01" }`),
				ShouldFail: false,
			},
			{
				Input:      `{ 1 ; 2 }`,
				Output:     json.RawMessage(`[ { "int": "1" }, { "int": "2" } ]`),
				ShouldFail: false,
			},
		})
	})

	t.Run("Types", func(t *testing.T) {
		runTests(t, []test{
			{
				Input:      "int",
				Output:     json.RawMessage(`{ "prim": "int" }`),
				ShouldFail: false,
			},
			{
				Input:      `string %abc`,
				Output:     json.RawMessage(`{ "prim": "string", "annots": ["%abc"] }`),
				ShouldFail: false,
			},
			{
				Input:      `(pair unit unit)`,
				Output:     json.RawMessage(`{ "prim": "pair", "args": [ { "prim": "unit" }, { "prim": "unit" } ] }`),
				ShouldFail: false,
			},
		})
	})

	t.Run("Contracts", func(t *testing.T) {
		runTests(t, []test{
			{
				Input: `
				{
					storage unit;
					parameter unit;
					code {
						DROP;
						UNIT;
						NIL operation;
						PAIR;
					};
					view "view1" unit unit {}
				}
				`,
				Output: json.RawMessage(`
				[
					{
						"prim": "storage",
						"args": [
							{
								"prim": "unit"
							}
						]
					},
					{
						"prim": "parameter",
						"args": [
							{
								"prim": "unit"
							}
						]
					},
					{
						"prim": "code",
						"args": [
							[
								{
									"prim": "DROP"
								},
								{
									"prim": "UNIT"
								},
								{
									"prim": "NIL",
									"args": [
										{
											"prim": "operation"
										}
									]
								},
								{
									"prim": "PAIR"
								}
							]
						]
					},
					{
						"prim": "view",
						"args": [
							{
								"string": "view1"
							},
							{
								"prim": "unit"
							},
							{
								"prim": "unit"
							},
							[]
						]
					}
				]
				`),
				ShouldFail: false,
			},
			{
				Input:      `string %abc`,
				Output:     json.RawMessage(`{ "prim": "string", "annots": ["%abc"] }`),
				ShouldFail: false,
			},
			{
				Input:      `(pair unit unit)`,
				Output:     json.RawMessage(`{ "prim": "pair", "args": [ { "prim": "unit" }, { "prim": "unit" } ] }`),
				ShouldFail: false,
			},
		})
	})
}
