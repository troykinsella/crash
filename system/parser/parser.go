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

Grammar (NEEDS UPDATING)
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

var traceEnabled = false

type Parser struct {
	s *scanner.Scanner

	trace bool
	indent int

	tok token.Token
	lit string
}

func New(r *bufio.Reader) *Parser {
	p := &Parser{
		s: scanner.New(r),
		trace: traceEnabled,
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
	p.tok, p.lit = p.s.Scan()

	if p.tok == token.WS {
		p.tok, p.lit = p.s.Scan()
	}
	if p.trace {
		fmt.Printf("next: %s, literal=%s\n", p.tok, p.lit)
	}
}

func (p *Parser) err(msg string) {
	panic(errors.New(msg))
}

func (p *Parser) errExpected(expected string) {
	ln, col := p.s.Pos()
	p.err(fmt.Sprintf("[%d:%d] found '%s' (%s), expected %s", ln, col, p.lit, p.tok, expected))
}

func (p *Parser) expect(t token.Token) {
	if p.tok != t {
		ln, col := p.s.Pos()
		p.err(fmt.Sprintf("[%d:%d] found '%s' (%s), expected %s", ln, col, p.lit, p.tok, t))
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
	if p.tok == token.EOF {
		return nil
	}

	// First, put back the current token so scanner.SeekInterp() can see it
	p.s.Unscan()

	tok, expectEnd, str := p.s.SeekInterp()
	if tok == token.EOF {
		return &ast.IString{
			Str: str,
		}
	}

	var expr *ast.Expression
	var ident *ast.Identifier
	var next *ast.IString

	p.next()

	if expectEnd {
		e, err := p.Expression()
		if err != nil {
			panic(err) // handled by IString
		}
		expr = e
		p.expect(token.INTERPOLATE_END)

	} else if p.tok != token.EOF {
		ident = p.parseIdentifier()
	}

	next = p.istr()

	return &ast.IString{
		Str: str,
		Expr: expr,
		Ident: ident,
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
	if p.tok == token.COMMENT {
		p.next()
		com, err = p.IString()
		if err != nil {
			return nil, err
		}
	} else {
		p.s.Unscan()
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
	if p.tok == token.NOT {
		n = true
		p.next()
	}

	// Operation
	op := p.lit
	p.expect(token.IDENT)

	// Arguments
	var arguments *ast.ExpressionList
	if p.tok != token.EOF && p.tok != token.COMMENT { // Dip back into statement rules. Doesn't belong here.
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
		if p.tok == token.COMMA {
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
		if p.tok == token.DOT {
			sel := p.parseSelectorExpr(r)
			r = &ast.PrimaryExpr{
				Selector: sel,
			}
		} else if p.tok == token.OPEN_SQUARE_BRACKET {
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

	t := p.tok

	if t == token.IDENT {
		i := p.parseIdentifier()
		pexpr = &ast.PrimaryExpr{
			Ident: i,
		}
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
	i := p.lit
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

func (p *Parser) parseIdentifier() *ast.Identifier {
	if p.trace {
		defer un(trace(p, "Literal"))
	}

	r := &ast.Identifier{
		Name: p.lit,
	}

	p.expect(token.IDENT)
	return r
}

func (p *Parser) parseLiteral() *ast.Literal {
	if p.trace {
		defer un(trace(p, "Literal"))
	}

	var r *ast.Literal

	switch p.tok {
	case token.STRING:
		p.printTrace("Str:", p.lit)
		r = &ast.Literal{
			Type: p.tok,
			Str: p.lit,
		}
	case token.NUMBER:
		i, err := strconv.ParseInt(p.lit, 10, 0)
		if err != nil {
			panic(err)
		}
		p.printTrace("Int:", i)
		r = &ast.Literal{
			Type: p.tok,
			Int: i,
		}
	case token.TRUE, token.FALSE:
		p.printTrace("Bool:", p.tok == token.TRUE)
		r = &ast.Literal{
			Type: p.tok,
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
