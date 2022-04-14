package micheline

import (
	"fmt"
	"unicode"

	"github.com/romarq/visualtez-testing/internal/business/michelson/micheline/token"
)

type (
	Scanner struct {
		// immutable state
		source string

		// scanning state
		char         rune // current character
		offset       int  // character offset
		rdOffset     int  // reading offset (position after current character)
		lineOffset   int  // current line offset
		columnOffset int  // current column offset

		panicOnError bool
		errors       []Error
	}
)

const (
	NUL = 0
)

func (s *Scanner) Init(micheline string) {
	s.source = micheline
	s.offset = 0
	s.rdOffset = 0
	s.lineOffset = 0
	s.columnOffset = 0
}

func (s *Scanner) Scan() (pos int, tk token.Kind, text string) {
	s.skipWhitespace()

	s.next()
	pos = s.offset
	text = string(s.char)

	switch s.char {
	case NUL:
		tk = token.Nul
	case '{':
		tk = token.Open_brace
	case '}':
		tk = token.Close_brace
	case '(':
		tk = token.Open_paren
	case ')':
		tk = token.Close_paren
	case ';':
		tk = token.Semi
	case '#':
		tk = token.Comment
		for s.peek() != '\n' {
			if s.IsAtEnd() {
				s.errorf(`Reached EOF while parsing comment. %s`, s.lineInfo())
				break
			}
			s.next()
			text = text + string(s.char)
		}
	case '%', ':', '@':
		tk = token.Annot
		for isIdentifier(s.peek()) {
			if s.IsAtEnd() {
				s.errorf(`Reached EOF while parsing annotation. %s`, s.lineInfo())
				break
			}
			s.next()
			text = text + string(s.char)
		}
	case '"': // String
		tk = token.String
		text = ""
		for s.peek() != '"' {
			if s.IsAtEnd() {
				s.errorf(`Reached EOF while parsing a string. %s`, s.lineInfo())
				break
			}
			s.next()
			text = text + string(s.char)
		}
		s.next() // Consume closing quote
	case '0': // Can be an integer or bytes
		for isIdentifier(s.peek()) {
			if s.IsAtEnd() {
				s.errorf(`Reached EOF while parsing hexadecimal. %s`, s.lineInfo())
				break
			}
			s.next()
			text = text + string(s.char)
		}

		if len(text) > 1 && text[0:2] == "0x" {
			tk = token.Bytes
		} else {
			tk = token.Int
		}
	default:
		for isIdentifier(s.peek()) {
			if s.IsAtEnd() {
				s.errorf(`Reached EOF while parsing value. %s`, s.lineInfo())
				break
			}
			s.next()
			text = text + string(s.char)
		}

		if unicode.IsDigit(s.char) {
			tk = token.Int
		} else {
			tk = token.Identifier
		}
	}

	return
}

func (s *Scanner) skipWhitespace() {
	for isWhitespace(s.peek()) {
		s.incrementPosition()
	}
}

func (s *Scanner) incrementLine()   { s.lineOffset = s.lineOffset + 1 }
func (s *Scanner) incrementColumn() { s.columnOffset = s.columnOffset + 1 }
func (s *Scanner) resetColumn()     { s.columnOffset = 0 }
func (s *Scanner) incrementPosition() {
	if s.char == '\n' {
		s.incrementLine()
		s.resetColumn()
	} else {
		s.incrementColumn()
	}

	s.offset = s.rdOffset
	s.rdOffset = s.rdOffset + 1
}

// next reads the next byte in the scanner sequence and increments the current position.
func (s *Scanner) next() {
	if s.char = s.peek(); s.char != NUL {
		s.incrementPosition()
	}
}

// peek returns the next byte in the scanner sequence without
// incrementing the current position. Returns (0 => NUL) if the scanning is over.
func (s Scanner) peek() rune {
	if s.IsAtEnd() {
		return NUL
	}
	return rune(s.source[s.rdOffset])
}

func (s Scanner) IsAtEnd() bool {
	// len(s.source) is cached and there is no overhead
	return s.rdOffset == len(s.source)
}

func (s Scanner) lineInfo() string {
	return fmt.Sprintf("(line: %d, column: %d)", s.lineOffset, s.columnOffset)
}

func isWhitespace(c rune) bool { return c == ' ' || c == '\n' || c == '\r' || c == '\t' }
func isIdentifier(c rune) bool { return '0' <= c && c <= '9' || c == '_' || 'A' <= c && c <= 'z' }

func (s *Scanner) errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if s.panicOnError {
		panic(msg)
	}
	s.errors = append(s.errors, Error{Message: msg})
}
