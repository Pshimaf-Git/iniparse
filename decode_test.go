package iniparse

import (
	"testing"
)

func TestDecodeSimpleStruct(t *testing.T) {
	type User struct {
		Name string `ini:"name"`
		Age  int    `ini:"age"`
	}

	input := `name = Alice
age = 30`

	var u User
	err := Unmarshal(b(input), &u)
	if err != nil {
		t.Fatal(err)
	}

	if u.Name != "Alice" {
		t.Errorf("expected Name=Alice, got %q", u.Name)
	}
	if u.Age != 30 {
		t.Errorf("expected Age=30, got %d", u.Age)
	}
}

func TestDecodeSection(t *testing.T) {
	type Database struct {
		Host string `ini:"host"`
		Port int    `ini:"port"`
	}

	type Server struct {
		Host string `ini:"host"`
		Port string `ini:"port"`
	}

	type Config struct {
		DB   Database `ini:"database"`
		Serv Server   `ini:"server"`
	}

	input := `
[server]
host = 0.0.0.0
port = :80

[database]
host = localhost
port = 5432`

	var cfg Config
	err := Unmarshal(b(input), &cfg)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.DB.Host != "localhost" {
		t.Errorf("expected Host=localhost, got %q", cfg.DB.Host)
	}
	if cfg.DB.Port != 5432 {
		t.Errorf("expected Port=5432, got %d", cfg.DB.Port)
	}

	if cfg.Serv.Host != "0.0.0.0" {
		t.Errorf("expected Host=localhost, got %q", cfg.Serv.Host)
	}
	if cfg.Serv.Port != ":80" {
		t.Errorf("expected Port=5432, got %s", cfg.Serv.Port)
	}
}

func TestDecodeBoolean(t *testing.T) {
	type Settings struct {
		Debug   bool `ini:"debug"`
		Verbose bool `ini:"verbose"`
	}

	input := `debug = true
verbose = no`

	var s Settings
	err := Unmarshal(b(input), &s)
	if err != nil {
		t.Fatal(err)
	}

	if !s.Debug {
		t.Error("expected Debug=true")
	}
	if s.Verbose {
		t.Error("expected Verbose=false")
	}
}

func TestDecodeFloat(t *testing.T) {
	type Metrics struct {
		Rate float64 `ini:"rate"`
	}

	input := `rate = 3.14`

	var m Metrics
	err := Unmarshal(b(input), &m)
	if err != nil {
		t.Fatal(err)
	}

	if m.Rate != 3.14 {
		t.Errorf("expected Rate=3.14, got %f", m.Rate)
	}
}

func TestDecodeNegativeInt(t *testing.T) {
	type Offset struct {
		Value int `ini:"value"`
	}

	input := `value = -42`

	var o Offset
	err := Unmarshal(b(input), &o)
	if err != nil {
		t.Fatal(err)
	}

	if o.Value != -42 {
		t.Errorf("expected Value=-42, got %d", o.Value)
	}
}

func TestDecodeUint(t *testing.T) {
	type Config struct {
		Port uint16 `ini:"port"`
	}

	input := `port = 8080`

	var cfg Config
	err := Unmarshal(b(input), &cfg)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Port != 8080 {
		t.Errorf("expected Port=8080, got %d", cfg.Port)
	}
}

func TestDecodeSlice(t *testing.T) {
	type Config struct {
		Tags []string `ini:"tags"`
	}

	input := `tags = go,ini,parser`

	var cfg Config
	err := Unmarshal(b(input), &cfg)
	if err != nil {
		t.Fatal(err)
	}

	if len(cfg.Tags) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(cfg.Tags))
	}
	expected := []string{"go", "ini", "parser"}
	for i, tag := range cfg.Tags {
		if tag != expected[i] {
			t.Errorf("tag %d: expected %q, got %q", i, expected[i], tag)
		}
	}
}

