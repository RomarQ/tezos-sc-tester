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
		p.errorf("Could not deserialize json: %s", err)
	}

	switch obj := r.(type) {
	case MichelsonJSON:
		if obj.isInt() {
			return ast.Int{
				Value: obj.Int,
			}
		}
		if obj.isString() {
			return ast.String{
				Value: obj.String,
			}
		}
		if obj.isBytes() {
			return ast.Bytes{
				Value: obj.Bytes,
			}
		}

		annotations := make([]ast.Annotation, len(obj.Annots))
		for i, el := range obj.Annots {
			annotations[i] = p.parseAnnotation(el)
		}
		arguments := make([]ast.Node, len(obj.Args))
		for i, el := range obj.Args {
			o, err := json.Marshal(el)
			if err != nil {
				p.errorf("%v", err)
			}
			arguments[i] = p.Parse(o)
		}
		return ast.Prim{
			Prim:        obj.Prim,
			Annotations: annotations,
			Arguments:   arguments,
		}
	case []json.RawMessage:
		elements := make([]ast.Node, len(obj))
		for i, el := range obj {
			o, err := json.Marshal(el)
			if err != nil {
				p.errorf("%v", err)
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

func (p *Parser) parseAnnotation(annot string) ast.Annotation {
	if len(annot) == 0 {
		p.errorf("Unexpected empty annotation.")
	}

	var annotationKind ast.AnnotationKind

	switch annot[0] {
	case ':':
		annotationKind = ast.TypeAnnotation
	case '@':
		annotationKind = ast.VariableAnnotation
	case '%':
		annotationKind = ast.FieldAnnotation
	default:
		p.errorf("Unexpected annotation: (%s)", annot)
	}

	return ast.Annotation{
		Kind:  annotationKind,
		Value: annot,
	}
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
