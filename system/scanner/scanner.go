package scanner

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/troykinsella/crash/system/token"
	"unicode/utf8"
	"unicode"
)

var eof = rune(0)


type Scanner struct {
	in *RuneReader
	col uint16
	line uint16
}

func New(in *bufio.Reader) *Scanner {
	return &Scanner{
		in: NewRuneReader(in),
		col: 0,
		line: 0,
	}
}

func (s *Scanner) Pos() (uint16, uint16) {
	return s.line, s.col
}

func (s *Scanner) read() rune {
	ch, err := s.in.Read()
	if err != nil {
		return eof
	}
	s.col++
	return ch
}

func (s *Scanner) unread() {
	s.col--
	s.in.Unread()
}

func (s *Scanner) peek() rune {
	ch := s.read()
	if ch != eof {
		s.unread()
	}
	return ch
}

func (s *Scanner) consumeWhitespace() string {
	var buf bytes.Buffer
	for {
		ch := s.read()
		if ch == eof {
			break
		}
		if !isWhitespace(ch) {
			s.unread()
			break
		}
		buf.WriteRune(ch)
	}
	return buf.String()
}

func (s *Scanner) consumeComment() token.Token {
	s.read()

	if ch := s.read(); ch != '/' {
		return token.ILLEGAL
	}

	// Skip initial whitespace
	/*for {
		ch := s.Read()
		if ch != ' ' && ch != '\t' {
			break
		}
	}*/

	return token.COMMENT
}

func (s *Scanner) consumeString() (token.Token, string) {
	var buf bytes.Buffer

	startQuote := s.read()

	for {
		ch := s.read()
		if ch == eof {
			panic(fmt.Errorf("Unexpected end of file: %s", buf.String()))
		}
		if ch == startQuote {
			break
		}
		buf.WriteRune(ch)
	}

	return token.STRING, buf.String()
}

func (s *Scanner) consumeIdentifier() (token.Token, string) {
	var buf bytes.Buffer
	ch := s.read()
	buf.WriteRune(ch)

	for {
		ch = s.read();
		if ch == eof {
			break
		}
		if isIdentifier(ch) {
			buf.WriteRune(ch)
		} else {
			s.unread()
			break
		}
	}

	str := buf.String()

	switch str {
	case "true":
		return token.TRUE, ""
	case "false":
		return token.FALSE, ""
	case "and":
		return token.AND, ""
	case "not":
		return token.NOT, ""
	case "or":
		return token.OR, ""
	case "xor":
		return token.XOR, ""
	}

	return token.IDENT, str
}

func (s *Scanner) consumeNumber() (token.Token, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		ch := s.read();
		if ch == eof {
			break
		}
		if !isDigit(ch) {
			s.unread()
			break
		}
		buf.WriteRune(ch)
	}

	return token.NUMBER, buf.String()
}

func (s *Scanner) Unscan() {
	s.in.Rewind()
}

func (s *Scanner) Scan() (token.Token, string) {
	s.in.Reset()

	ch := s.peek()
	if ch == eof {
		return token.EOF, ""
	}

	if isWhitespace(ch) {
		return token.WS, s.consumeWhitespace()
	}
	if isIdentifier(ch) {
		return s.consumeIdentifier()
	}
	if isString(ch) {
		return s.consumeString()
	}
	if isDigit(ch) {
		return s.consumeNumber()
	}
	if ch == '/' {
		return s.consumeComment(), ""
	}

	ch = s.read()

	switch ch {
	case '(': return token.OPEN_BRACKET, ""
	case ')': return token.CLOSE_BRACKET, ""
	case '[': return token.OPEN_SQUARE_BRACKET, ""
	case ']': return token.CLOSE_SQUARE_BRACKET, ""
	case ',': return token.COMMA, ""
	case '.': return token.DOT, ""
	case '}': return token.INTERPOLATE_END, ""
	case '\\': return token.ESCAPE, ""
	}

	if ch == '$' {
		lit := "$"
		ch = s.peek()
		if ch == '{' {
			lit += "{"
			s.read()
		}
		return token.INTERPOLATE_BEGIN, lit
	}

	return token.ILLEGAL, string(ch)
}

// return bool is "expect interpolate end"
func (s *Scanner) SeekInterp() (token.Token, bool, string) {
	var buf bytes.Buffer

	escape := false
	close := false
	for {
		ch := s.read()
		if ch == eof {
			return token.EOF, false, buf.String()
		}

		if ch == '\\' {
			escape = true
		} else if ch == '$' && !escape {
			p := s.peek()

			if p == '{' {
				s.read()
				close = true
			}

			// Special case: '$' before EOF -> don't bother returning interpolate_begin
			p = s.peek()
			if p == eof {
				buf.WriteRune('$')
				if close {
					buf.WriteRune('{')
				}
				return token.EOF, false, buf.String()
			}

			return token.INTERPOLATE_BEGIN, close, buf.String()
		} else {
			escape = false
		}

		buf.WriteRune(ch)
	}
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9') ||
	       (ch >= utf8.RuneSelf && unicode.IsDigit(ch))
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') ||
	       (ch >= 'A' && ch <= 'Z') ||
	       (ch >= utf8.RuneSelf && unicode.IsLetter(ch))
}

func isIdentifier(ch rune) bool {
	return isLetter(ch) || ch == '_' || ch == '-'
}

func isString(ch rune) bool {
	return ch == '"' || ch == '\'' || ch == '`'
}
