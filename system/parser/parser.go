package parser

import (
	"fmt"
	"bufio"
	"github.com/troykinsella/crash/system/ast"
	"github.com/troykinsella/crash/system/scanner"
	"github.com/troykinsella/crash/system/token"
	"strconv"
	"errors"
)

// Thanks to https://github.com/benbjohnson for ideas

/*

The syntax is specified using Extended Backus-Naur Form (EBNF):

Production  = production_name "=" [ Expression ] "." .
Expression  = Alternative { "|" Alternative } .
Alternative = Term { Term } .
Term        = production_name | token [ "…" token ] | Group | Option | Repetition .
Group       = "(" Expression ")" .
Option      = "[" Expression "]" .
Repetition  = "{" Expression "}" .

Grammar
=======

Statement       = Operation [ Comment ] .
Comment         = "//" .* .
Operation       = ExpressionList Identifier [ ExpressionList ] .

ExpressionList  = Expression { "," Expression } .

Expression      = PrimaryExp

PrimaryExpr     = Operand | SelectorExpr | IndexExpr
Operand         = Literal | identifier
Selector        = "." identifier .
SelectorExpr    = PrimaryExpr Selector
Index           = "[" Expression "]"
IndexExpr       = PrimaryExpr Index



//SubjectList     = Subject { "," Subject } .
//Subject         = Identifier | String .

//ArgumentList    = Argument { "," Argument } .
//Argument        = Identifier | String | Number .

identifier      = ['a'-'z''A'-'Z''-''_''$']+ .
Literal         = number | string .
number          = [+-]?['0'-'9']+ .
string          = '"' .* '"' .
boolean         = "true" | "false" .

//Group           = "(" Operation ")" .

##########

Index         = "[" Expression "]"

AndExpr       = statement "and" statement
or            = statement "or" statement

UnaryExpr     = PrimaryExpr | unary_op UnaryExpr

unary_op      = "+" | "-" | "not" .
binary_op     = "and" | "or" | "xor" .

*/

type tokenScanner struct {
	s   *scanner.Scanner

	buf bool
	tok token.Token
	lit string
}

func (ts *tokenScanner) scan() (token.Token, string) {
	if ts.buf {
		ts.buf = false
		return ts.tok, ts.lit
	}

	ts.tok, ts.lit = ts.s.Scan()
	return ts.tok, ts.lit
}

func (ts *tokenScanner) unscan() {
	ts.buf = true
}

type Parser struct {
	ts *tokenScanner

	trace bool
	indent int
}

func New(r *bufio.Reader) *Parser {
	p := &Parser{
		ts: &tokenScanner{
			s: scanner.New(r),
		},
		trace: false, // TEMP
	}
	p.next()
	return p
}

// Excellent pattern shamelessly lifted from go/parser
func (p *Parser) printTrace(a ...interface{}) {
	if !p.trace {
		return
	}

	const dots = ". . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . "
	const n = len(dots)

	fmt.Printf("%5d:%3d: ", 1, 1)
	i := 2 * p.indent
	for i > n {
		fmt.Print(dots)
		i -= n
	}
	// i <= n
	fmt.Print(dots[0:i])
	fmt.Println(a...)
}

// Excellent pattern shamelessly lifted from go/parser
func trace(p *Parser, msg string) *Parser {
	p.printTrace(msg, "(")
	p.indent++
	return p
}

// Excellent pattern shamelessly lifted from go/parser
// Usage pattern: defer un(trace(p, "..."))
func un(p *Parser) {
	p.indent--
	p.printTrace(")")
}

func (p *Parser) next() {
	p.ts.scan()

	if p.ts.tok == token.WS {
		p.ts.scan()
	}

	if p.trace {
		fmt.Printf("next: %s, literal=%s\n", p.ts.tok, p.ts.lit)
	}
}

func (p *Parser) err(msg string) {
	panic(errors.New(msg))
}

func (p *Parser) errExpected(expected string) {
	ln, col := p.ts.s.Pos()
	p.err(fmt.Sprintf("[%d:%d] found '%s' (%s), expected %s", ln, col, p.ts.lit, p.ts.tok, expected))
}

func (p *Parser) expect(t token.Token) {
	if p.ts.tok != t {
		ln, col := p.ts.s.Pos()
		p.err(fmt.Sprintf("[%d:%d] found '%s' (%s), expected %s", ln, col, p.ts.lit, p.ts.tok, t))
	}
	p.next()
}

func (p *Parser) IString() (a *ast.IString, err error) {
	if p.trace {
		defer un(trace(p, "IString"))
	}

	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	a = p.istr()
	return
}

func (p *Parser) istr() *ast.IString {
	tok, str := p.ts.s.SeekInterp()
	if tok == token.EOF {
		if str == "" {
			return nil
		}
		return &ast.IString{
			Str: str,
		}
	}

	p.next()
	expr, err := p.Expression()
	if err != nil {
		panic(err) // handled by IString
	}

	// Expect } (but don't call next() after)
	if p.ts.tok != token.INTERPOLATE_END {
		ln, col := p.ts.s.Pos()
		p.err(fmt.Sprintf("[%d:%d] found '%s' (%s), expected %s", ln, col, p.ts.lit, p.ts.tok, token.INTERPOLATE_END))
	}

	next := p.istr()

	return &ast.IString{
		Str: str,
		Expr: expr,
		Next: next,
	}
}

