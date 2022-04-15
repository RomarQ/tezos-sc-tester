package json

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/romarq/visualtez-testing/internal/business/michelson/ast"
	"github.com/romarq/visualtez-testing/pkg/utils"
)

type Parser struct {
	errors []string
}

func (p *Parser) Parse(raw []byte) ast.Node {
	r, err := unmarshal(raw)
	if err != nil {
		p.errorf("Could not deserialize JSON: %s", err)
	}

	switch obj := r.(type) {
	case MichelsonJSON:
		switch {
		case obj.isInt():
			return ast.Int{
				Value: obj.Int,
			}
		case obj.isString():
			return ast.String{
				Value: obj.String,
			}
		case obj.isBytes():
			return ast.Bytes{
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
					p.errorf("Could not parse argument of prim. %v", err)
					break
				}
				arguments[i] = p.Parse(o)
			}
			return ast.Prim{
				Prim:        obj.Prim,
				Annotations: annotations,
				Arguments:   arguments,
			}
		}
	case []json.RawMessage:
		elements := make([]ast.Node, len(obj))
		for i, el := range obj {
			o, err := json.Marshal(el)
			if err != nil {
				p.errorf("Could not parse element of sequence. %v", err)
				break
			}
			elements[i] = p.Parse(o)
		}
		return ast.Sequence{
			Elements: elements,
		}
	}

	p.errorf("Unexpected Michelson JSON: %s", utils.PrettifyJSON(raw))
	return nil
}

func (p *Parser) HasErrors() bool {
	return len(p.errors) > 0
}

func (p *Parser) Error() error {
	return fmt.Errorf(strings.Join(p.errors, ";\n"))
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
		p.errorf("Unexpected annotation: (%s)", annot)
	}

	return
}

func (p *Parser) errorf(format string, args ...interface{}) {
	p.errors = append(p.errors, fmt.Sprintf(format, args...))
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
