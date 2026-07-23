package iniparse

import (
	"testing"
)

func TestLexerSimpleKeyValue(t *testing.T) {
	input := "name = Alice"
	lex := newLexer(input)
	tokens, err := lex.tokenize()
	if err != nil {
		t.Fatal(err)
	}

	assertTokenKinds(t, tokens, []tokenKind{
		tokenKey, tokenAssign, tokenValue, tokenEOF,
	})
	assertTokenLiterals(t, tokens, []string{
		"name", "=", "Alice", "",
	})
}

func TestLexerSection(t *testing.T) {
	input := `[database]
host = localhost
port = 5432`

	lex := newLexer(input)
	tokens, err := lex.tokenize()
	if err != nil {
		t.Fatal(err)
	}

	assertTokenKinds(t, tokens, []tokenKind{
		tokenSection, tokenKey, tokenAssign, tokenValue,
		tokenKey, tokenAssign, tokenInteger, tokenEOF,
	})
	if tokens[0].literal != "database" {
		t.Errorf("expected section name %q, got %q", "database", tokens[0].literal)
	}
}

func TestLexerComment(t *testing.T) {
	input := `; this is a comment
name = value
# another comment`

	lex := newLexer(input)
	tokens, err := lex.tokenize()
	if err != nil {
		t.Fatal(err)
	}

	assertTokenKinds(t, tokens, []tokenKind{
		tokenComment, tokenKey, tokenAssign, tokenValue,
		tokenComment, tokenEOF,
	})
}

func TestLexerQuotedString(t *testing.T) {
	input := `greeting = "hello world"`
	lex := newLexer(input)
	tokens, err := lex.tokenize()
	if err != nil {
		t.Fatal(err)
	}

	assertTokenKinds(t, tokens, []tokenKind{
		tokenKey, tokenAssign, tokenString, tokenEOF,
	})
	if tokens[2].literal != "hello world" {
		t.Errorf("expected %q, got %q", "hello world", tokens[2].literal)
	}
}

func TestLexerUnterminatedSection(t *testing.T) {
	input := `[section`
	lex := newLexer(input)
	_, err := lex.tokenize()
	if err == nil {
		t.Fatal("expected error for unterminated section")
	}
}

func TestLexerBooleans(t *testing.T) {
	input := `debug = true
verbose = no`
	lex := newLexer(input)
	tokens, err := lex.tokenize()
	if err != nil {
		t.Fatal(err)
	}

	assertTokenKinds(t, tokens, []tokenKind{
		tokenKey, tokenAssign, tokenBoolean,
		tokenKey, tokenAssign, tokenBoolean, tokenEOF,
	})
}

func TestLexerFloat(t *testing.T) {
	input := `rate = 3.14`
	lex := newLexer(input)
	tokens, err := lex.tokenize()
	if err != nil {
		t.Fatal(err)
	}

	assertTokenKinds(t, tokens, []tokenKind{
		tokenKey, tokenAssign, tokenFloat, tokenEOF,
	})
}

func TestLexerNegativeInteger(t *testing.T) {
	input := `offset = -42`
	lex := newLexer(input)
	tokens, err := lex.tokenize()
	if err != nil {
		t.Fatal(err)
	}

	assertTokenKinds(t, tokens, []tokenKind{
		tokenKey, tokenAssign, tokenInteger, tokenEOF,
	})
	if tokens[2].literal != "-42" {
		t.Errorf("expected %q, got %q", "-42", tokens[2].literal)
	}
}

func TestLexerDefaultSection(t *testing.T) {
	input := `name = Alice
age = 30`
	lex := newLexer(input)
	tokens, err := lex.tokenize()
	if err != nil {
		t.Fatal(err)
	}

	assertTokenKinds(t, tokens, []tokenKind{
		tokenKey, tokenAssign, tokenValue,
		tokenKey, tokenAssign, tokenInteger, tokenEOF,
	})
}

func assertTokenKinds(t *testing.T, tokens []token, expected []tokenKind) {
	t.Helper()
	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}
	for i, tok := range tokens {
		if tok.kind != expected[i] {
			t.Errorf("token %d: expected kind %s, got %s", i, expected[i], tok.kind)
		}
	}
}

func assertTokenLiterals(t *testing.T, tokens []token, expected []string) {
	t.Helper()
	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}
	for i, tok := range tokens {
		if tok.literal != expected[i] {
			t.Errorf("token %d: expected literal %q, got %q", i, expected[i], tok.literal)
		}
	}
}
