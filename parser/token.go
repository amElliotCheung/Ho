package interpreter

import "strconv"

const (
	// Identifiers + literals
	IDENT  = "IDENT"  // add, foobar, x, y, ...
	INT    = "INT"    // 1343456
	STRING = "STRING" // "foobar"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	MOD      = "%"

	LT = "<"
	GT = ">"

	EQ  = "=="
	NEQ = "!="

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

type Token interface {
	GetLineNumber() int
	IsOperator() bool
	IsIdentifier() bool
	IsNumber() bool
	IsString() bool
	GetNumber() int
	GetText() string
	GetType() string
}

var EOF = BaseToken{lineNumber: -1}
var EOL = NewStrToken(0, "\\n")

type BaseToken struct {
	lineNumber int
}

func (t BaseToken) GetLineNumber() int {
	return t.lineNumber
}

func (t BaseToken) GetType() string {
	return ""
}

func (t BaseToken) IsIdentifier() bool {
	return false
}

func (t BaseToken) IsNumber() bool {
	return false
}

func (t BaseToken) IsString() bool {
	return false
}
func (t BaseToken) IsOperator() bool {
	return false
}

func (t BaseToken) GetNumber() int {
	panic("not a number!")
}

func (t BaseToken) GetText() string {
	return ""
}

// ========================== Numer
type NumToken struct {
	BaseToken
	value int
}

func NewNumToken(lineNo, val int) NumToken {
	return NumToken{
		BaseToken: BaseToken{lineNumber: lineNo},
		value:     val,
	}
}

func (t NumToken) IsNumber() bool {
	return true
}

func (t NumToken) GetType() string {
	return "INT"
}

func (t NumToken) GetNumber() int {
	return t.value
}

func (t NumToken) GetText() string {
	return strconv.Itoa(t.value)
}

// ===================== Identifier
type IdToken struct {
	BaseToken
	text string
}

func NewIdToken(lineNo int, text string) IdToken {
	return IdToken{
		BaseToken: BaseToken{lineNumber: lineNo},
		text:      text,
	}
}
func (t IdToken) IsIdentifier() bool {
	return true
}

func (t IdToken) GetText() string {
	return t.text
}

func (t IdToken) GetType() string {
	return "IDENT"
}

// =========================== string
type StrToken struct {
	BaseToken
	literal string
}

func NewStrToken(lineNo int, str string) StrToken {
	return StrToken{
		BaseToken: BaseToken{lineNumber: lineNo},
		literal:   str,
	}
}
func (t StrToken) IsString() bool {
	return true
}

func (t StrToken) GetText() string {
	return t.literal
}

func (t StrToken) GetType() string {
	return "STRING"
}

// =========================== Operator
type OpToken struct {
	BaseToken
	op string
}

func NewOpToken(lineNo int, str string) OpToken {
	return OpToken{
		BaseToken: BaseToken{lineNumber: lineNo},
		op:        str,
	}
}

func (o OpToken) GetText() string {
	return o.op
}

func (o OpToken) GetType() string {
	return o.op
}

func (o OpToken) IsOperator() bool {
	return true
}
