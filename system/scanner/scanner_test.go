package scanner

import (
	"testing"
	"strings"
	"bufio"
	"github.com/troykinsella/crash/system/token"
)

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
		{"${",     token.INTERPOLATE_BEGIN,    ""},
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
	}{
		{"",   token.EOF, ""},
		{"$",  token.EOF, "$"},
		{"${", token.INTERPOLATE_BEGIN, ""},
		{" ${", token.INTERPOLATE_BEGIN, " "},
		{"asdf${", token.INTERPOLATE_BEGIN, "asdf"},
		{"$${", token.INTERPOLATE_BEGIN, "$"},
		{"$$${", token.INTERPOLATE_BEGIN, "$$"},
		{"{${", token.INTERPOLATE_BEGIN, "{"},
	}

	for i, test := range tests {
		s := New(bufio.NewReader(strings.NewReader(test.str)))
		tok, text := s.SeekInterp()
		if tok != test.tok {
			t.Errorf("%d. unexpected token: expected=%q actual=%q (%q)", i, test.tok, tok, text)
		} else if text != test.text {
			t.Errorf("%d. unexpected text: expected=%q actual=%q", i, test.text, text)
		}

		if tok, lit := s.Scan(); tok != token.EOF {
			t.Errorf("%d. expected eof: actual=%q", i, lit)
		}
	}
}
