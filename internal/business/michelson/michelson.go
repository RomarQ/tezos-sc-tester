package michelson

import (
	"encoding/json"

	"github.com/romarq/tezos-sc-tester/internal/business/michelson/ast"
	MichelsonJSON "github.com/romarq/tezos-sc-tester/internal/business/michelson/json"
	"github.com/romarq/tezos-sc-tester/internal/business/michelson/micheline"
)

// JSONOfMicheline converts Michelson from "micheline" to "json" format
func JSONOfMicheline(michelsonMicheline string) (json.RawMessage, error) {
	ast, err := ParseMicheline(michelsonMicheline)
	if err != nil {
		return nil, err
	}
	return MichelsonJSON.Print(ast, "", "  ")
}

// MichelineOfJSON converts Michelson from "json" to "micheline" format
func MichelineOfJSON(michelsonJSON json.RawMessage) (string, error) {
	ast, err := ParseJSON(michelsonJSON)
	if err != nil {
		return "", err
	}
	return micheline.Print(ast, ""), nil
}

// ParseJSON parses Michelson from "json" format into an AST
func ParseJSON(michelsonJSON json.RawMessage) (ast.Node, error) {
	parser := MichelsonJSON.Parser{}
	return parser.Parse(michelsonJSON)
}

// ParseMicheline parses Michelson from "micheline" format into an AST
func ParseMicheline(michelsonMicheline string) (ast.Node, error) {
	parser := micheline.InitParser(michelsonMicheline)
	ast := parser.Parse()
	return ast, parser.Error()
}
