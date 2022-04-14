package token

import "fmt"

type (
	Kind  int8
	Token interface {
		Kind() Kind
		Text() string
	}
	token struct {
		kind Kind
		text string
	}
)

const (
	String Kind = iota
	Bytes
	Int
	Identifier
	Annot
	Comment
	Semi
	Open_paren
	Close_paren
	Open_brace
	Close_brace
	Nul
)

// String returns the string corresponding to the token kind.
func (k Kind) String() string {
	switch k {
	case String:
		return "String"
	case Bytes:
		return "Bytes"
	case Int:
		return "Int"
	case Identifier:
		return "Identifier"
	case Annot:
		return "Annot"
	case Comment:
		return "Comment"
	case Semi:
		return "Semi"
	case Open_paren:
		return "Open_paren"
	case Close_paren:
		return "Close_paren"
	case Open_brace:
		return "Open_brace"
	case Close_brace:
		return "Close_brace"
	case Nul:
		return "NUL"
	}

	panic(fmt.Sprintf("Unexpected token kind (%d)", k))
}

func (t token) Kind() Kind {
	return t.kind
}
func (t token) Text() string {
	return t.text
}
