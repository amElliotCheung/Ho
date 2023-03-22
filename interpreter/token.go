package interpreter

const (
	// Identifiers + literals
	IDENTIFIER = "IDENTIFIER" // add, x, y, ...
	INTEGER    = "INTEGER"    // 1343456
	STRING     = "STRING"     // "foobar"
	BOOLEAN    = "BOOLEAN"
	// Operators
	OPERATOR = "OPERATOR"
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	MOD      = "%"

	LT  = "<"
	GT  = ">"
	LTE = "<="
	GTE = ">="

	EQ  = "=="
	NEQ = "!="

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"
	COMMA    = ","

	// Keywords
	RESERVED = "RESERVED"
	FUNCTION = "func"
	IF       = "if"
	ELSE     = "else"
	WHILE    = "while"
	TRUE     = "true"
	FALSE    = "false"
	HOPE     = "hope"
	FUZZING  = "fuzzing"
)

type Token interface {
	LineNumber() int
	Type() string
	Literal() string
}

var EOF = helperToken{BaseToken{lineNumber: -1, literal: "EOF"}}
var EOL = helperToken{BaseToken{lineNumber: 0, literal: "EOL"}}

type BaseToken struct {
	lineNumber int
	literal    string
}

func (t BaseToken) LineNumber() int {
	return t.lineNumber
}

func (t BaseToken) Literal() string {
	return t.literal
}

// ===================== Identifier
type IdToken struct {
	BaseToken
}

func (i IdToken) Type() string {
	return IDENTIFIER
}

func NewIdToken(lineNo int, literal string) *IdToken {
	return &IdToken{
		BaseToken: BaseToken{
			lineNumber: lineNo,
			literal:    literal,
		},
	}
}

// ========================== Numer
type NumToken struct {
	BaseToken
}

func NewNumToken(lineNo int, literal string) *NumToken {
	return &NumToken{
		BaseToken: BaseToken{
			lineNumber: lineNo,
			literal:    literal,
		},
	}
}

func (n NumToken) Type() string {
	return INTEGER
}

// =========================== string
type StrToken struct {
	BaseToken
}

func NewStrToken(lineNo int, literal string) *StrToken {
	return &StrToken{
		BaseToken: BaseToken{
			lineNumber: lineNo,
			literal:    literal,
		},
	}
}

func (s StrToken) Type() string {
	return STRING
}

// =========================== bool
type BooleanToken struct {
	BaseToken
}

func NewBooleanToken(lineNo int, literal string) *BooleanToken {
	return &BooleanToken{
		BaseToken: BaseToken{lineNumber: lineNo,
			literal: literal},
	}
}

func (b BooleanToken) Type() string {
	return BOOLEAN
}

// =========================== Operator
type OpToken struct {
	BaseToken
}

func NewOpToken(lineNo int, literal string) *OpToken {
	return &OpToken{
		BaseToken: BaseToken{lineNumber: lineNo,
			literal: literal},
	}
}

func (o OpToken) Type() string {
	return OPERATOR
}

// ======== reserved token
type ReservedToken struct {
	BaseToken
}

func NewReservedToken(lineNo int, literal string) *ReservedToken {
	return &ReservedToken{
		BaseToken: BaseToken{lineNumber: lineNo,
			literal: literal},
	}
}

func (r ReservedToken) Type() string {
	return r.Literal()
}

// ======= helper token
type helperToken struct {
	BaseToken
}

func (r helperToken) Type() string {
	return "helper"
}
