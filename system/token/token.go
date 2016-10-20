package token

type Token uint8

const (
	ILLEGAL Token = iota
	EOF
	WS

	IDENT
	NUMBER
	STRING

	COMMA
	COMMENT

	OPEN_BRACKET
	CLOSE_BRACKET

	OPEN_SQUARE_BRACKET
	CLOSE_SQUARE_BRACKET

	TRUE
	FALSE

	DOT
	AND
	OR
	XOR
	NOT

	INTERPOLATE_BEGIN
	INTERPOLATE_END

	ESCAPE
)

func (t Token) String() string {
	switch t {
	case ILLEGAL: return "illegal"
	case EOF: return "eof"
	case WS: return "whitespace"
	case IDENT: return "identifier"
	case NUMBER: return "number"
	case STRING: return "string"
	case COMMA: return "comma"
	case COMMENT: return "comment"
	case OPEN_BRACKET: return "("
	case CLOSE_BRACKET: return ")"
	case OPEN_SQUARE_BRACKET: return "["
	case CLOSE_SQUARE_BRACKET: return "]"
	case TRUE: return "true"
	case FALSE: return "false"
	case AND: return "and"
	case OR: return "or"
	case XOR: return "xor"
	case NOT: return "not"
	case INTERPOLATE_BEGIN: return "interpolate_begin"
	case INTERPOLATE_END: return "interpolate_end"
	case ESCAPE: return "escape"
	}
	return ""
}
