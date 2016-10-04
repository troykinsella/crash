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
	r *bufio.Reader
	col uint16
	line uint16
}

func New(r *bufio.Reader) *Scanner {
	return &Scanner{
		r: r,
		col: 0,
		line: 0,
	}
}

func (s *Scanner) Pos() (uint16, uint16) {
	return s.line, s.col
}

func (s *Scanner) Read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	s.col++
	return ch
}

func (s *Scanner) Unread() {
	s.col--
	s.r.UnreadRune()
}

func (s *Scanner) Peek() rune {
	ch := s.Read()
	s.Unread()
	return ch
}

func (s *Scanner) consumeWhitespace() string {
	var buf bytes.Buffer
	for {
		ch := s.Read()
		if ch == eof {
			break
		}
		if !isWhitespace(ch) {
			s.Unread()
			break
		}
		buf.WriteRune(ch)
	}
	return buf.String()
}

func (s *Scanner) consumeComment() token.Token {
	s.Read()

	if ch := s.Read(); ch != '/' {
		panic(fmt.Errorf("invalid comment: %s", ch))
	}

	// Skip initial whitespace
	for {
		ch := s.Read()
		if ch != ' ' || ch == '\t' {
			s.Unread()
			break
		}
	}

	return token.COMMENT
}

func (s *Scanner) consumeString() (token.Token, string) {
	var buf bytes.Buffer

	startQuote := s.Read()

	for {
		ch := s.Read()
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
	ch := s.Read()
	dollar := ch == '$'
	buf.WriteRune(ch)

	for {
		ch = s.Read();
		if ch == eof {
			break
		}
		if (dollar && isDigit(ch)) || isIdentifier(ch) {
			buf.WriteRune(ch)
		} else {
			s.Unread()
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
	buf.WriteRune(s.Read())

	for {
		ch := s.Read();
		if ch == eof {
			break
		}
		if !isDigit(ch) {
			s.Unread()
			break
		}
		buf.WriteRune(ch)
	}

	return token.NUMBER, buf.String()
}

func (s *Scanner) Scan() (token.Token, string) {
	ch := s.Peek()
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

	ch = s.Read()

	switch ch {
	case '(': return token.OPEN_BRACKET, ""
	case ')': return token.CLOSE_BRACKET, ""
	case '[': return token.OPEN_SQUARE_BRACKET, ""
	case ']': return token.CLOSE_SQUARE_BRACKET, ""
	case ',': return token.COMMA, ""
	case '.': return token.DOT, ""
	case '}': return token.INTERPOLATE_END, ""
	}

	if ch == '$' {
		ch = s.Read()
		if ch == '{' {
			return token.INTERPOLATE_BEGIN, ""
		}
		s.Unread()
	}

	return token.ILLEGAL, string(ch)
}

func (s *Scanner) SeekInterp() (token.Token, string) {
	var buf bytes.Buffer

	for {
		ch := s.Read()
		if ch == eof {
			return token.EOF, buf.String()
		}
		if ch == '$' {
			if s.Peek() == '{' {
				s.Read()
				return token.INTERPOLATE_BEGIN, buf.String()
			}
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
