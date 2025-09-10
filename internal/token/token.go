package token

type TokenKind string

type Token struct {
	Kind    TokenKind
	Literal string
	Line    int
	Column  int
}

const (
	ILLEGAL   TokenKind = "ILLEGAL"
	INDEX     TokenKind = "INDEX"
	TIMESTAMP TokenKind = "TIMESTAMP"
	ARROW     TokenKind = "ARROW"
	TEXT      TokenKind = "TEXT"
	LF        TokenKind = "LF"
	EOC       TokenKind = "EOC"
	EOF       TokenKind = "EOF"
)

func NewToken(kind TokenKind, literal string, line, column int) Token {
	return Token{
		Kind:    kind,
		Literal: literal,
		Line:    line,
		Column:  column,
	}
}
