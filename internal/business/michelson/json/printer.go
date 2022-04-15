package json

import (
	"encoding/json"
	"fmt"

	"github.com/romarq/visualtez-testing/internal/business/michelson/ast"
)

func Print(n ast.Node, prefix string, indent string) (json.RawMessage, error) {
	switch node := n.(type) {
	case ast.Bytes:
		return MichelsonJSON{
			Bytes: node.Value[2:],
		}.Marshal(prefix, indent)
	case ast.Int:
		return MichelsonJSON{
			Int: fmt.Sprint(node.Value),
		}.Marshal(prefix, indent)
	case ast.String:
		return MichelsonJSON{
			String: node.Value,
		}.Marshal(prefix, indent)
	case ast.Prim:
		return printPrim(node, prefix, indent)
	case ast.Sequence:
		return printSequence(node, prefix, indent)
	}
	fmt.Print(n)
	return nil, fmt.Errorf("Unexpected AST Node (%s).", n.String())
}

func printPrim(n ast.Prim, prefix string, indent string) (json.RawMessage, error) {
	prim := MichelsonJSON{
		Prim: n.Prim,
	}
	for _, el := range n.Annotations {
		prim.Annots = append(prim.Annots, el.Value)
	}
	for _, el := range n.Arguments {
		b, err := Print(el, prefix, indent)
		if err != nil {
			return nil, err
		}
		prim.Args = append(prim.Args, b)
	}

	return prim.Marshal(prefix, indent)
}

func printSequence(n ast.Sequence, prefix string, indent string) (json.RawMessage, error) {
	sequence := make([]json.RawMessage, 0)
	for _, el := range n.Elements {
		b, err := Print(el, prefix, indent)
		if err != nil {
			return nil, err
		}
		sequence = append(sequence, b)
	}

	return json.MarshalIndent(sequence, prefix, indent)
}
