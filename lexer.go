package iniparse

import (
	"fmt"
	"strings"
	"unicode"
)

type lexer struct {
	input       []rune
	pos         int
	line        int
	col         int
	tokens      []token
	expectValue bool
}

func newLexer(input string) *lexer {
	return &lexer{
		input: []rune(input),
		line:  1,
		col:   1,
	}
}

func (l *lexer) peek() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *lexer) advance() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	ch := l.input[l.pos]
	l.pos++
	if ch == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	return ch
}

func (l *lexer) skipWhitespace() {
	for l.pos < len(l.input) && unicode.IsSpace(l.input[l.pos]) {
		l.advance()
	}
}

func (l *lexer) skipHorizontalWhitespace() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == ' ' || ch == '\t' || ch == '\r' {
			l.advance()
		} else {
			break
		}
	}
}

func (l *lexer) skipLineComment() bool {
	ch := l.peek()
	if ch == ';' || ch == '#' {
		startLine, startCol := l.line, l.col
		l.advance()
		for l.pos < len(l.input) && l.input[l.pos] != '\n' {
			l.advance()
		}
		l.tokens = append(l.tokens, newToken(tokenComment, "", startLine, startCol))
		return true
	}
	return false
}

func (l *lexer) readUntil(delimiter rune) (string, bool) {
	var sb strings.Builder
	for l.pos < len(l.input) {
		ch := l.peek()
		if ch == delimiter {
			return sb.String(), true
		}
		if ch == '\n' {
			return "", false
		}
		sb.WriteRune(l.advance())
	}
	return sb.String(), false
}

func (l *lexer) readUnquotedValue() string {
	var sb strings.Builder
	for l.pos < len(l.input) {
		ch := l.peek()
		if ch == '\n' || ch == ';' || ch == '#' {
			break
		}
		sb.WriteRune(l.advance())
	}
	return strings.TrimRight(sb.String(), " \t\r")
}

func (l *lexer) tokenizeSection() error {
	startLine, startCol := l.line, l.col
	l.advance() // skip '['

	name, ok := l.readUntil(']')
	if !ok {
		return fmt.Errorf("%d:%d: unterminated section header", startLine, startCol)
	}

	if l.pos < len(l.input) {
		l.advance() // skip ']'
	}

	// skip trailing whitespace until newline or EOF
	for l.pos < len(l.input) && l.input[l.pos] != '\n' {
		ch := l.input[l.pos]
		if ch == ';' || ch == '#' {
			break
		}
		if !unicode.IsSpace(ch) {
			return fmt.Errorf("%d:%d: unexpected character after section header: %c", l.line, l.col, ch)
		}
		l.advance()
	}

	l.tokens = append(l.tokens, newToken(tokenSection, strings.TrimSpace(name), startLine, startCol))
	return nil
}

func (l *lexer) tokenizeAssign() {
	startLine, startCol := l.line, l.col
	l.advance() // skip '='
	l.tokens = append(l.tokens, newToken(tokenAssign, "=", startLine, startCol))
	l.expectValue = true
}

func (l *lexer) tokenizeKey() {
	startLine, startCol := l.line, l.col
	var sb strings.Builder
	for l.pos < len(l.input) {
		ch := l.peek()
		if ch == '=' || unicode.IsSpace(ch) || ch == '\n' || ch == ';' || ch == '#' {
			break
		}
		sb.WriteRune(l.advance())
	}
	l.tokens = append(l.tokens, newToken(tokenKey, sb.String(), startLine, startCol))
}

func (l *lexer) tokenizeValue() error {
	l.skipHorizontalWhitespace()
	if l.pos >= len(l.input) || l.peek() == '\n' {
		return nil
	}

	ch := l.peek()
	if ch == '"' {
		if err := l.tokenizeQuotedString(); err != nil {
			return err
		}
		l.expectValue = false
		return nil
	}

	startLine, startCol := l.line, l.col
	raw := l.readUnquotedValue()
	l.expectValue = false

	switch {
	case raw == "":
		return nil
	case isBoolean(raw):
		l.tokens = append(l.tokens, newToken(tokenBoolean, strings.ToLower(raw), startLine, startCol))
	case isInteger(raw):
		l.tokens = append(l.tokens, newToken(tokenInteger, raw, startLine, startCol))
	case isFloat(raw):
		l.tokens = append(l.tokens, newToken(tokenFloat, raw, startLine, startCol))
	default:
		l.tokens = append(l.tokens, newToken(tokenValue, raw, startLine, startCol))
	}
	return nil
}

func (l *lexer) tokenizeQuotedString() error {
	startLine, startCol := l.line, l.col
	l.advance() // skip opening quote

	var sb strings.Builder
	for l.pos < len(l.input) {
		ch := l.advance()
		if ch == '"' {
			l.tokens = append(l.tokens, newToken(tokenString, sb.String(), startLine, startCol))
			return nil
		}
		sb.WriteRune(ch)
	}
	return fmt.Errorf("%d:%d: unterminated quoted string", startLine, startCol)
}

func (l *lexer) tokenize() ([]token, error) {
	for l.pos < len(l.input) {
		if l.expectValue {
			l.skipHorizontalWhitespace()
		} else {
			l.skipWhitespace()
		}
		if l.pos >= len(l.input) {
			break
		}

		if l.skipLineComment() {
			continue
		}

		ch := l.peek()

		switch ch {
		case '[':
			if err := l.tokenizeSection(); err != nil {
				return nil, err
			}
			l.expectValue = false

		case '=':
			l.tokenizeAssign()

		case '\n':
			l.advance()
			l.expectValue = false

		default:
			if l.expectValue {
				if err := l.tokenizeValue(); err != nil {
					return nil, err
				}
			} else {
				l.tokenizeKey()
			}
		}
	}

	l.tokens = append(l.tokens, newToken(tokenEOF, "", l.line, l.col))
	return l.tokens, nil
}

func isBoolean(s string) bool {
	switch strings.ToLower(s) {
	case "true", "false", "yes", "no", "on", "off":
		return true
	}
	return false
}

func isInteger(s string) bool {
	if s == "" {
		return false
	}
	start := 0
	if s[0] == '-' || s[0] == '+' {
		if len(s) == 1 {
			return false
		}
		start = 1
	}
	for i := start; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

func isFloat(s string) bool {
	if s == "" {
		return false
	}
	start := 0
	if s[0] == '-' || s[0] == '+' {
		if len(s) == 1 {
			return false
		}
		start = 1
	}
	dotSeen := false
	hasDigitBefore := false
	hasDigitAfter := false
	for i := start; i < len(s); i++ {
		if s[i] == '.' {
			if dotSeen {
				return false
			}
			dotSeen = true
			continue
		}
		if s[i] < '0' || s[i] > '9' {
			return false
		}
		if dotSeen {
			hasDigitAfter = true
		} else {
			hasDigitBefore = true
		}
	}
	return dotSeen && hasDigitBefore && hasDigitAfter
}
