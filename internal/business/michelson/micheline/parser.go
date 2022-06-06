package micheline

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/romarq/tezos-sc-tester/internal/business/michelson/ast"
	"github.com/romarq/tezos-sc-tester/internal/business/michelson/micheline/token"
)

type Parser struct {
	token_position int
	token_kind     token.Kind
	token_text     string

	scanner Scanner

	trace bool
}

var (
	regex_bytes  = regexp.MustCompile("^0x[0-9a-fA-F]+$")
	regex_number = regexp.MustCompile("^-?[0-9]+$")
)

func InitParser(micheline string) (parser Parser) {
	parser.scanner = InitScanner(micheline)
	return
}

func (p *Parser) Parse() ast.Node {
	p.next()

	switch kind := p.token_kind; {
	case kind == token.Bytes:
		return p.parseBytes()
	case kind == token.String:
		return p.parseString()
	case kind == token.Int:
		return p.parseInt()
	case kind == token.Open_paren:
		return p.parseParenthesis()
	case kind == token.Identifier:
		return p.parsePrim()
	case kind == token.Open_brace:
		return p.parseSequence()
	default:
		p.scanner.errorf("Unexpected token (%s) as sequence child.", kind.String())
	}

	return nil
}

func (p *Parser) HasErrors() bool {
	return len(p.scanner.errors) > 0
}

func (p *Parser) Error() error {
	if !p.HasErrors() {
		return nil
	}
	_errors := make([]string, len(p.scanner.errors))
	for _, err := range p.scanner.errors {
		_errors = append(_errors, err.Message)
	}
	return errors.New(strings.Join(_errors, ";\n"))
}

func (p *Parser) next() {
	p.token_position, p.token_kind, p.token_text = p.scanner.Scan()
	if p.trace {
		fmt.Printf("[Scanner] (%s) with text (%s)\n", p.token_kind.String(), p.token_text)
	}
}

func (p *Parser) parseBytes() ast.Bytes {
	if p.trace {
		fmt.Println("[Parsing|IN] Bytes")
		defer fmt.Println("[Parsing|OUT] Bytes")
	}

	position := p.expect(token.Bytes)
	defer p.next() // Consume next token

	bytes := p.token_text
	if isBytes(p.token_text) || len(p.token_text)%2 == 0 {
		bytes = p.token_text[2:]
	} else {
		p.scanner.errorf("Invalid bytes: %s. %v", p.token_text, position)
	}

	return ast.Bytes{
		Position: ast.Position{
			Pos: position,
			End: position + len(p.token_text) - 1,
		},
		Value: bytes,
	}
}

func (p *Parser) parseString() ast.String {
	if p.trace {
		fmt.Println("[Parsing|IN] String")
		defer fmt.Println("[Parsing|OUT]  String")
	}

	position := p.expect(token.String)
	defer p.next() // Consume next token

	return ast.String{
		Position: ast.Position{
			Pos: position,
			End: position + len(p.token_text) + /* Count quotes */ 1,
		},
		Value: p.token_text,
	}
}

func (p *Parser) parseInt() ast.Int {
	if p.trace {
		fmt.Println("[Parsing|IN] Int")
		defer fmt.Println("[Parsing|OUT] Int")
	}

	position := p.expect(token.Int)
	defer p.next() // Consume next token

	if !isNumber(p.token_text) {
		p.scanner.errorf("Invalid number: %s. %v", p.token_text, position)
	}

	return ast.Int{
		Position: ast.Position{
			Pos: position,
			End: position + len(p.token_text) - 1,
		},
		Value: p.token_text,
	}
}

