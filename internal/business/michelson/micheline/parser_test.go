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
			parser := InitParser(test.Input)

			node := parser.Parse()
			assert.Empty(t, parser.scanner.errors, "No errors expected")
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
				Output: `Bytes(01)`,
			},
			{
				Input:  `0x01ab`,
				Output: `Bytes(01ab)`,
			},
			{
				Input:  `0xabcd`,
				Output: `Bytes(abcd)`,
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
	t.Run("Parse right combs", func(t *testing.T) {
		runTests(t, []test{
			{
				Input: `Pair (Pair "tz1f2k9M3ztqtbCTk5EmEepboEJxksXvafaU" True)
				(Pair {} {} {} {} 0)
				{}`,
				Output: `Prim(Pair, [], [Prim(Pair, [], [String(tz1f2k9M3ztqtbCTk5EmEepboEJxksXvafaU), Prim(True, [], [])]), Prim(Pair, [], [Sequence([]), Sequence([]), Sequence([]), Sequence([]), Int(0)]), Sequence([])])`,
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
			parser := InitParser(test.Input)

			node := parser.Parse()
			assert.Empty(t, parser.scanner.errors, "No errors expected")
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
				Output: "Sequence([Prim(storage, [], [Prim(unit, [], [])]), Prim(parameter, [], [Prim(unit, [], [])]), Prim(code, [], [Sequence([Prim(DROP, [], []), Prim(UNIT, [], []), Prim(NIL, [], [Prim(operation, [], [])]), Prim(PAIR, [], [])])]), Prim(view, [], [String(view1), Prim(unit, [], []), Prim(unit, [], []), Sequence([])])])",
			},
			{
				Input: `
				{
					storage (unit %abc);
					parameter (unit %ep);
					code {
						DROP;
						UNIT;
						NIL operation;
						PAIR;
					};
					view "view1" unit (pair int int) {
						DROP;
						PUSH (pair int int) (Pair 1 1)
					}
				}
				`,
				Output: "Sequence([Prim(storage, [], [Prim(unit, [%abc], [])]), Prim(parameter, [], [Prim(unit, [%ep], [])]), Prim(code, [], [Sequence([Prim(DROP, [], []), Prim(UNIT, [], []), Prim(NIL, [], [Prim(operation, [], [])]), Prim(PAIR, [], [])])]), Prim(view, [], [String(view1), Prim(unit, [], []), Prim(pair, [], [Prim(int, [], []), Prim(int, [], [])]), Sequence([Prim(DROP, [], []), Prim(PUSH, [], [Prim(pair, [], [Prim(int, [], []), Prim(int, [], [])]), Prim(Pair, [], [Int(1), Int(1)])])])])])",
			},
		})
	})
}
