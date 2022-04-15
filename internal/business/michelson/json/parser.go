package json

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/romarq/visualtez-testing/internal/business/michelson/ast"
	"github.com/romarq/visualtez-testing/internal/utils"
)

type Parser struct {
	errors []string
}

// Parse parses raw JSON into a Michelson AST
func (p *Parser) Parse(raw []byte) (node ast.Node, err error) {
	defer func() {
		if len(p.errors) > 0 {
			err = fmt.Errorf(strings.Join(p.errors, ";\n"))
		}
	}()

	r, err := unmarshal(raw)
	if err != nil {
		p.errorf("could not deserialize JSON: %s.", err)
	}

	switch obj := r.(type) {
	case MichelsonJSON:
		switch {
		case obj.isInt():
			node = ast.Int{
				Value: obj.Int,
			}
		case obj.isString():
			node = ast.String{
				Value: obj.String,
			}
		case obj.isBytes():
			node = ast.Bytes{
				Value: obj.Bytes,
			}
		case obj.isPrim():
			annotations := make([]ast.Annotation, len(obj.Annots))
			for i, el := range obj.Annots {
				annotations[i] = p.parseAnnotation(el)
			}
			arguments := make([]ast.Node, len(obj.Args))
			for i, el := range obj.Args {
				o, err := json.Marshal(el)
				if err != nil {
					p.errorf("could not parse argument of prim: %s.", err)
					break
				}
				arguments[i], _ = p.Parse(o)
			}
			node = ast.Prim{
				Prim:        obj.Prim,
				Annotations: annotations,
				Arguments:   arguments,
			}
		default:
			p.errorf("unexpected Michelson JSON: %s.", utils.PrettifyJSON(raw))
		}
	case []json.RawMessage:
		elements := make([]ast.Node, len(obj))
		for i, el := range obj {
			o, err := json.Marshal(el)
			if err != nil {
				p.errorf("could not parse element of sequence: %s.", err)
				break
			}
			elements[i], _ = p.Parse(o)
		}
		node = ast.Sequence{
			Elements: elements,
		}
	default:
		p.errorf("unexpected Michelson JSON: %s.", utils.PrettifyJSON(raw))
	}

	// Errors found during parsing will be aggregated on defer
	return
}

func (p *Parser) parseAnnotation(annot string) (annotation ast.Annotation) {
	annotation.Value = annot

	if len(annot) == 0 {
		p.errorf("Unexpected empty annotation.")
		return
	}

	switch annot[0] {
	case ':':
		annotation.Kind = ast.TypeAnnotation
	case '@':
		annotation.Kind = ast.VariableAnnotation
	case '%':
		annotation.Kind = ast.FieldAnnotation
	default:
		p.errorf("Unexpected annotation (%s).", annot)
	}

	return
}

func (p *Parser) errorf(format string, args ...interface{}) {
	p.errors = append(p.errors, fmt.Sprintf(format, args...))
}

// Deserialize raw JSON into a Michelson JSON struct.
// It can be single object or a slice.
func unmarshal(raw json.RawMessage) (interface{}, error) {
	var seq []json.RawMessage
	if err := json.Unmarshal(raw, &seq); err == nil {
		return seq, nil
	}
	var prim MichelsonJSON
	err := json.Unmarshal(raw, &prim)
	return prim, err
}
