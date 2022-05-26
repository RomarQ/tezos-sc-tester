package micheline

import (
	"fmt"
	"strings"

	"github.com/romarq/tezos-sc-tester/internal/business/michelson/ast"
	"github.com/romarq/tezos-sc-tester/internal/business/michelson/utils"
)

type (
	context struct {
		indent string
		depth  int
	}
)

func Print(n ast.Node, indent string) string {
	p := context{
		indent: indent,
		depth:  0,
	}
	return p.print(n, true)
}

func (p *context) print(n ast.Node, allowParenthesis bool) (micheline string) {
	switch node := n.(type) {
	case ast.Bytes:
		micheline = fmt.Sprintf("0x%s", node.Value)
	case ast.Int:
		micheline = node.Value
	case ast.String:
		micheline = fmt.Sprintf(`"%s"`, node.Value)
	case ast.Prim:
		micheline = p.printPrim(node, allowParenthesis)
	case ast.Sequence:
		micheline = p.printSequence(node)
	}
	return
}

func (p *context) printPrim(n ast.Prim, allowParenthesis bool) string {
	args := []string{n.Prim}
	for _, el := range n.Annotations {
		args = append(args, el.Value)
	}
	for _, el := range n.Arguments {
		args = append(args, p.print(el, true))
	}

	if !allowParenthesis || utils.IsInstruction(n.Prim) || utils.IsReservedWord(n.Prim) {
		return fmt.Sprintf("%s", strings.Join(args, " "))
	}

	return fmt.Sprintf("(%s)", strings.Join(args, " "))
}

func (p *context) printSequence(n ast.Sequence) string {
	prevIndent := p.getIndent()

	// Increase depth used for indentation
	p.depth += 1
	defer func() {
		// Decrease depth on function exit
		p.depth -= 1
	}()

	elements := make([]string, 0)
	for _, el := range n.Elements {
		elements = append(elements, p.print(el, false))
	}

	return fmt.Sprintf("{%s%s%s}", p.getIndent(), strings.Join(elements, fmt.Sprintf(";%s", p.getIndent())), prevIndent)
}

func (p *context) getIndent() string {
	if p.indent != "" {
		return fmt.Sprintf("\n%s", strings.Repeat(p.indent, p.depth))
	} else {
		return " "
	}

}