func (p *Parser) parseSequence() ast.Sequence {
	if p.trace {
		fmt.Println("[Parsing|IN] Sequence")
		defer fmt.Println("[Parsing|OUT] Sequence")
	}

	begin := p.expect(token.Open_brace)
	p.next() // Consume next token

	elements := make([]ast.Node, 0)
	for p.token_kind != token.Close_brace {
		switch p.token_kind {
		case token.Bytes:
			elements = append(elements, p.parseBytes())
		case token.String:
			elements = append(elements, p.parseString())
		case token.Int:
			elements = append(elements, p.parseInt())
		case token.Identifier:
			elements = append(elements, p.parsePrim())
		case token.Open_brace:
			elements = append(elements, p.parseSequence())
		default:
			p.scanner.errorf("Unexpected token (%s) as sequence child.", p.token_kind.String())
		}

		if p.token_kind != token.Close_brace {
			p.expect(token.Semi) // Semicolon is used to separate elements in sequences
			p.next()             // Consume next token
		}
	}
	end := p.expect(token.Close_brace)
	defer p.next() // Consume next token

	// TODO: Sort sequences of comparable values
	// Michelson enforces maps and sets to be sorted

	return ast.Sequence{
		Position: ast.Position{
			Pos: begin,
			End: end,
		},
		Elements: elements,
	}
}

func (p *Parser) parsePrim() ast.Prim {
	if p.trace {
		fmt.Printf("[Parsing|IN] Prim (%s)\n", p.token_text)
		defer fmt.Printf("[Parsing|OUT] Prim (%s)\n", p.token_text)
	}

	begin := p.expect(token.Identifier)
	identifier := p.token_text

	p.next() // Consume next token

	// Check annotations (Annotations can only appear right after an identifier)
	annotations := p.parseAnnotations()

	arguments := make([]ast.Node, 0)

	for {
		switch p.token_kind {
		case token.Bytes:
			arguments = append(arguments, p.parseBytes())
			continue
		case token.String:
			arguments = append(arguments, p.parseString())
			continue
		case token.Int:
			arguments = append(arguments, p.parseInt())
			continue
		case token.Open_paren:
			arguments = append(arguments, p.parseParenthesis())
			continue
		case token.Identifier:
			for p.token_kind == token.Identifier {
				identBegin := p.token_position
				arguments = append(arguments, ast.Prim{
					Position: ast.Position{
						Pos: identBegin,
						End: identBegin + len(p.token_text) - 1,
					},
					Prim:        p.token_text,
					Annotations: p.parseAnnotations(),
				})
				p.next() // Consume next token
			}
			continue
		case token.Open_brace:
			arguments = append(arguments, p.parseSequence())
			continue
		}
		break
	}

	return ast.Prim{
		Position: ast.Position{
			Pos: begin,
			End: p.token_position,
		},
		Prim:        identifier,
		Annotations: annotations,
		Arguments:   arguments,
	}
}

func (p *Parser) parseParenthesis() ast.Prim {
	begin := p.expect(token.Open_paren)

	p.next() // Consume next token
	if p.token_kind != token.Identifier {
		p.scanner.errorf("Expected token (%s), but received (%s).", token.Identifier.String(), p.token_kind.String())
	}

	node := p.parsePrim()
	end := p.expect(token.Close_paren)
	p.next() // Consume next token

	node.Position.Pos = begin
	node.Position.End = end

	return node
}

// parseAnnotations parses annotations (Annotations can only appear right after an identifier)
func (p *Parser) parseAnnotations() (annotations []ast.Annotation) {
	for p.token_kind == token.Annot {
		annotations = append(annotations, p.parseAnnotation())
		p.next() // Consume next token
	}
	return
}

func (p *Parser) parseAnnotation() ast.Annotation {
	position := p.expect(token.Annot)

	var annotationKind ast.AnnotationKind

	if len(p.token_text) == 0 {
		p.scanner.errorf("Unexpected empty")
	} else {
		switch p.token_text[0] {
		case ':':
			annotationKind = ast.TypeAnnotation
		case '@':
			annotationKind = ast.VariableAnnotation
		case '%':
			annotationKind = ast.FieldAnnotation
		default:
			p.scanner.errorf("Unexpected annotation: (%s)", p.token_text)
		}
	}
	return ast.Annotation{
		Position: ast.Position{
			Pos: position,
			End: position + len(p.token_text) - 1,
		},
		Kind:  annotationKind,
		Value: p.token_text,
	}
}

func (p *Parser) expect(kind token.Kind) (pos int) {
	if p.token_kind == kind {
		pos = p.token_position
	} else {
		p.scanner.errorf("Expected token kind (%s), but received (%s).", kind.String(), p.token_kind.String())
	}
	return
}

func isBytes(text string) bool  { return regex_bytes.MatchString(text) }
func isNumber(text string) bool { return regex_number.MatchString(text) }
