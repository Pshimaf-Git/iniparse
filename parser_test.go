package iniparse

import (
	"testing"
)

func TestParseSimpleKeyValue(t *testing.T) {
	input := `name = Alice
age = 30`
	p := newParser()
	ini, err := p.parse(input)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := ini.get("", "name")
	if !ok || val != "Alice" {
		t.Errorf("expected name=Alice, got %q (found=%v)", val, ok)
	}

	val, ok = ini.get("", "age")
	if !ok || val != "30" {
		t.Errorf("expected age=30, got %q (found=%v)", val, ok)
	}
}

func TestParseSection(t *testing.T) {
	input := `[database]
host = localhost
port = 5432

[server]
host = 0.0.0.0
port = 8080`

	p := newParser()
	ini, err := p.parse(input)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := ini.get("database", "host")
	if !ok || val != "localhost" {
		t.Errorf("expected database.host=localhost, got %q", val)
	}

	val, ok = ini.get("server", "port")
	if !ok || val != "8080" {
		t.Errorf("expected server.port=8080, got %q", val)
	}
}

func TestParseComments(t *testing.T) {
	input := `; comment
name = Alice
# another comment
age = 30`

	p := newParser()
	ini, err := p.parse(input)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := ini.get("", "name")
	if !ok || val != "Alice" {
		t.Errorf("expected name=Alice, got %q", val)
	}
}

func TestParseQuotedString(t *testing.T) {
	input := `greeting = "hello world"`
	p := newParser()
	ini, err := p.parse(input)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := ini.get("", "greeting")
	if !ok || val != "hello world" {
		t.Errorf("expected greeting='hello world', got %q", val)
	}
}

func TestParseBooleanValues(t *testing.T) {
	input := `debug = true
verbose = no`
	p := newParser()
	ini, err := p.parse(input)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := ini.get("", "debug")
	if !ok || val != "true" {
		t.Errorf("expected debug=true, got %q", val)
	}

	val, ok = ini.get("", "verbose")
	if !ok || val != "no" {
		t.Errorf("expected verbose=no, got %q", val)
	}
}

func TestParseFloat(t *testing.T) {
	input := `rate = 3.14`
	p := newParser()
	ini, err := p.parse(input)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := ini.get("", "rate")
	if !ok || val != "3.14" {
		t.Errorf("expected rate=3.14, got %q", val)
	}
}

func TestParseInteger(t *testing.T) {
	input := `offset = -42`
	p := newParser()
	ini, err := p.parse(input)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := ini.get("", "offset")
	if !ok || val != "-42" {
		t.Errorf("expected offset=-42, got %q", val)
	}
}

func TestParseEmptySection(t *testing.T) {
	input := `[empty]

name = value`
	p := newParser()
	ini, err := p.parse(input)
	if err != nil {
		t.Fatal(err)
	}

	if !ini.hasSection("empty") {
		t.Error("expected section 'empty' to exist")
	}

	val, ok := ini.get("empty", "name")
	if !ok || val != "value" {
		t.Errorf("expected empty.name=value, got %q", val)
	}
}

func TestParseSections(t *testing.T) {
	input := `[alpha]
x = 1

[beta]
y = 2

[gamma]
z = 3`

	p := newParser()
	ini, err := p.parse(input)
	if err != nil {
		t.Fatal(err)
	}

	secs := ini.sectionNames()
	if len(secs) != 3 {
		t.Fatalf("expected 3 sections, got %d", len(secs))
	}
	expected := []string{"alpha", "beta", "gamma"}
	for i, s := range secs {
		if s != expected[i] {
			t.Errorf("section %d: expected %q, got %q", i, expected[i], s)
		}
	}
}

func TestParseKeys(t *testing.T) {
	input := `[db]
host = localhost
port = 5432
name = test`

	p := newParser()
	ini, err := p.parse(input)
	if err != nil {
		t.Fatal(err)
	}

	k := ini.keys("db")
	if len(k) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(k))
	}
	expected := []string{"host", "port", "name"}
	for i, key := range k {
		if key != expected[i] {
			t.Errorf("key %d: expected %q, got %q", i, expected[i], key)
		}
	}
}

func TestParseMissingEquals(t *testing.T) {
	input := `name Alice`
	p := newParser()
	_, err := p.parse(input)
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}
