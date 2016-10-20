package scanner

import (
	"testing"
	"strings"
	"bufio"
	"github.com/troykinsella/crash/system/token"
	"reflect"
)

func boolEq(actual bool, expected bool, t *testing.T) {
	if actual != expected {
		t.Errorf("Unexpected bool: expected=%t actual=%t", expected, actual)
	}
}

func tokEq(actual token.Token, expected token.Token, t *testing.T) {
	if actual != expected {
		t.Errorf("Unexpected token: expected=%s actual=%s", expected.String(), actual.String())
	}
}

func strEq(actual string, expected string, t *testing.T) {
	if actual != expected {
		t.Errorf("Unexpected string: expected=%s actual=%s", expected, actual)
	}
}

func TestScanner_Scan(t *testing.T) {
	var tests = []struct {
		str string
		tok token.Token
		lit string
	}{
		{"",       token.EOF,                  ""},
		{" ",      token.WS,                   " "},
		{"\t",     token.WS,                   "\t" },
		{"\r",     token.WS,                   "\r" },
		{"\n",     token.WS,                   "\n" },
		{"foo",    token.IDENT,                "foo"},
		{"1",      token.NUMBER,               "1"},
		{"123",    token.NUMBER,               "123"},
		{"'foo'",  token.STRING,               "foo"},
		{`"foo"`,  token.STRING,               "foo"},
		{"`foo`",  token.STRING,               "foo"},
		{",",      token.COMMA,                ""},
		{"(",      token.OPEN_BRACKET,         ""},
		{")",      token.CLOSE_BRACKET,        ""},
		{"[",      token.OPEN_SQUARE_BRACKET,  ""},
		{"]",      token.CLOSE_SQUARE_BRACKET, ""},
		{"//",     token.COMMENT,              ""},
		{".",      token.DOT,                  ""},
		{"true",   token.TRUE,                 ""},
		{"false",  token.FALSE,                ""},
		{"and",    token.AND,                  ""},
		{"or",     token.OR,                   ""},
		{"xor",    token.XOR,                  ""},
		{"not",    token.NOT,                  ""},
		{"$",      token.INTERPOLATE_BEGIN,    "$"},
		{"${",     token.INTERPOLATE_BEGIN,    "${"},
		{"}",      token.INTERPOLATE_END,      ""},
	}

	for i, test := range tests {
		s := New(bufio.NewReader(strings.NewReader(test.str)))
		tok, lit := s.Scan()
		if tok != test.tok {
			t.Errorf("%d. unexpected token: expected=%q actual=%q (%q)", i, test.tok, tok, lit)
		} else if lit != test.lit {
			t.Errorf("%d. unexpected literal: expected=%q actual=%q", i, test.lit, lit)
		}

		if tok, lit = s.Scan(); tok != token.EOF {
			t.Errorf("%d. expected eof: actual=%q", i, lit)
		}
	}
}

func TestScanner_SeekInterp(t *testing.T) {
	var tests = []struct {
		str string
		tok token.Token
		text string
		expectEnd bool
	}{
		{"",   token.EOF, "", false},
		{"$",  token.EOF, "$", false},
		{" $",  token.EOF, " $", false},
		{"${", token.EOF, "${", false},
		{" ${", token.EOF, " ${", false},
		{"asdf$", token.EOF, "asdf$", false},
		{"asdf${", token.EOF, "asdf${", false},

		{"asdf",   token.EOF, "asdf", false},
		{"$asdf",  token.INTERPOLATE_BEGIN, "", false},
		{"${asdf", token.INTERPOLATE_BEGIN, "", true},
		{" $asdf", token.INTERPOLATE_BEGIN, " ", false},
		{" ${asdf", token.INTERPOLATE_BEGIN, " ", true},
		{"asdf$fdsa", token.INTERPOLATE_BEGIN, "asdf", false},
		{"asdf${fdsa", token.INTERPOLATE_BEGIN, "asdf", true},

		{"\\$",    token.EOF, "\\$", false},
		{"\\${",   token.EOF, "\\${", false},
		{" \\${",  token.EOF, " \\${", false},

		{"$$", token.INTERPOLATE_BEGIN, "", false},
		{"${$", token.INTERPOLATE_BEGIN, "", true},
		{"{$", token.EOF, "{$", false},
		{"{${", token.EOF, "{${", false},
		{"$$asdf", token.INTERPOLATE_BEGIN, "", false},
		{"${$asdf", token.INTERPOLATE_BEGIN, "", true},
		{"{$asdf", token.INTERPOLATE_BEGIN, "{", false},
		{"{${asdf", token.INTERPOLATE_BEGIN, "{", true},
	}

	for i, test := range tests {
		s := New(bufio.NewReader(strings.NewReader(test.str)))
		tok, expectEnd, text := s.SeekInterp()
		if tok != test.tok {
			t.Errorf("%d. unexpected token: expected=%q actual=%q text=%q", i, test.tok, tok, text)
		} else if text != test.text {
			t.Errorf("%d. unexpected text: expected=%q actual=%q", i, test.text, text)
		} else if expectEnd != test.expectEnd {
			t.Errorf("%d. unexpected expectEnd: expected=%t actual=%t", i, test.expectEnd, expectEnd)
		}
	}
}

