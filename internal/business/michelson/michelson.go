package michelson

import (
	"encoding/json"

	MichelsonJSON "github.com/romarq/visualtez-testing/internal/business/michelson/json"
	"github.com/romarq/visualtez-testing/internal/business/michelson/micheline"
)

// JSONOfMicheline converts Michelson from "micheline" to "json" format
func JSONOfMicheline(michelsonMicheline string) (json.RawMessage, error) {
	parser := micheline.InitParser(michelsonMicheline)
	ast := parser.Parse()
	if parser.HasErrors() {
		return nil, parser.Error()
	}
	return MichelsonJSON.Print(ast, "", "  ")
}

// MichelineOfJSON converts Michelson from "json" to "micheline" format
func MichelineOfJSON(michelsonJSON json.RawMessage) (string, error) {
	parser := MichelsonJSON.Parser{}
	ast, err := parser.Parse(michelsonJSON)
	if err != nil {
		return "", err
	}
	return micheline.Print(ast), nil
}
