package token

import "strconv"

type Token interface {
	GetLineNumber() int
	IsIdentifier() bool
	IsNumber() bool
	IsString() bool
	GetNumber() int
	GetText() string
}

const (
	ident = iota
	str
)

type BaseToken struct {
	lineNumber int
}

func (t BaseToken) GetLineNumber() int {
	return t.lineNumber
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

func (t BaseToken) GetNumber() int {
	panic("not a number!")
}

func (t BaseToken) getText() string {
	return ""
}

// ========================== Numer
type NumToken struct {
	BaseToken
	value int
}

func (t NumToken) IsNumber() bool {
	return true
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

func (t IdToken) IsIdentifier() bool {
	return true
}

func (t IdToken) GetText() string {
	return t.text
}

// =========================== string
type StrToken struct {
	BaseToken
	literal string
}

func (t StrToken) IsString() bool {
	return true
}

func (t StrToken) GetText() string {
	return t.literal
}
