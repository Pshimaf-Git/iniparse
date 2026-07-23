package iniparse

import "fmt"

type tokenKind int

const (
	tokenEOF tokenKind = iota

	tokenSection
	tokenKey
	tokenAssign
	tokenValue
	tokenComment

	tokenString
	tokenInteger
	tokenFloat
	tokenBoolean
)

var tokenKindNames = map[tokenKind]string{
	tokenEOF:     "<eof>",
	tokenSection: "<section>",
	tokenKey:     "<key>",
	tokenAssign:  "<assign>",
	tokenValue:   "<value>",
	tokenComment: "<comment>",
	tokenString:  "<string>",
	tokenInteger: "<integer>",
	tokenFloat:   "<float>",
	tokenBoolean: "<boolean>",
}

func (tk tokenKind) String() string {
	if name, ok := tokenKindNames[tk]; ok {
		return name
	}
	return "<unknown>"
}

type token struct {
	kind    tokenKind
	literal string
	line    int
	col     int
}

func newToken(kind tokenKind, lit string, line, col int) token {
	return token{
		kind:    kind,
		literal: lit,
		line:    line,
		col:     col,
	}
}

func (t token) String() string {
	return fmt.Sprintf("token{%s %q at %d:%d}", t.kind, t.literal, t.line, t.col)
}