func (p *Parser) Statement() (a *ast.Statement, err error) {
	if p.trace {
		defer un(trace(p, "Statement"))
	}

	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	op := p.operation()

	var com *ast.IString
	if p.ts.tok == token.COMMENT {
		// No next() here

		com, err = p.IString()
		if err != nil {
			return nil, err
		}
	} else {
		p.ts.unscan()
	}

	a = &ast.Statement{
		Operation: op,
		Message: com,
	}

	return
}

func (p *Parser) operation() *ast.Operation {
	if p.trace {
		defer un(trace(p, "Operation"))
	}

	// Subjects
	subjectList := p.parseExpressionList()

	// Optional: "not"
	n := false
	if p.ts.tok == token.NOT {
		n = true
		p.next()
	}

	// Operation
	op := p.ts.lit
	p.expect(token.IDENT)

	// Arguments
	var arguments *ast.ExpressionList
	if p.ts.tok != token.EOF && p.ts.tok != token.COMMENT { // Dip back into statement rules. Doesn't belong here.
		arguments = p.parseExpressionList()
	}

	return &ast.Operation{
		Subjects: subjectList,
		Negate: n,
		Operation: op,
		Arguments: arguments,
	}
}

func (p *Parser) parseExpressionList() *ast.ExpressionList {
	if p.trace {
		defer un(trace(p, "ExpressionList"))
	}

	exprs := make([]*ast.Expression, 0)

	for {
		n, err := p.Expression()
		if err != nil {
			panic(err) // handled by p.Statement
		}
		exprs = append(exprs, n)
		if p.ts.tok == token.COMMA {
			p.next()
		} else {
			break
		}
	}

	return &ast.ExpressionList{
		Expressions: exprs,
	}
}

func (p *Parser) Expression() (a *ast.Expression, err error) {
	if p.trace {
		defer un(trace(p, "Expression"))
	}

	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	pexpr := p.parsePrimaryExpr()

	return &ast.Expression{
		PExpr: pexpr,
	}, nil
}

func (p *Parser) parsePrimaryExpr() *ast.PrimaryExpr {
	if p.trace {
		defer un(trace(p, "PrimaryExpr"))
	}

	operand := p.parseOperand()
	r := operand

	for {
		if p.ts.tok == token.DOT {
			sel := p.parseSelectorExpr(r)
			r = &ast.PrimaryExpr{
				Selector: sel,
			}
		} else if p.ts.tok == token.OPEN_SQUARE_BRACKET {
			idx := p.parseIndexExpr(r)
			r = &ast.PrimaryExpr{
				Index: idx,
			}
		} else {
			break
		}
	}

	return r
}

func (p *Parser) parseOperand() *ast.PrimaryExpr {
	if p.trace {
		defer un(trace(p, "Operand"))
	}

	var pexpr *ast.PrimaryExpr

	t := p.ts.tok

	if t == token.IDENT {
		l := p.ts.lit
		p.printTrace("Ident:", l)
		pexpr = &ast.PrimaryExpr{
			Ident: l,
		}
		p.next()
	} else if isLiteral(t) {
		l := p.parseLiteral()
		pexpr = &ast.PrimaryExpr{
			Literal: l,
		}
	} else {
		p.errExpected("identifier or literal")
		return nil
	}

	return pexpr
}

func (p *Parser) parseSelectorExpr(operand *ast.PrimaryExpr) *ast.SelectorExpr {
	if p.trace {
		defer un(trace(p, "SelectorExpr"))
	}

	p.expect(token.DOT)
	i := p.ts.lit
	p.expect(token.IDENT)

	p.printTrace("Ident:", i)

	return &ast.SelectorExpr{
		Operand: operand,
		Ident: i,
	}
}

func (p *Parser) parseIndexExpr(operand *ast.PrimaryExpr) *ast.IndexExpr {
	if p.trace {
		defer un(trace(p, "IndexExpr"))
	}

	p.expect(token.OPEN_SQUARE_BRACKET)

	idx := p.parsePrimaryExpr()

	p.expect(token.CLOSE_SQUARE_BRACKET)

	return &ast.IndexExpr{
		Operand: operand,
		Index: idx,
	}
}

func (p *Parser) parseLiteral() *ast.Literal {
	if p.trace {
		defer un(trace(p, "Literal"))
	}

	var r *ast.Literal

	t := p.ts.tok
	l := p.ts.lit

	switch p.ts.tok {
	case token.STRING:
		p.printTrace("Str:", l)
		r = &ast.Literal{
			Type: t,
			Str: l,
		}
	case token.NUMBER:
		i, err := strconv.ParseInt(l, 10, 0)
		if err != nil {
			panic(err)
		}
		p.printTrace("Int:", i)
		r = &ast.Literal{
			Type: t,
			Int: i,
		}
	case token.TRUE, token.FALSE:
		p.printTrace("Bool:", t == token.TRUE)
		r = &ast.Literal{
			Type: t,
		}
	default:
		p.err("string or number or boolean")
	}

	p.next()
	return r
}

func isLiteral(t token.Token) bool {
	return t == token.STRING || t == token.NUMBER || t == token.TRUE || t == token.FALSE
}
