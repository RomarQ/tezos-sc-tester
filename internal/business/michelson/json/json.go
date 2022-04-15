package json

import (
	"encoding/json"
	"fmt"
	"strings"

	MichelsonUtils "github.com/romarq/visualtez-testing/internal/business/michelson/utils"
	"github.com/romarq/visualtez-testing/pkg/utils"
)

type (
	MichelsonJSON struct {
		Prim   string        `json:"prim,omitempty"`
		Int    string        `json:"int,omitempty"`
		String string        `json:"string,omitempty"`
		Bytes  string        `json:"bytes,omitempty"`
		Args   []interface{} `json:"args,omitempty"`
		Annots []string      `json:"annots,omitempty"`
	}
)

// Convert Michelson from JSON to Micheline format
func MichelineOfJSON(raw json.RawMessage) (string, error) {
	return toMicheline(raw)
}

func (json MichelsonJSON) isInt() bool {
	return json.Int != ""
}
func (json MichelsonJSON) isString() bool {
	return json.String != ""
}
func (json MichelsonJSON) isBytes() bool {
	return json.Bytes != ""
}
func (json MichelsonJSON) isPrim() bool {
	return json.Prim != ""
}
func (json MichelsonJSON) hasArg() bool {
	return len(json.Args) > 0
}
func (json MichelsonJSON) hasAnnots() bool {
	return len(json.Annots) > 0
}

func toMichelineInt(json MichelsonJSON) (string, error) {
	if !json.isInt() {
		return "", fmt.Errorf("Expected (Int), but received: %v", utils.PrettifyJSON(json))
	}
	return json.Int, nil
}
func toMichelineString(json MichelsonJSON) (string, error) {
	if !json.isString() {
		return "", fmt.Errorf("Expected (String), but received: %v", utils.PrettifyJSON(json))
	}
	return fmt.Sprintf(`"%s"`, json.String), nil
}
func toMichelineBytes(json MichelsonJSON) (string, error) {
	if !json.isBytes() {
		return "", fmt.Errorf("Expected (Bytes), but received: %v", utils.PrettifyJSON(json))
	}
	return json.Bytes, nil
}
func toMichelineSeq(seq []json.RawMessage) (string, error) {
	elements := make([]string, 0)
	for _, rawElement := range seq {
		element, err := toMicheline(rawElement)
		if err != nil {
			return "", err
		}
		elements = append(elements, element)
	}

	return fmt.Sprintf("{ %s }", strings.Join(elements, " ; ")), nil
}
func toMichelineAnnots(annots []string) string {
	return strings.Join(annots, " ")
}

func toMichelinePrim(michelson MichelsonJSON) (string, error) {
	if !michelson.isPrim() {
		return "", fmt.Errorf("Invalid (prim): %v", utils.PrettifyJSON(michelson))
	}

	tokens := []string{michelson.Prim}
	if len(michelson.Annots) > 0 {
		tokens = append(tokens, toMichelineAnnots(michelson.Annots))
	}

	for _, raw := range michelson.Args {
		j, err := json.Marshal(raw)
		if err != nil {
			return "", err
		}
		token, err := toMicheline(j)
		if err != nil {
			return "", err
		}
		tokens = append(tokens, token)
	}

	format := "%s"
	if michelson.supportsParenthesis() && len(tokens) > 1 {
		format = "(%s)"
	}

	return fmt.Sprintf(format, strings.Join(tokens, " ")), nil
}

func toMicheline(raw json.RawMessage) (string, error) {
	r, err := unmarshal(raw)
	if err != nil {
		return "", err
	}

	switch json := r.(type) {
	case MichelsonJSON:
		if json.isInt() {
			return toMichelineInt(json)
		}
		if json.isString() {
			return toMichelineString(json)
		}
		if json.isBytes() {
			return toMichelineBytes(json)
		}

		return toMichelinePrim(json)
	case []json.RawMessage:
		return toMichelineSeq(json)
	}

	return "", fmt.Errorf("Unexpected (michelson JSON): %v", utils.PrettifyJSON(raw))
}

func (json MichelsonJSON) supportsParenthesis() bool {
	return json.isPrim() &&
		// Cannot be a contract root
		!MichelsonUtils.IsReservedWord(json.Prim) &&
		// Match type token regex
		!MichelsonUtils.IsInstruction(json.Prim)
}

func unmarshal(raw json.RawMessage) (interface{}, error) {
	var seq []json.RawMessage
	if err := json.Unmarshal(raw, &seq); err == nil {
		return seq, nil
	}
	var prim MichelsonJSON
	err := json.Unmarshal(raw, &prim)
	return prim, err
}