func TestScanner_tokenStream(t *testing.T) {
	var tests = []struct {
		str string
		toks []token.Token
		lits []string
	}{
		{
			"foo op // bar",
			[]token.Token{
				token.IDENT,
				token.WS,
				token.IDENT,
				token.WS,
				token.COMMENT,
				token.WS,
				token.IDENT,
			},
			[]string{
				"foo",
				" ",
				"op",
				" ",
				"",
				" ",
				"bar",
			},
		},
	}

	for i, test := range tests {
		s := New(bufio.NewReader(strings.NewReader(test.str)))

		toks := make([]token.Token, 0)
		lits := make([]string, 0)

		for {
			tok, lit := s.Scan()
			if tok == token.EOF {
				break
			}

			toks = append(toks, tok)
			lits = append(lits, lit)
		}

		if !reflect.DeepEqual(toks, test.toks) {
			t.Errorf("%d. \"%s\" unexpected token: expected=%#v actual=%#v", i, test.str, test.toks, toks)
		} else if !reflect.DeepEqual(lits, test.lits) {
			t.Errorf("%d. \"%s\" unexpected literal: expected=%#v actual=%#v", i, test.str, test.lits, lits)
		}
	}
}

func TestScanner_usage(t *testing.T) {

	str := "foo 123 bar //cool $as fuck ${n} shit"
	s := New(bufio.NewReader(strings.NewReader(str)))

	tok, lit := s.Scan()
	tokEq(tok, token.IDENT, t)
	strEq(lit, "foo", t)

	tok, lit = s.Scan()
	tokEq(tok, token.WS, t)
	strEq(lit, " ", t)

	tok, lit = s.Scan()
	tokEq(tok, token.NUMBER, t)
	strEq(lit, "123", t)

	tok, lit = s.Scan()
	tokEq(tok, token.WS, t)
	strEq(lit, " ", t)

	tok, lit = s.Scan()
	tokEq(tok, token.IDENT, t)
	strEq(lit, "bar", t)

	tok, lit = s.Scan()
	tokEq(tok, token.WS, t)
	strEq(lit, " ", t)

	tok, lit = s.Scan()
	tokEq(tok, token.COMMENT, t)
	strEq(lit, "", t)

	// At this point we want to start processing the interpolated comment string
	// with SeekInterp(). Mimic how the parser grabs the next immediately after processing a
	// previous one. i.e. "process" -> "next", rather than "next" -> "process"
	tok, lit = s.Scan()
	tokEq(tok, token.IDENT, t)
	strEq(lit, "cool", t)

	// Now put back "cool", because we've already consumed it, but we want
	// SeekInterp() to see it.
	s.Unscan()

	tok, b, str := s.SeekInterp()
	tokEq(tok, token.INTERPOLATE_BEGIN, t)
	boolEq(b, false, t)
	strEq(str, "cool ", t)

	// Simulate the interpolated expression as just an identifier, because
	// we're not a fully-fledged parser.
	tok, lit = s.Scan()
	tokEq(tok, token.IDENT, t)
	strEq(lit, "as", t)

	tok, b, str = s.SeekInterp()
	tokEq(tok, token.INTERPOLATE_BEGIN, t)
	boolEq(b, true, t)
	strEq(str, " fuck ", t)

	// Simulate the interpolated expression again
	tok, lit = s.Scan()
	tokEq(tok, token.IDENT, t)
	strEq(lit, "n", t)

	// Should be a close brace now
	tok, lit = s.Scan()
	tokEq(tok, token.INTERPOLATE_END, t)
	strEq(lit, "", t)

	tok, b, str = s.SeekInterp()
	tokEq(tok, token.EOF, t)
	boolEq(b, false, t)
	strEq(str, " shit", t)

	// For shits
	tok, lit = s.Scan()
	tokEq(tok, token.EOF, t)
	strEq(lit, "", t)
}
