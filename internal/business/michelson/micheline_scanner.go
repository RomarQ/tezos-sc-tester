package michelson

import (
	"fmt"
)

type (
	MichelineScanner struct {
		position     uint
		row          uint
		column       uint
		sourceLength uint
		source       string
	}
	token struct {
		Kind tokenKind
		Text string
	}
	tokenKind int8
)

const (
	String tokenKind = iota
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
)

func InitMichelineScanner(micheline string) MichelineScanner {
	return MichelineScanner{
		position:     0,
		row:          0,
		column:       0,
		sourceLength: uint(len(micheline)),
		source:       micheline,
	}
}

func (s *MichelineScanner) Tokenize() ([]token, error) {
	tokens := make([]token, 0)
	for !s.isEnd() {
		c := s.next()

		switch c {
		default:
			text := c
			if isDigit(c) {
				for isDigit(s.lookNext()) {
					if s.isEnd() {
						return tokens, fmt.Errorf(`Reached EOF while parsing integer. %s`, s.line())
					}
					text = text + s.next()
				}
				tokens = append(tokens, buildToken(Int, text))
			} else if isIdentifier(c) {
				for isIdentifier(s.lookNext()) {
					if s.isEnd() {
						return tokens, fmt.Errorf(`Reached EOF while parsing identifier. %s`, s.line())
					}
					text = text + s.next()
				}
				tokens = append(tokens, buildToken(Identifier, text))
			}
		case " ", "\r", "\t":
		case "\n":
			s.incrementRow()
			s.resetColumn()
		case "{":
			tokens = append(tokens, buildToken(Open_brace, c))
		case "}":
			tokens = append(tokens, buildToken(Close_brace, c))
		case "(":
			tokens = append(tokens, buildToken(Open_paren, c))
		case ")":
			tokens = append(tokens, buildToken(Close_paren, c))
		case ";":
			tokens = append(tokens, buildToken(Semi, c))
		case "#":
			text := c
			for s.lookNext() != "\n" {
				if s.isEnd() {
					return tokens, fmt.Errorf(`Reached EOF while parsing comment. %s`, s.line())
				}
				text = text + s.next()
			}
			tokens = append(tokens, buildToken(Comment, text))
		case "%", ":", "@":
			text := c
			for isIdentifier(s.lookNext()) {
				if s.isEnd() {
					return tokens, fmt.Errorf(`Reached EOF while parsing annotation. %s`, s.line())
				}
				text = text + s.next()
			}
			tokens = append(tokens, buildToken(Annot, text))
		case "\"": // String
			text := ""
			c = s.next()
			for c != "\"" {
				if s.isEnd() {
					return tokens, fmt.Errorf(`Reached EOF while parsing a string. %s`, s.line())
				}
				text = text + c
				c = s.next()
			}
			tokens = append(tokens, buildToken(String, text))
		case "0": // Can be an integer or bytes
			text := c
			for isIdentifier(s.lookNext()) {
				if s.isEnd() {
					return tokens, fmt.Errorf(`Reached EOF while parsing hexadecimal. %s`, s.line())
				}
				text = text + s.next()
			}

			if len(text) > 1 && text[0:2] == "0x" {
				// Validate hexadecimal string
				byteString := text[2:]
				if !isHex(byteString) {
					return tokens, fmt.Errorf(`Invalid hexadecimal value (%s). %s`, byteString, s.line())
				}
				tokens = append(tokens, buildToken(Bytes, byteString))
			} else {
				// Validate integer
				if !isDigit(text) {
					return tokens, fmt.Errorf(`Invalid integer value (%s). %s`, text, s.line())
				}
				tokens = append(tokens, buildToken(Int, text))
			}
		}
	}

	return tokens, nil
}

func buildToken(kind tokenKind, text string) token {
	return token{
		Kind: kind,
		Text: text,
	}
}

func (s *MichelineScanner) incrementRow() {
	s.row = s.row + 1
}
func (s *MichelineScanner) incrementColumn() {
	s.column = s.column + 1
}
func (s *MichelineScanner) resetColumn() {
	s.column = 0
}
func (s *MichelineScanner) incrementPosition() {
	s.position = s.position + 1
	s.incrementColumn()
}

func (s *MichelineScanner) next() string {
	if s.isEnd() {
		return ""
	}

	// Increment position
	s.incrementPosition()

	// Get next character
	return string(s.source[s.position-1])
}

func (s MichelineScanner) lookNext() string {
	if s.isEnd() {
		return ""
	}
	// Get next character without incrementing position
	return string(s.source[s.position])
}

func (s MichelineScanner) isEnd() bool {
	return s.position == s.sourceLength
}

func (s MichelineScanner) line() string {
	return fmt.Sprintf("(line: %d, column: %d)", s.row, s.column)
}
