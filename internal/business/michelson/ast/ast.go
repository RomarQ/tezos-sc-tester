package ast

import (
	"fmt"
	"strings"
)

type Node interface {
	String() string
}

type Position struct {
	Pos int
	End int
}

type Int struct {
	Position
	Value string
}

type String struct {
	Position
	Value string
}

type Bytes struct {
	Position
	Value string
}

type Sequence struct {
	Position
	Elements []Node // list of elements in the sequence
}

type Prim struct {
	Position
	Prim        string
	Annotations []Annotation
	Arguments   []Node
}

type (
	AnnotationKind uint8
	Annotation     struct {
		Position
		Kind  AnnotationKind // Annotation kind (https://tezos.gitlab.io/alpha/michelson.html#annotations)
		Value string
	}
)

const (
	TypeAnnotation AnnotationKind = iota
	VariableAnnotation
	FieldAnnotation
)

func (n Bytes) String() string  { return fmt.Sprintf("Bytes(%s)", n.Value) }
func (n String) String() string { return fmt.Sprintf("String(%s)", n.Value) }
func (n Int) String() string    { return fmt.Sprintf("Int(%s)", n.Value) }
func (n Prim) String() string {
	annotations := make([]string, 0)
	for _, el := range n.Annotations {
		annotations = append(annotations, el.Value)
	}
	args := make([]string, 0)
	for _, el := range n.Arguments {
		args = append(args, el.String())
	}
	return fmt.Sprintf("Prim(%s, [%s], [%s])", n.Prim, strings.Join(annotations, ", "), strings.Join(args, ", "))
}
func (n Sequence) String() string {
	elements := make([]string, 0)
	for _, el := range n.Elements {
		elements = append(elements, el.String())
	}
	return fmt.Sprintf("Sequence([%s])", strings.Join(elements, ", "))
}
