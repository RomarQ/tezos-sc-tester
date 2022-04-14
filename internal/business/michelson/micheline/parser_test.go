package micheline

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseValues(t *testing.T) {

	type test struct {
		Input  string
		Output string
	}

	runTests := func(t *testing.T, list []test) {
		for _, test := range list {
			scanner := Scanner{}
			scanner.Init(test.Input)

			parser := Parser{}
			parser.Init(scanner)

			node := parser.Parse()
			assert.Empty(t, scanner.errors, "No errors expected")
			assert.Equal(t, node.String(), test.Output, "Verify AST")
		}
	}

	t.Run("Parse Int", func(t *testing.T) {
		runTests(t, []test{
			{
				Input:  `10`,
				Output: `Int(10)`,
			},
			{
				Input:  `-10`,
				Output: `Int(-10)`,
			},
		})
	})
	t.Run("Parse Bytes", func(t *testing.T) {
		runTests(t, []test{
			{
				Input:  `0x01`,
				Output: `Bytes(0x01)`,
			},
			{
				Input:  `0x01ab`,
				Output: `Bytes(0x01ab)`,
			},
			{
				Input:  `0xabcd`,
				Output: `Bytes(0xabcd)`,
			},
		})
	})
	t.Run("Parse String", func(t *testing.T) {
		runTests(t, []test{
			{
				Input:  `"Hello World"`,
				Output: `String(Hello World)`,
			},
		})
	})
	t.Run("Parse Sequence", func(t *testing.T) {
		runTests(t, []test{
			{
				Input:  `{ 1; 2; 3 }`,
				Output: `Sequence([Int(1), Int(2), Int(3)])`,
			},
		})
	})
}

func TestParseContracts(t *testing.T) {

	type test struct {
		Input  string
		Output string
	}

	runTests := func(t *testing.T, list []test) {
		for _, test := range list {
			scanner := Scanner{}
			scanner.Init(test.Input)

			parser := Parser{}
			parser.Init(scanner)

			node := parser.Parse()
			assert.Empty(t, scanner.errors, "No errors expected")
			assert.Equal(t, node.String(), test.Output, "Verify AST")
		}
	}

	t.Run("Parse Contract", func(t *testing.T) {
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
				Output: "Sequence([Prim(storage, [], [Prim(unit, [], [])]), Prim(parameter, [], [Prim(unit, [], [])]), Prim(code, [], [Sequence([Prim(DROP, [], []), Prim(UNIT, [], []), Prim(NIL, [], [Prim(operation, [], [])]), Prim(PAIR, [], [])])]), Prim(view, [], [String(view1), Prim(unit, [], [Prim(unit, [], [Sequence([])])])])])",
			},
		})
	})
}
