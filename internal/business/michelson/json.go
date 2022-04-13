package michelson

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/romarq/visualtez-testing/pkg/utils"
)

type (
	MichelsonJSON struct {
		Prim   *string           `json:"prim,omitempty"`
		Int    *string           `json:"int,omitempty"`
		String *string           `json:"string,omitempty"`
		Bytes  *string           `json:"bytes,omitempty"`
		Bool   *string           `json:"bool,omitempty"`
		Args   []json.RawMessage `json:"args,omitempty"`
		Annots []string          `json:"annots,omitempty"`
	}
)

// Convert Michelson from JSON to Micheline format
func MichelineOfJSON(raw json.RawMessage) (string, error) {
	return toMicheline(raw)
}

func (json MichelsonJSON) isInt() bool {
	return json.Int != nil
}
func (json MichelsonJSON) isString() bool {
	return json.String != nil
}
func (json MichelsonJSON) isBytes() bool {
	return json.Bytes != nil
}
func (json MichelsonJSON) isBool() bool {
	return json.Bool != nil
}
func (json MichelsonJSON) isPrim() bool {
	return json.Prim != nil
}
func (json MichelsonJSON) hasArg() bool {
	return len(json.Args) > 0
}
func (json MichelsonJSON) hasAnnots() bool {
	return len(json.Annots) > 0
}

func toMichelineInt(json MichelsonJSON) (string, error) {
	if json.Int == nil {
		return "", fmt.Errorf("Expected (Int), but received: %v", utils.PrettifyJSON(json))
	}
	return *json.Int, nil
}
func toMichelineString(json MichelsonJSON) (string, error) {
	if json.String == nil {
		return "", fmt.Errorf("Expected (String), but received: %v", utils.PrettifyJSON(json))
	}
	return fmt.Sprintf(`"%s"`, *json.String), nil
}
func toMichelineBytes(json MichelsonJSON) (string, error) {
	if json.Bytes == nil {
		return "", fmt.Errorf("Expected (Bytes), but received: %v", utils.PrettifyJSON(json))
	}
	return *json.Bytes, nil
}
func toMichelineBool(json MichelsonJSON) (string, error) {
	if json.Bool == nil {
		return "", fmt.Errorf("Expected (Bool), but received: %v", utils.PrettifyJSON(json))
	}
	return *json.Bool, nil
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

func toMichelinePrim(json MichelsonJSON) (string, error) {
	if json.Prim == nil {
		return "", fmt.Errorf("Invalid (prim): %v", utils.PrettifyJSON(json))
	}

	tokens := []string{*json.Prim}
	if len(json.Annots) > 0 {
		tokens = append(tokens, toMichelineAnnots(json.Annots))
	}

	for _, raw := range json.Args {
		token, err := toMicheline(raw)
		if err != nil {
			return "", err
		}
		tokens = append(tokens, token)
	}

	format := "%s"
	if json.supportsParenthesis() && len(tokens) > 1 {
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
		if json.isBool() {
			return toMichelineBool(json)
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
		!utils.Contains(reserved_words, *json.Prim) &&
		// Match type token regex
		!isInstruction(*json.Prim)
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