func TestDecodePointerToStruct(t *testing.T) {
	type Inner struct {
		Value string `ini:"value"`
	}

	type Outer struct {
		Inner *Inner `ini:"section"`
	}

	input := `[section]
value = hello`

	var o Outer
	err := Unmarshal(b(input), &o)
	if err != nil {
		t.Fatal(err)
	}

	if o.Inner == nil {
		t.Fatal("expected Inner to be non-nil")
	}
	if o.Inner.Value != "hello" {
		t.Errorf("expected Value=hello, got %q", o.Inner.Value)
	}
}

func TestDecodeTagWithSection(t *testing.T) {
	type Config struct {
		Host string `ini:"server.host"`
		Port int    `ini:"server.port"`
	}

	input := `[server]
host = 0.0.0.0
port = 8080`

	var cfg Config
	err := Unmarshal(b(input), &cfg)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Host != "0.0.0.0" {
		t.Errorf("expected Host=0.0.0.0, got %q", cfg.Host)
	}
	if cfg.Port != 8080 {
		t.Errorf("expected Port=8080, got %d", cfg.Port)
	}
}

func TestDecodeSkipsMissingKeys(t *testing.T) {
	type Config struct {
		Name    string `ini:"name"`
		Missing string `ini:"missing"`
	}

	input := `name = Alice`

	var cfg Config
	err := Unmarshal(b(input), &cfg)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Name != "Alice" {
		t.Errorf("expected Name=Alice, got %q", cfg.Name)
	}
	if cfg.Missing != "" {
		t.Errorf("expected Missing to be empty, got %q", cfg.Missing)
	}
}

func TestDecodeSkipsIgnoredFields(t *testing.T) {
	type Config struct {
		Name    string `ini:"name"`
		Ignored string `ini:"-"`
	}

	input := `name = Alice
ignored = should not be set`

	var cfg Config
	err := Unmarshal(b(input), &cfg)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Name != "Alice" {
		t.Errorf("expected Name=Alice, got %q", cfg.Name)
	}
	if cfg.Ignored != "" {
		t.Errorf("expected Ignored to be empty, got %q", cfg.Ignored)
	}
}

func TestDecodeNilPointer(t *testing.T) {
	type Inner struct {
		Value string `ini:"value"`
	}

	type Outer struct {
		Inner *Inner `ini:"section"`
	}

	input := `[section]
value = hello`

	var o Outer
	err := Unmarshal(b(input), &o)
	if err != nil {
		t.Fatal(err)
	}

	if o.Inner == nil {
		t.Fatal("expected Inner to be non-nil after unmarshal")
	}
	if o.Inner.Value != "hello" {
		t.Errorf("expected Value=hello, got %q", o.Inner.Value)
	}
}

func TestDecodeNonStructPointer(t *testing.T) {
	var s string
	err := Unmarshal(b("key = value"), &s)
	if err == nil {
		t.Fatal("expected error for non-struct pointer")
	}
}

func TestDecodeNilPointerArg(t *testing.T) {
	err := Unmarshal(b("key = value"), nil)
	if err == nil {
		t.Fatal("expected error for nil pointer")
	}
}

func TestDecodeMultipleSections(t *testing.T) {
	type Database struct {
		Host string `ini:"host"`
	}

	type Server struct {
		Host string `ini:"host"`
	}

	type Config struct {
		DB     Database `ini:"database"`
		Server Server   `ini:"server"`
	}

	input := `[database]
host = db.example.com

[server]
host = api.example.com`

	var cfg Config
	err := Unmarshal(b(input), &cfg)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.DB.Host != "db.example.com" {
		t.Errorf("expected Db.Host=db.example.com, got %q", cfg.DB.Host)
	}
	if cfg.Server.Host != "api.example.com" {
		t.Errorf("expected Server.Host=api.example.com, got %q", cfg.Server.Host)
	}
}

func b(s string) []byte {
	return []byte(s)
}
