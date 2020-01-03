package token

type Token struct {
	Type   int
	Value  interface{}
	Lexeme string
	Line   int
}

const (
	Number = iota
	Plus
	Minus
	Multiplication
	Division
	Space
	Print
	Semicolon
	Identifier
	Let
	EOF
	Equal
	LeftBrace
	RightBrace
	EqualEqual
	BangEqual
	Bang
	Greater
	GreaterEqual
	Less
	LessEqual
	RightParen
	LeftParen
	String
	If
	Else
	Or
	And
	Boolean
	While
	For
	Comma
	Function
	Return
)
