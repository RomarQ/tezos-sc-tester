package micheline

import (
	"fmt"
	"strings"

	"github.com/romarq/visualtez-testing/internal/business/michelson/ast"
	"github.com/romarq/visualtez-testing/internal/business/michelson/utils"
)

func Print(n ast.Node) (micheline string) {
	switch node := n.(type) {
	case ast.Bytes:
		micheline = node.Value
	case ast.Int:
		micheline = fmt.Sprint(node.Value)
	case ast.String:
		micheline = fmt.Sprintf(`"%s"`, node.Value)
	case ast.Prim:
		micheline = printPrim(node)
	case ast.Sequence:
		micheline = printSequence(node)
	}
	return
}

func printPrim(n ast.Prim) string {
	args := []string{n.Prim}
	for _, el := range n.Annotations {
		args = append(args, el.Value)
	}
	for _, el := range n.Arguments {
		args = append(args, Print(el))
	}

	if utils.IsInstruction(n.Prim) || utils.IsReservedWord(n.Prim) {
		return fmt.Sprintf("%s", strings.Join(args, " "))
	}

	return fmt.Sprintf("(%s)", strings.Join(args, " "))
}

func printSequence(n ast.Sequence) string {
	elements := make([]string, 0)
	for _, el := range n.Elements {
		elements = append(elements, Print(el))
	}

	return fmt.Sprintf("{ %s }", strings.Join(elements, " ; "))
}
